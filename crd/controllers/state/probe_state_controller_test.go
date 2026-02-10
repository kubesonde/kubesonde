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

package state

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "kubesonde.io/api/v1"
)

func TestStateManagerCreation(t *testing.T) {
	t.Run("Test NewStateManager creation", func(t *testing.T) {
		sm := NewStateManager()
		assert.NotNil(t, sm)
		assert.NotNil(t, sm.probeOutput.Items)
		assert.NotNil(t, sm.probeOutput.Errors)
		assert.NotNil(t, sm.probeOutput.PodNetworking)
		assert.NotNil(t, sm.probeOutput.PodNetworkingV2)
		assert.NotNil(t, sm.probeOutput.PodConfigurationNetworking)
		assert.NotNil(t, sm.podsWithNetstat)
		assert.Equal(t, defaultLockTimeout, sm.lockTimeout)
	})

	t.Run("Test GetDefaultManager", func(t *testing.T) {
		sm1 := GetDefaultManager()
		sm2 := GetDefaultManager()
		assert.Equal(t, sm1, sm2) // Should return the same instance (singleton)
		assert.NotNil(t, sm1)
	})
}

func TestStateManagerNetstatOperations(t *testing.T) {
	t.Run("Test SetNestatPod and GetNetstatPods", func(t *testing.T) {
		sm := NewStateManager()

		// Add pods
		sm.SetNestatPod("pod1")
		sm.SetNestatPod("pod2")

		// Get pods
		pods := sm.GetNetstatPods()
		assert.Len(t, pods, 2)
		assert.Contains(t, pods, "pod1")
		assert.Contains(t, pods, "pod2")
	})

	t.Run("Test DeleteNetstatPod", func(t *testing.T) {
		sm := NewStateManager()

		// Add pods
		sm.SetNestatPod("pod1")
		sm.SetNestatPod("pod2")
		sm.SetNestatPod("pod3")

		// Delete one pod
		sm.DeleteNetstatPod("pod2")

		// Check remaining pods
		pods := sm.GetNetstatPods()
		assert.Len(t, pods, 2)
		assert.Contains(t, pods, "pod1")
		assert.Contains(t, pods, "pod3")
		assert.NotContains(t, pods, "pod2")
	})

	t.Run("concurrent access to netstat operations", func(t *testing.T) {
		sm := NewStateManager()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			sm.SetNestatPod("concurrent-pod-1")
		}()

		go func() {
			defer wg.Done()
			sm.SetNestatPod("concurrent-pod-2")
		}()

		wg.Wait()

		pods := sm.GetNetstatPods()
		assert.Len(t, pods, 2)
	})

}

func TestStateManagerProbeStateOperations(t *testing.T) {
	t.Run("Test SetProbeState and GetProbeState", func(t *testing.T) {
		sm := NewStateManager()

		// Create test probe output
		probeOutput := &v1.ProbeOutput{
			Items: []v1.ProbeOutputItem{
				{
					Type: v1.PROBE,
					Source: v1.ProbeEndpointInfo{
						Type: v1.POD,
						Name: "correct",
					},
					Destination: v1.ProbeEndpointInfo{
						Type: v1.POD,
					},
					Protocol:        "TCP",
					Port:            "80",
					ResultingAction: v1.ALLOW,
				},
			},
			Errors: []v1.ProbeOutputError{
				{
					Value: v1.ProbeOutputItem{
						Type: v1.PROBE,
						Source: v1.ProbeEndpointInfo{
							Type:      v1.POD,
							IPAddress: "error-ip",
						},
						Destination: v1.ProbeEndpointInfo{
							Type: v1.POD,
						},
						Protocol:        "UDP",
						Port:            "53",
						ResultingAction: v1.DENY,
					},
					Reason: "test error",
				},
			},
		}

		// Set probe state
		err := sm.SetProbeState(probeOutput)
		assert.NoError(t, err)

		// Get probe state
		state := sm.GetProbeState()
		assert.Len(t, state.Items, 1)
		assert.Len(t, state.Errors, 1)
		assert.Equal(t, "correct", state.Items[0].Source.Name)
		assert.Equal(t, "error-ip", state.Errors[0].Value.Source.IPAddress)
	})

	t.Run("Test SetProbeState with nil input", func(t *testing.T) {
		sm := NewStateManager()
		err := sm.SetProbeState(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "probes cannot be nil")
	})
}

