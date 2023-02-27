// The dispatcher module is responsible for scheduling the probes
// it maintains an internal priority queue that continuously runs the probes
package dispatcher

import (
	"container/heap"
	"context"
	"time"

	"github.com/samber/lo"
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
	lo.Must0(dispatcherSemaphore.Acquire(context.Background(), 1))
	for index, probe := range probes {
		heap.Push(&pq, &Item{
			value:    probe,
			index:    index,
			priority: int(priority),
		})
	}
	dispatcherSemaphore.Release(1)

}

// Main routine. Starts the probe running loop.
func Run(apiClient kubernetes.Interface) {
	const probesPerSecond = time.Second / 10
	heap.Init(&pq)
	for { // FIXME: this could also be event based maybe
		lo.Must0(dispatcherSemaphore.Acquire(context.Background(), 1))
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
