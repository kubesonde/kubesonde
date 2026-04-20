// The dispatcher module is responsible for scheduling the probes
// it maintains an internal priority queue that continuously runs the probes
package dispatcher

import (
	"container/heap"
	"context"
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
func SendToQueue(commands []probe_command.KubesondeCommand, priority Priority) {
	dispatcherSemaphore.Acquire(context.Background(), 1)
	defer dispatcherSemaphore.Release(1)

	inQueue := make(map[probe_command.ComparableKubesondeCommand]bool, len(pq))
	for _, item := range pq {
		inQueue[item.value.ToComparableCommand()] = true
	}

	for _, command := range commands {
		if !inQueue[command.ToComparableCommand()] {
			heap.Push(&pq, &Item{
				value:    command,
				priority: int(priority),
			})
		}
	}
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
	const probeInterval = 50 * time.Millisecond
	heap.Init(&pq)
	for {
		dispatcherSemaphore.Acquire(context.Background(), 1)
		if pq.Len() == 0 {
			dispatcherSemaphore.Release(1)
			time.Sleep(probeInterval)
			continue
		}
		item := heap.Pop(&pq).(*Item)
		dispatcherSemaphore.Release(1)

		start := time.Now()
		inner.InspectAndStoreResult(apiClient, []probe_command.KubesondeCommand{item.value})
		duration := time.Since(start)
		if duration < probeInterval {
			time.Sleep(probeInterval - duration)
		}
	}
}