func TestStateManagerAppendOperations(t *testing.T) {
	t.Run("Test AppendProbes", func(t *testing.T) {
		sm := NewStateManager()

		// Add initial probe
		initialItems := []v1.ProbeOutputItem{
			{
				Type: v1.PROBE,
				Source: v1.ProbeEndpointInfo{
					Type: v1.POD,
				},
				Destination: v1.ProbeEndpointInfo{
					Type: v1.POD,
				},
				Protocol:        "TCP",
				Port:            "80",
				ResultingAction: v1.ALLOW,
			},
		}

		err := sm.AppendProbes(&initialItems)
		assert.NoError(t, err)

		// Add more probes
		newItems := []v1.ProbeOutputItem{
			{
				Type: v1.PROBE,
				Source: v1.ProbeEndpointInfo{
					Type: v1.POD,
				},
				Destination: v1.ProbeEndpointInfo{
					Type: v1.POD,
				},
				Protocol:        "UDP",
				Port:            "53",
				ResultingAction: v1.DENY,
			},
		}

		err = sm.AppendProbes(&newItems)
		assert.NoError(t, err)

		// Check that both items are present
		state := sm.GetProbeState()
		assert.Len(t, state.Items, 2)
	})

	t.Run("Test AppendProbes with nil input", func(t *testing.T) {
		sm := NewStateManager()
		err := sm.AppendProbes(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "items cannot be nil")
	})

	t.Run("Test AppendErrors", func(t *testing.T) {
		sm := NewStateManager()

		// Add initial error
		initialErrors := []v1.ProbeOutputError{
			{
				Value: v1.ProbeOutputItem{
					Type: v1.PROBE,
					Source: v1.ProbeEndpointInfo{
						Type: v1.POD,
					},
					Destination: v1.ProbeEndpointInfo{
						Type: v1.POD,
					},
					Protocol:        "TCP",
					Port:            "80",
					ResultingAction: v1.ALLOW,
				},
				Reason: "initial error",
			},
		}

		err := sm.AppendErrors(&initialErrors)
		assert.NoError(t, err)

		// Add more errors
		newErrors := []v1.ProbeOutputError{
			{
				Value: v1.ProbeOutputItem{
					Type: v1.PROBE,
					Source: v1.ProbeEndpointInfo{
						Type: v1.POD,
					},
					Destination: v1.ProbeEndpointInfo{
						Type: v1.POD,
					},
					Protocol:        "UDP",
					Port:            "53",
					ResultingAction: v1.DENY,
				},
				Reason: "new error",
			},
		}

		err = sm.AppendErrors(&newErrors)
		assert.NoError(t, err)

		// Check that both errors are present
		state := sm.GetProbeState()
		assert.Len(t, state.Errors, 2)
	})

	t.Run("Test AppendErrors with nil input", func(t *testing.T) {
		sm := NewStateManager()
		err := sm.AppendErrors(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "items cannot be nil")
	})
}

func TestStateManagerNetInfoOperations(t *testing.T) {
	t.Run("Test AppendNetInfoV2", func(t *testing.T) {
		sm := NewStateManager()

		// Add initial net info
		key := "test-key"
		initialItems := []v1.PodNetworkingItem{
			{
				IP:       "192.168.1.2",
				Port:     "80",
				Protocol: "TCP",
			},
		}

		err := sm.AppendNetInfoV2(key, &initialItems)
		assert.NoError(t, err)

		// Add more items
		newItems := []v1.PodNetworkingItem{
			{
				IP:       "192.168.1.1",
				Port:     "443",
				Protocol: "TCP",
			},
		}

		err = sm.AppendNetInfoV2(key, &newItems)
		assert.NoError(t, err)

		// Check that both items are present
		state := sm.GetProbeState()
		assert.Len(t, state.PodNetworkingV2[key], 2)
	})

	t.Run("Test AppendNetInfoV2 with nil input", func(t *testing.T) {
		sm := NewStateManager()
		err := sm.AppendNetInfoV2("test-key", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "items cannot be nil")
	})

	t.Run("Test AppendNetInfoV2 with empty key", func(t *testing.T) {
		sm := NewStateManager()
		err := sm.AppendNetInfoV2("", nil)
		assert.Error(t, err)
	})

	t.Run("Test SetConfig", func(t *testing.T) {
		sm := NewStateManager()

		// Set config
		podName := "test-pod"
		items := []v1.PodNetworkingItem{
			{
				IP:       "192.168.1.5",
				Port:     "80",
				Protocol: "TCP",
			},
		}

		err := sm.SetConfig(podName, &items)
		assert.NoError(t, err)

		// Check that config is set
		state := sm.GetProbeState()
		assert.Len(t, state.PodConfigurationNetworking[podName], 1)
	})

	t.Run("Test SetConfig with nil input", func(t *testing.T) {
		sm := NewStateManager()
		err := sm.SetConfig("test-pod", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "items cannot be nil")
	})

	t.Run("Test SetConfig with empty pod name", func(t *testing.T) {
		sm := NewStateManager()
		err := sm.SetConfig("", nil)
		assert.Error(t, err)
	})

	t.Run("Test SetNetInfoV2", func(t *testing.T) {
		sm := NewStateManager()

		// Set net info
		key := "test-key"
		items := []v1.PodNetworkingItem{
			{
				IP:       "192.168.1.5",
				Port:     "80",
				Protocol: "TCP",
			},
		}

		err := sm.SetNetInfoV2(key, &items)
		assert.NoError(t, err)

		// Check that net info is set
		state := sm.GetProbeState()
		assert.Len(t, state.PodNetworkingV2[key], 1)
	})

	t.Run("Test SetNetInfoV2 with nil input", func(t *testing.T) {
		sm := NewStateManager()
		err := sm.SetNetInfoV2("test-key", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "items cannot be nil")
	})

	t.Run("Test SetNetInfoV2 with empty key", func(t *testing.T) {
		sm := NewStateManager()
		err := sm.SetNetInfoV2("", nil)
		assert.Error(t, err)
	})
}
