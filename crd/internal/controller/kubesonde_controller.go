/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"sync"
	"time"

	kubesondev1 "kubesonde.io/api/v1"
	kubesondeDispatcher "kubesonde.io/controllers/dispatcher"
	kubesondeEvents "kubesonde.io/controllers/events"
	kubesondemetrics "kubesonde.io/controllers/metrics"
	kubesondemonitor "kubesonde.io/controllers/monitor"
	"kubesonde.io/controllers/state"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	recursiveprobing "kubesonde.io/controllers/recursive-probing"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// dispatcherOnce ensures the probe worker pool is started only once, regardless
// of how many times Reconcile runs.
var dispatcherOnce sync.Once

// startupOnce ensures the per-object probing/events/monitor goroutines are
// started only once. Reconcile is requeued periodically to refresh status, and
// those long-running goroutines must not be re-spawned on every requeue.
var startupOnce sync.Once

// statusRefreshInterval is how often Reconcile is requeued to refresh the
// completeness status shown by `kubectl get kubesondes`.
const statusRefreshInterval = 15 * time.Second

// KubesondeReconciler reconciles a Kubesonde object
type KubesondeReconciler struct {
	client.Client
	Log              logr.Logger
	Scheme           *runtime.Scheme
	KubernetesClient kubernetes.Interface
	// TODO: Add fake clock  for testing purposes
}

// +kubebuilder:rbac:groups=*,resources=*,verbs=*

func (r *KubesondeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Kubesonde-controller", req.NamespacedName)
	apiClient := r.KubernetesClient

	var Kubesonde kubesondev1.Kubesonde
	if err := r.Get(ctx, req.NamespacedName, &Kubesonde); err != nil {
		log.Error(err, "unable to fetch Kubesonde")
		// Ignore not found errors as we do not want to support this usecase
		// Use 	apierrors "k8s.io/apimachinery/pkg/api/errors" to find out if resource was deleted.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// TODOs
	/*
		1) Handle pod deletion. When a pod is deleted also the probes that regard that pod should be removed
		2) Handle resource deletion. When a kubesonde resource is removed, the state should be cleared
	*/

	// Dispatcher. Start the worker pool only once across reconciles; it drains a
	// process-wide shared queue and must not be started multiple times.
	dispatcherOnce.Do(func() {
		go kubesondeDispatcher.Run(context.Background(), apiClient)
	})

	// Start the long-running per-object goroutines only once; requeues below are
	// only for refreshing status and must not re-spawn these loops.
	startupOnce.Do(func() {
		// Events
		go kubesondeEvents.InitEventListener(apiClient, Kubesonde)

		// Probing
		go recursiveprobing.RecursiveProbing(Kubesonde, 20*time.Second)

		// Monitor
		go kubesondemonitor.RunMonitorContainers(apiClient)
	})

	// Refresh the completeness status so it is visible via `kubectl get kubesondes`.
	completeness := state.GetCompleteness()
	Kubesonde.Status.Complete = completeness.Complete
	Kubesonde.Status.ProbeCount = completeness.Count
	if completeness.Count > 0 {
		now := metav1.Now()
		Kubesonde.Status.LastProbeTime = &now
	}
	if err := r.Status().Update(ctx, &Kubesonde); err != nil {
		log.Error(err, "unable to update Kubesonde status")
		return ctrl.Result{RequeueAfter: statusRefreshInterval}, nil
	}

	return ctrl.Result{RequeueAfter: statusRefreshInterval}, nil
}

func (r *KubesondeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubesondev1.Kubesonde{}).
		Complete(r)
}

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(kubesondemetrics.MetricsSummary)
	metrics.Registry.MustRegister(kubesondemetrics.DurationSummary)
	metrics.Registry.MustRegister(kubesondemetrics.TargetedMetricsSummary)
}
