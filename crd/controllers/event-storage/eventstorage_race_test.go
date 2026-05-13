package eventstorage

import (
	"sync"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubesonde.io/controllers/probe_command"
)

// TestRaceCommandsMap confirms that the commands map is properly protected
// by a mutex.
func TestRaceCommandsMap(t *testing.T) {
	// Clear commands before test using exported API
	for range 10 {
		cmd := probe_command.KubesondeCommand{
			SourcePodName: "cleanup-pod",
			Command:       "ping",
		}
		AddProbe(cmd)
	}
	for range 10 {
		cmd := probe_command.KubesondeCommand{
			SourcePodName: "cleanup-pod",
			Command:       "ping",
		}
		AddProbe(cmd)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			cmd := probe_command.KubesondeCommand{
				SourcePodName:        "pod-1",
				Command:              "ping",
				DestinationIPAddress: "10.0.0.1",
				DestinationPort:      "80",
				Protocol:             "TCP",
			}
			AddProbe(cmd)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			GetProbes()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			cmd := probe_command.KubesondeCommand{
				SourcePodName:        "pod-2",
				Command:              "ping",
				DestinationIPAddress: "10.0.0.2",
				DestinationPort:      "80",
				Protocol:             "TCP",
			}
			ProbeAvailable(cmd)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			cmd := probe_command.KubesondeCommand{
				SourcePodName:        "pod-3",
				Command:              "ping",
				DestinationIPAddress: "10.0.0.3",
				DestinationPort:      "80",
				Protocol:             "TCP",
			}
			AddProbe(cmd)
		}
	}()

	wg.Wait()
}

// TestRaceActivePods confirms that _activePods map is properly protected
// by a mutex. Run with: go test -race -run TestRaceActivePods
func TestRaceActivePods(t *testing.T) {
	// Clear active pods before test using exported API
	for range 10 {
		AddActivePod("cleanup-pod-1", CreatedPodRecord{})
	}
	for range 10 {
		DeleteActivePod("cleanup-pod-1")
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			AddActivePod("active-pod-1", CreatedPodRecord{
				Pod: v1.Pod{
					ObjectMeta: metav1.ObjectMeta{Name: "pod-1"},
				},
				DeploymentName:    "deploy-1",
				ReplicaSetName:    "rs-1",
				CreationTimestamp: 12345,
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			GetActivePods()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			GetActivePodNames()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			GetActivePodByName("active-pod-1")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			DeleteActivePod("active-pod-1")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			ClearEventStorage()
		}
	}()

	wg.Wait()
}

// TestRaceDeletedPods confirms that _deletedPods map is properly protected
// by a mutex.
func TestRaceDeletedPods(t *testing.T) {
	// Clear deleted pods before test using exported API
	for range 10 {
		AddDeletedPod("cleanup-del-1", DeletedPodRecord{})
	}
	for range 10 {
		DeleteActivePod("cleanup-del-1")
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			AddDeletedPod("del-pod-1", DeletedPodRecord{
				Pod:               v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "del-pod-1"}},
				DeploymentName:    "deploy-1",
				CreationTimestamp: 100,
				DeletionTimestamp: 200,
			})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			GetDeletedPodNames()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			ClearEventStorage()
		}
	}()

	wg.Wait()
}

// TestRaceServices confirms that _services slice is properly protected
// by a mutex.
func TestRaceServices(t *testing.T) {
	// Clear services before test using exported API
	for range 10 {
		AddService(v1.Service{})
	}
	for range 10 {
		ClearEventStorage()
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			AddService(v1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc-1"}})
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			GetServices()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 50 {
			ClearEventStorage()
		}
	}()

	wg.Wait()
}

// TestMixedReadWrite confirms that mixed concurrent read/write/delete
// on the SAME map is properly protected by a mutex.
func TestMixedReadWrite(t *testing.T) {
	// Clear state before test using exported API
	for range 10 {
		cmd := probe_command.KubesondeCommand{
			SourcePodName: "cleanup-pod",
			Command:       "ping",
		}
		AddProbe(cmd)
	}
	for range 10 {
		DeleteActivePod("cleanup-pod")
	}
	for range 10 {
		ClearEventStorage()
	}

	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmd := probe_command.KubesondeCommand{
				SourcePodName:        "mixed-pod",
				Command:              "ping",
				DestinationIPAddress: "10.0.0.1",
				DestinationPort:      "80",
				Protocol:             "TCP",
			}
			AddProbe(cmd)
			AddActivePod("pod-1", CreatedPodRecord{Pod: v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod-1"}}})
			DeleteActivePod("pod-1")
		}()
	}

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			GetProbes()
			GetActivePods()
			GetActivePodNames()
		}()
	}

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ClearEventStorage()
		}()
	}

	wg.Wait()
}
