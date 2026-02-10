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
	"time"

	kubesondev1 "kubesonde.io/api/v1"
	kubesondeDispatcher "kubesonde.io/controllers/dispatcher"
	kubesondeEvents "kubesonde.io/controllers/events"
	kubesondemetrics "kubesonde.io/controllers/metrics"
	kubesondemonitor "kubesonde.io/controllers/monitor"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	recursiveprobing "kubesonde.io/controllers/recursive-probing"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

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

	// Dispatcher
	go kubesondeDispatcher.Run(apiClient)

	// Events
	go kubesondeEvents.InitEventListener(apiClient, Kubesonde)

	// Probing
	go recursiveprobing.RecursiveProbing(Kubesonde, 20*time.Second)

	// Monitor
	go kubesondemonitor.RunMonitorContainers(apiClient)

	return ctrl.Result{}, nil
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
