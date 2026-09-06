// The dispatcher module is responsible for scheduling the probes
// it maintains an internal priority queue that continuously runs the probes
package dispatcher

import (
	"container/heap"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/semaphore"
	"k8s.io/client-go/kubernetes"
	"kubesonde.io/controllers/inner"
	"kubesonde.io/controllers/probe_command"
)

var dispatcherLog = log.New(log.Writer(), "dispatcher: ", log.LstdFlags)

type Priority int

const (
	LOW  Priority = 1
	HIGH Priority = 2
)

var (
	dispatcherSemaphore = semaphore.NewWeighted(1)
	pq                  = make(PriorityQueue, 0, 1000)
)

// defaultProbeWorkers is the number of probes executed concurrently. It is kept
// modest to stay well within the kubelet's exec/attach concurrency limits.
const defaultProbeWorkers = 10

// probeWorkers reads the desired worker-pool size from the environment,
// falling back to defaultProbeWorkers when unset or invalid.
func probeWorkers() int {
	if raw := os.Getenv("KUBESONDE_PROBE_WORKERS"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			return n
		}
		dispatcherLog.Printf("Invalid KUBESONDE_PROBE_WORKERS=%q, using default %d", raw, defaultProbeWorkers)
	}
	return defaultProbeWorkers
}

// executeProbe runs a single probe command. It is a package-level variable so
// that tests can substitute a fake executor to observe scheduling behaviour.
var executeProbe = func(apiClient kubernetes.Interface, command probe_command.KubesondeCommand) {
	inner.InspectAndStoreResult(apiClient, []probe_command.KubesondeCommand{command})
}

// Add probes to queue
func SendToQueue(commands []probe_command.KubesondeCommand, priority Priority) {
	if err := dispatcherSemaphore.Acquire(context.Background(), 1); err != nil {
		dispatcherLog.Printf("Failed to acquire semaphore: %v", err)
		return
	}
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

// popNext removes and returns the highest-priority probe from the queue.
// The boolean is false when the queue is empty. Access to the shared heap is
// serialised through dispatcherSemaphore, which also guards SendToQueue.
func popNext() (*Item, bool) {
	if err := dispatcherSemaphore.Acquire(context.Background(), 1); err != nil {
		dispatcherLog.Printf("Failed to acquire semaphore: %v", err)
		return nil, false
	}
	defer dispatcherSemaphore.Release(1)

	if pq.Len() == 0 {
		return nil, false
	}
	return heap.Pop(&pq).(*Item), true
}

// Run starts the probe worker pool. Probes are executed concurrently by a
// bounded number of workers (see KUBESONDE_PROBE_WORKERS) that all drain the
// shared priority queue. Workers run until ctx is cancelled.
//
// Run must be started only once for the shared queue; the caller is responsible
// for that (see the reconciler, which guards startup with sync.Once).
func Run(ctx context.Context, apiClient kubernetes.Interface) {
	// Initialise the heap under the same lock that guards all other queue
	// access, so it cannot race with a concurrent SendToQueue.
	if err := dispatcherSemaphore.Acquire(ctx, 1); err != nil {
		dispatcherLog.Printf("Failed to acquire semaphore: %v", err)
		return
	}
	heap.Init(&pq)
	dispatcherSemaphore.Release(1)

	workers := probeWorkers()
	dispatcherLog.Printf("Starting probe dispatcher with %d workers", workers)
	for i := 0; i < workers; i++ {
		go worker(ctx, apiClient)
	}
}

// worker continuously drains the queue and executes probes until ctx is
// cancelled. When the queue is empty it backs off briefly to avoid busy-spinning.
func worker(ctx context.Context, apiClient kubernetes.Interface) {
	const idleBackoff = 50 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		item, ok := popNext()
		if !ok {
			select {
			case <-ctx.Done():
				return
			case <-time.After(idleBackoff):
			}
			continue
		}
		executeProbe(apiClient, item.value)
	}
}
