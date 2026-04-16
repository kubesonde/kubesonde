package dispatcher

import (
	"container/heap"
	"sync"
	"testing"

	v1 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/probe_command"
)

// --- Helpers ---

func newItem(val probe_command.KubesondeCommand, priority int) *Item {
	return &Item{
		value:    val,
		priority: priority,
	}
}

// --- Tests ---

var (
	low_command = probe_command.KubesondeCommand{
		Destination:          "test-destination",
		DestinationPort:      "80",
		DestinationHostnames: nil,
		DestinationNamespace: "default",
		SourcePodName:        "test-pod",
		ContainerName:        "debugger",
		Namespace:            "default",
		Command:              "low",
		Action:               v1.ALLOW,
	}
	medium_command = probe_command.KubesondeCommand{
		Destination:          "test-destination",
		DestinationPort:      "80",
		DestinationHostnames: nil,
		DestinationNamespace: "default",
		SourcePodName:        "test-pod",
		ContainerName:        "debugger",
		Namespace:            "default",
		Command:              "medium",
		Action:               v1.ALLOW,
	}

	high_command = probe_command.KubesondeCommand{
		Destination:          "test-destination",
		DestinationPort:      "80",
		DestinationHostnames: nil,
		DestinationNamespace: "default",
		SourcePodName:        "test-pod",
		ContainerName:        "debugger",
		Namespace:            "default",
		Command:              "high",
		Action:               v1.ALLOW,
	}
)

func TestPriorityQueue_Order(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, newItem(low_command, 1))
	heap.Push(&pq, newItem(medium_command, 5))
	heap.Push(&pq, newItem(high_command, 10))

	item := heap.Pop(&pq).(*Item)
	if item.value.Command != "high" {
		t.Fatalf("expected highest priority first, got %v", item.value)
	}

	item = heap.Pop(&pq).(*Item)
	if item.value.Command != "medium" {
		t.Fatalf("expected second highest, got %v", item.value)
	}

	item = heap.Pop(&pq).(*Item)
	if item.value.Command != "low" {
		t.Fatalf("expected lowest last, got %v", item.value)
	}
}

func TestPriorityQueue_MultipleInsertions(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	for i := 0; i < 100; i++ {
		heap.Push(&pq, newItem(low_command, i))
	}

	prev := 1<<31 - 1 // max int
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		if item.priority > prev {
			t.Fatalf("heap property violated: %d came after %d", item.priority, prev)
		}
		prev = item.priority
	}
}

func TestPriorityQueue_Fix(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	itemA := newItem(low_command, 1)
	itemB := newItem(medium_command, 2)

	heap.Push(&pq, itemA)
	heap.Push(&pq, itemB)

	// Promote A to highest priority
	itemA.priority = 100
	heap.Fix(&pq, itemA.index)

	top := heap.Pop(&pq).(*Item)
	if top.value.Command != "low" {
		t.Fatalf("expected A after priority increase, got %v", top.value.Command)
	}
}

func TestPriorityQueue_IndexIntegrity(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	items := []*Item{
		newItem(low_command, 1),
		newItem(medium_command, 5),
		newItem(high_command, 3),
	}

	for _, it := range items {
		heap.Push(&pq, it)
	}

	for i, it := range pq {
		if it.index != i {
			t.Fatalf("index mismatch: expected %d, got %d", i, it.index)
		}
	}
}

func TestPriorityQueue_PushPopConsistency(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, newItem(low_command, 10))
	heap.Push(&pq, newItem(medium_command, 20))
	heap.Push(&pq, newItem(high_command, 15))

	if pq.Len() != 3 {
		t.Fatalf("expected length 3, got %d", pq.Len())
	}

	heap.Pop(&pq)
	heap.Pop(&pq)
	heap.Pop(&pq)

	if pq.Len() != 0 {
		t.Fatalf("expected empty queue, got %d", pq.Len())
	}
}

func TestPriorityQueue_ConcurrentAccess(t *testing.T) {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			mu.Lock()
			heap.Push(&pq, newItem(low_command, p))
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	prev := 1<<31 - 1
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		if item.priority > prev {
			t.Fatalf("concurrent heap corruption detected")
		}
		prev = item.priority
	}
}
