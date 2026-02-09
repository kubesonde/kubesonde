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

package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKubesondeSpec(t *testing.T) {
	t.Run("Test KubesondeSpec fields", func(t *testing.T) {
		spec := KubesondeSpec{
			DebuggerImage: "test-debugger-image",
			MonitorImage:  "test-monitor-image",
			Namespace:     "test-namespace",
			Probe:         "all",
		}

		assert.Equal(t, "test-debugger-image", spec.DebuggerImage)
		assert.Equal(t, "test-monitor-image", spec.MonitorImage)
		assert.Equal(t, "test-namespace", spec.Namespace)
		assert.Equal(t, "all", spec.Probe)
	})

	t.Run("Test KubesondeSpec with exclude and include", func(t *testing.T) {
		exclude := []ExcludedItem{
			{
				FromPodSelector: "from-pod-selector",
				ToPodSelector:   "to-pod-selector",
				Port:            "8080",
				Protocol:        "TCP",
			},
		}

		include := []IncludedItem{
			{
				FromPodSelector: "from-pod-selector",
				ToPodSelector:   "to-pod-selector",
				Port:            "8080",
				Protocol:        "TCP",
				ExpectedAction:  ALLOW,
			},
		}

		spec := KubesondeSpec{
			Exclude: exclude,
			Include: include,
		}

		assert.Len(t, spec.Exclude, 1)
		assert.Len(t, spec.Include, 1)
		assert.Equal(t, "from-pod-selector", spec.Exclude[0].FromPodSelector)
		assert.Equal(t, "to-pod-selector", spec.Include[0].ToPodSelector)
		assert.Equal(t, ALLOW, spec.Include[0].ExpectedAction)
	})
}

func TestKubesondeStatus(t *testing.T) {
	t.Run("Test KubesondeStatus fields", func(t *testing.T) {
		// Test with nil LastProbeTime
		status := KubesondeStatus{}
		assert.Nil(t, status.LastProbeTime)

		// Test with a specific time
		now := metav1.Now()
		status = KubesondeStatus{
			LastProbeTime: &now,
		}
		assert.NotNil(t, status.LastProbeTime)
		assert.Equal(t, now, *status.LastProbeTime)
	})
}

func TestKubesonde(t *testing.T) {
	t.Run("Test Kubesonde object creation", func(t *testing.T) {
		// Test basic object creation
		kubesonde := Kubesonde{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Kubesonde",
				APIVersion: "security.kubesonde.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-kubesonde",
				Namespace: "default",
			},
			Spec: KubesondeSpec{
				Namespace: "default",
				Probe:     "all",
			},
			Status: KubesondeStatus{},
		}

		assert.Equal(t, "test-kubesonde", kubesonde.Name)
		assert.Equal(t, "default", kubesonde.Namespace)
		assert.Equal(t, "security.kubesonde.io/v1", kubesonde.APIVersion)
		assert.Equal(t, "Kubesonde", kubesonde.Kind)
	})
}

func TestActionTypeConstants(t *testing.T) {
	t.Run("Test ActionType constants", func(t *testing.T) {
		assert.Equal(t, ALLOW, ActionType("Allow"))
		assert.Equal(t, DENY, ActionType("Deny"))
	})
}

func TestProbeTypeConstants(t *testing.T) {
	t.Run("Test ProbeType constants", func(t *testing.T) {
		assert.Equal(t, ALL, ProbeType("all"))
		assert.Equal(t, NONE, ProbeType("none"))
	})
}
