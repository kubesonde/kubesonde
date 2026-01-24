package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"
	v1 "kubesonde.io/api/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	defaultLockTimeout = 5 * time.Second
)

var (
	log            = logf.Log.WithName("controllers.state")
	defaultManager *StateManager
	once           sync.Once
)

// StateManager handles concurrent access to probe state
type StateManager struct {
	mu                  sync.RWMutex
	probeOutput         v1.ProbeOutput
	podsWithNetstat     []string
	podsWithNetstatLock sync.RWMutex
	lockTimeout         time.Duration
}

// NewStateManager creates a new state manager instance
func NewStateManager() *StateManager {
	return &StateManager{
		probeOutput: v1.ProbeOutput{
			Items:                      []v1.ProbeOutputItem{},
			Errors:                     []v1.ProbeOutputError{},
			PodNetworking:              []v1.PodNetworkingInfo{},
			PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
			PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
		},
		podsWithNetstat: []string{},
		lockTimeout:     defaultLockTimeout,
	}
}

// GetDefaultManager returns the singleton state manager instance
func GetDefaultManager() *StateManager {
	once.Do(func() {
		defaultManager = NewStateManager()
	})
	return defaultManager
}

// SetNestatPod adds a pod to the netstat tracking list
func (sm *StateManager) SetNestatPod(pod string) {
	sm.podsWithNetstatLock.Lock()
	defer sm.podsWithNetstatLock.Unlock()
	sm.podsWithNetstat = append(sm.podsWithNetstat, pod)
}

// DeleteNetstatPod removes a pod from the netstat tracking list
func (sm *StateManager) DeleteNetstatPod(pod string) {
	sm.podsWithNetstatLock.Lock()
	defer sm.podsWithNetstatLock.Unlock()
	sm.podsWithNetstat = lo.Filter(sm.podsWithNetstat, func(s string, i int) bool {
		return s != pod
	})
}

// GetNetstatPods returns a copy of the netstat pod list
func (sm *StateManager) GetNetstatPods() []string {
	sm.podsWithNetstatLock.RLock()
	defer sm.podsWithNetstatLock.RUnlock()
	// Return a copy to prevent external modifications
	pods := make([]string, len(sm.podsWithNetstat))
	copy(pods, sm.podsWithNetstat)
	return pods
}

// SetProbeState replaces the entire probe state
func (sm *StateManager) SetProbeState(probes *v1.ProbeOutput) error {
	if probes == nil {
		return fmt.Errorf("probes cannot be nil")
	}

	log.Info(fmt.Sprintf("Setting probe state. Items: %d - Errors: %d",
		len(probes.Items), len(probes.Errors)))

	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.probeOutput = *probes
	return nil
}

// GetProbeState returns a copy of the current probe state
func (sm *StateManager) GetProbeState() v1.ProbeOutput {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// Deep copy to prevent external modifications
	return v1.ProbeOutput{
		Items:                      append([]v1.ProbeOutputItem{}, sm.probeOutput.Items...),
		Errors:                     append([]v1.ProbeOutputError{}, sm.probeOutput.Errors...),
		PodNetworking:              append([]v1.PodNetworkingInfo{}, sm.probeOutput.PodNetworking...),
		PodNetworkingV2:            copyNetworkingMapV2(sm.probeOutput.PodNetworkingV2),
		PodConfigurationNetworking: copyNetworkingMapV2(sm.probeOutput.PodConfigurationNetworking),
		Start:                      sm.probeOutput.Start,
		End:                        sm.probeOutput.End,
	}
}

