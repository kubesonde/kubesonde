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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	recursiveprobing "kubesonde.io/controllers/recursive-probing"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// probingMu guards probingCancel, which tracks the lifecycle of the background
// probing goroutines (dispatcher, event listener, recursive probing, monitor).
// They are started when a Kubesonde exists and cancelled when it is deleted, so
// probing actually stops (rather than refilling the cleared state).
var (
	probingMu     sync.Mutex
	probingCancel context.CancelFunc
)

// startProbing launches the background goroutines under a fresh cancelable
// context. It is idempotent: if probing is already running it does nothing.
func startProbing(apiClient kubernetes.Interface, k kubesondev1.Kubesonde) {
	probingMu.Lock()
	defer probingMu.Unlock()
	if probingCancel != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	probingCancel = cancel

	go kubesondeDispatcher.Run(ctx, apiClient)
	go kubesondeEvents.InitEventListener(ctx, apiClient, k)
	go recursiveprobing.RecursiveProbing(ctx, k, 20*time.Second)
	go kubesondemonitor.RunMonitorContainers(ctx, apiClient)
}

// stopProbing cancels the background goroutines if they are running.
func stopProbing() {
	probingMu.Lock()
	defer probingMu.Unlock()
	if probingCancel != nil {
		probingCancel()
		probingCancel = nil
	}
}

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
		if apierrors.IsNotFound(err) {
			// The Kubesonde resource was deleted: stop the background probing
			// goroutines and clear the accumulated state so it is not refilled
			// and no stale results are served by the REST API.
			log.Info("Kubesonde resource deleted, stopping probing and clearing state")
			stopProbing()
			state.ClearState()

			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch Kubesonde")
		return ctrl.Result{}, err
	}
	// Start (or keep running) the background probing goroutines for this
	// Kubesonde. Idempotent across the periodic status requeues.
	startProbing(apiClient, Kubesonde)

	// Refresh the completeness status so it is visible via `kubectl get kubesondes`.
	completeness := state.GetCompleteness()
	Kubesonde.Status.Complete = completeness.Complete
	Kubesonde.Status.ProbeCount = completeness.Count
	if completeness.Count > 0 {
		now := metav1.Now()
		Kubesonde.Status.LastProbeTime = &now
	}

	// Mirror completeness into a standard Condition so users can block on it with
	// `kubectl wait --for=condition=Complete kubesonde/<name>`.
	condStatus := metav1.ConditionFalse
	reason := "Probing"
	message := "Probing is in progress"
	if completeness.Complete {
		condStatus = metav1.ConditionTrue
		reason = "Quiesced"
		message = "Recorded probe count has been stable for the completeness window"
	}
	meta.SetStatusCondition(&Kubesonde.Status.Conditions, metav1.Condition{
		Type:               kubesondev1.ConditionComplete,
		Status:             condStatus,
		ObservedGeneration: Kubesonde.Generation,
		Reason:             reason,
		Message:            message,
	})

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
