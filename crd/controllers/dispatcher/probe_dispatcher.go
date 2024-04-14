// The dispatcher module is responsible for scheduling the probes
// it maintains an internal priority queue that continuously runs the probes
package dispatcher

import (
	"container/heap"
	"time"

	"golang.org/x/sync/semaphore"
	"k8s.io/client-go/kubernetes"
	"kubesonde.io/controllers/inner"
	"kubesonde.io/controllers/probe_command"
)

type Priority int

const (
	LOW  Priority = 1
	HIGH Priority = 2
)

var (
	dispatcherSemaphore = semaphore.NewWeighted(1)
	pq                  = make(PriorityQueue, 0, 1000)
)

// Add probes to queue
func SendToQueue(probes []probe_command.KubesondeCommand, priority Priority) {

	for result := dispatcherSemaphore.TryAcquire(1); !result; result = dispatcherSemaphore.TryAcquire(1) {
		// Keep trying to acquire
	}

	for index, probe := range probes {
		heap.Push(&pq, &Item{
			value:    probe,
			index:    index,
			priority: int(priority),
		})
	}
	dispatcherSemaphore.Release(1)

}
func QueueSize() int {
	for result := dispatcherSemaphore.TryAcquire(1); !result; result = dispatcherSemaphore.TryAcquire(1) {
		// Keep trying to acquire
	}
	size := pq.Len()
	dispatcherSemaphore.Release(1)
	return size
}

// Main routine. Starts the probe running loop.
func Run(apiClient kubernetes.Interface) {
	const probesPerSecond = 300 * time.Millisecond
	heap.Init(&pq)
	for { // FIXME: this could also be event based maybe
		result := dispatcherSemaphore.TryAcquire(1)
		if result == false {
			continue
		}
		for pq.Len() > 0 {
			start := time.Now()
			item := heap.Pop(&pq).(*Item)
			inner.InspectAndStoreResult(apiClient, []probe_command.KubesondeCommand{item.value})
			duration := time.Since(start)
			if duration < probesPerSecond {
				time.Sleep(probesPerSecond - duration)
			}

		}
		dispatcherSemaphore.Release(1)
	}
}
