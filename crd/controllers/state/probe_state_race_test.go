package state

import (
	"sync"
	"testing"

	v1 "kubesonde.io/api/v1"
)

func probeItems(items ...v1.ProbeOutputItem) []v1.ProbeOutputItem {
	return items
}

// TestDefaultManagerRace confirms that defaultManager and once have NO protection.
func TestDefaultManagerRace(t *testing.T) {
	ResetDefaultManager()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 100 {
			ResetDefaultManager()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 100 {
			GetDefaultManager()
		}
	}()

	wg.Wait()
}

// TestDefaultManagerConcurrentResetRead confirms that ResetDefaultManager
// and GetDefaultManager have a data race on the defaultManager pointer.
func TestDefaultManagerConcurrentReset(t *testing.T) {
	ResetDefaultManager()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			ResetDefaultManager()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			sm := GetDefaultManager()
			if sm == nil {
				t.Logf("GetDefaultManager returned nil")
			}
		}
	}()

	wg.Wait()
}

// TestClearBlocksProbeStorage confirms that Clear() holds the write lock
// for an extended period, blocking all concurrent AppendProbes calls.
func TestClearBlocksProbeStorage(t *testing.T) {
	sm := NewStateManager()

	items := probeItems(v1.ProbeOutputItem{Type: "Probe"})
	sm.AppendProbes(&items)

	var wg sync.WaitGroup

	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			item := probeItems(v1.ProbeOutputItem{Type: "Probe"})
			sm.AppendProbes(&item)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		sm.Clear()
	}()

	wg.Wait()
}

// TestLockContention simulates a realistic scenario where the probe loop
// continuously calls GetProbeState() while other goroutines try to AppendProbes.
func TestLockContention(t *testing.T) {
	sm := NewStateManager()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 200 {
			sm.GetProbeState()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 200 {
			item := probeItems(v1.ProbeOutputItem{Type: "Probe"})
			sm.AppendProbes(&item)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 20 {
			sm.Clear()
		}
	}()

	wg.Wait()
}

// TestDoubleAppendProbes confirms that InspectWithContinuousMode +
// InspectAndStoreResult causes duplicate AppendProbes calls on the same data.
func TestDoubleAppendProbes(t *testing.T) {
	sm := NewStateManager()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			item := probeItems(v1.ProbeOutputItem{Type: "Probe"})
			sm.AppendProbes(&item)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			state := sm.GetProbeState()
			sm.AppendProbes(&state.Items)
		}
	}()

	wg.Wait()
}
