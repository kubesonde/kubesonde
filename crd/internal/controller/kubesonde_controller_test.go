/*
Copyright 2025.

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
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	kubesondev1 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/state"
)

// Mock for testing
type MockClient struct {
	client.Client
	mock.Mock
}

func TestKubesondeReconciler(t *testing.T) {
	t.Run("Test Reconcile with not found error", func(t *testing.T) {
		// Create a fake client
		scheme := runtime.NewScheme()
		_ = kubesondev1.AddToScheme(scheme)

		fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

		// Create reconciler
		reconciler := &KubesondeReconciler{
			Client:           fakeClient,
			Log:              logr.Discard(),
			Scheme:           scheme,
			KubernetesClient: kubernetesfake.NewSimpleClientset(),
		}

		// Create a request for a non-existent resource
		req := ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      "non-existent",
				Namespace: "default",
			},
		}

		// Call Reconcile
		result, err := reconciler.Reconcile(context.Background(), req)

		// Verify no error and no requeue
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
	})

	t.Run("Test Reconcile clears probe state when resource is deleted", func(t *testing.T) {
		// Ensure a clean singleton for this test.
		state.ResetDefaultManager()

		// Seed some probe state as if probing had already run.
		state.SetProbeState(&kubesondev1.ProbeOutput{
			Items: []kubesondev1.ProbeOutputItem{
				{Port: "80", Protocol: "TCP"},
			},
		})
		require.NotEmpty(t, state.GetProbeState().Items, "precondition: state should be seeded")

		scheme := runtime.NewScheme()
		_ = kubesondev1.AddToScheme(scheme)

		// Empty client => Get returns NotFound, simulating a deleted resource.
		fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

		reconciler := &KubesondeReconciler{
			Client:           fakeClient,
			Log:              logr.Discard(),
			Scheme:           scheme,
			KubernetesClient: kubernetesfake.NewSimpleClientset(),
		}

		req := ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      "deleted-kubesonde",
				Namespace: "default",
			},
		}

		result, err := reconciler.Reconcile(context.Background(), req)

		// No error, no requeue, and the probe state is cleared.
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{}, result)
		assert.Empty(t, state.GetProbeState().Items, "probe state should be cleared after deletion")
	})

	t.Run("Test SetupWithManager", func(t *testing.T) {
		// Create a fake client
		scheme := runtime.NewScheme()
		_ = kubesondev1.AddToScheme(scheme)

		fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

		// Create reconciler
		reconciler := &KubesondeReconciler{
			Client:           fakeClient,
			Log:              logr.Discard(),
			Scheme:           scheme,
			KubernetesClient: kubernetesfake.NewSimpleClientset(),
		}

		// Create a fake manager
		// This test mostly ensures that the method doesn't panic
		err := reconciler.SetupWithManager(nil)
		// We expect an error since we're passing nil manager
		assert.Error(t, err)
	})
}

func TestKubesondeReconcilerWithValidResource(t *testing.T) {
	t.Run("Test Reconcile with valid resource", func(t *testing.T) {
		// Create a fake client
		scheme := runtime.NewScheme()
		_ = kubesondev1.AddToScheme(scheme)

		// Create a Kubesonde resource
		kubesonde := &kubesondev1.Kubesonde{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-kubesonde",
				Namespace: "default",
			},
			Spec: kubesondev1.KubesondeSpec{
				Namespace: "default",
				Probe:     "all",
			},
		}

		fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(kubesonde).
			WithStatusSubresource(kubesonde).Build()

		// Create reconciler
		reconciler := &KubesondeReconciler{
			Client:           fakeClient,
			Log:              logr.Discard(),
			Scheme:           scheme,
			KubernetesClient: kubernetesfake.NewSimpleClientset(),
		}

		// Create a request
		req := ctrl.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-kubesonde",
				Namespace: "default",
			},
		}

		// Call Reconcile
		result, err := reconciler.Reconcile(context.Background(), req)

		// Verify no error and that reconcile requeues to refresh status.
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{RequeueAfter: statusRefreshInterval}, result)
	})
}