// AppendProbes adds unique probe items to the state
func (sm *StateManager) AppendProbes(items *[]v1.ProbeOutputItem) error {
	if items == nil {
		return fmt.Errorf("items cannot be nil")
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	newItems := append(sm.probeOutput.Items, *items...)
	sm.probeOutput.Items = lo.UniqBy(newItems, func(poi v1.ProbeOutputItem) v1.ComparableProbeOutputItem {
		return poi.ToComparableProbe()
	})

	return nil
}

// AppendErrors adds unique error items to the state
func (sm *StateManager) AppendErrors(items *[]v1.ProbeOutputError) error {
	if items == nil {
		return fmt.Errorf("items cannot be nil")
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	newItems := append(sm.probeOutput.Errors, *items...)
	sm.probeOutput.Errors = lo.UniqBy(newItems, func(poe v1.ProbeOutputError) v1.ComparableProbeOutputItem {
		return poe.Value.ToComparableProbe()
	})

	return nil
}

// AppendNetInfoV2 adds networking items to a specific key (union operation)
func (sm *StateManager) AppendNetInfoV2(key string, items *[]v1.PodNetworkingItem) error {
	if items == nil {
		return fmt.Errorf("items cannot be nil")
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.probeOutput.PodNetworkingV2[key] = lo.Union(sm.probeOutput.PodNetworkingV2[key], *items)
	return nil
}

// SetConfig sets the configuration networking for a specific pod
func (sm *StateManager) SetConfig(podName string, items *[]v1.PodNetworkingItem) error {
	if items == nil {
		return fmt.Errorf("items cannot be nil")
	}
	if podName == "" {
		return fmt.Errorf("podName cannot be empty")
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.probeOutput.PodConfigurationNetworking[podName] = *items
	return nil
}

// SetNetInfoV2 replaces networking items for a specific key
func (sm *StateManager) SetNetInfoV2(key string, items *[]v1.PodNetworkingItem) error {
	if items == nil {
		return fmt.Errorf("items cannot be nil")
	}
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.probeOutput.PodNetworkingV2[key] = *items
	return nil
}

// Helper function to deep copy networking map
func copyNetworkingMapV2(src v1.PodNetworkingInfoV2) v1.PodNetworkingInfoV2 {
	if src == nil {
		return make(v1.PodNetworkingInfoV2)
	}

	dst := make(v1.PodNetworkingInfoV2, len(src))
	for k, v := range src {
		dst[k] = append([]v1.PodNetworkingItem{}, v...)
	}
	return dst
}

// Package-level functions for backward compatibility
func SetNestatPod(pod string) {
	GetDefaultManager().SetNestatPod(pod)
}

func DeleteNetstatPod(pod string) {
	GetDefaultManager().DeleteNetstatPod(pod)
}

func GetNetstatPods() []string {
	return GetDefaultManager().GetNetstatPods()
}

func SetProbeState(probes *v1.ProbeOutput) {
	if err := GetDefaultManager().SetProbeState(probes); err != nil {
		log.Error(err, "Failed to set probe state")
	}
}

func GetProbeState() v1.ProbeOutput {
	return GetDefaultManager().GetProbeState()
}

func AppendProbes(items *[]v1.ProbeOutputItem) {
	if err := GetDefaultManager().AppendProbes(items); err != nil {
		log.Error(err, "Failed to append probes")
	}
}

func AppendErrors(items *[]v1.ProbeOutputError) {
	if err := GetDefaultManager().AppendErrors(items); err != nil {
		log.Error(err, "Failed to append errors")
	}
}

func AppendNetInfoV2(key string, items *[]v1.PodNetworkingItem) {
	if err := GetDefaultManager().AppendNetInfoV2(key, items); err != nil {
		log.Error(err, "Failed to append net info v2")
	}
}

func SetConfig(podName string, items *[]v1.PodNetworkingItem) {
	if err := GetDefaultManager().SetConfig(podName, items); err != nil {
		log.Error(err, "Failed to set config")
	}
}

func SetNetInfoV2(key string, items *[]v1.PodNetworkingItem) {
	if err := GetDefaultManager().SetNetInfoV2(key, items); err != nil {
		log.Error(err, "Failed to set net info v2")
	}
}
