package dispatcher

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	testclient "k8s.io/client-go/kubernetes/fake"
	v1 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/probe_command"
)

func TestContinuousMode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Probe dispatcher")
}

var _ = Describe("SendToQueue", func() {
	BeforeEach(func() {

	})
	It("Updates queue", func() {
		command := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationHostnames: nil,
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}
		commands := []probe_command.KubesondeCommand{command}

		// WHEN
		SendToQueue(commands, LOW)

		// THEN
		item, ok := popNext()
		Expect(ok).To(BeTrue())
		Expect(item.value).To(Equal(command))
		Expect(QueueSize()).To(Equal(0))
	})
})

// This test guards against the bottleneck described in finding #1: the
// dispatcher must not execute probes strictly one-at-a-time. With N independent
// probes that each take a fixed amount of time, the bounded worker pool runs
// them in parallel (observed max concurrency > 1). Against the previous serial
// implementation max concurrency was always 1 and this test failed.
var _ = Describe("Run executes probes concurrently", func() {
	It("runs independent probes in parallel", func() {
		const (
			numProbes     = 10
			probeDuration = 100 * time.Millisecond
		)

		var (
			inFlight    int32
			maxInFlight int32
		)

		var wg sync.WaitGroup
		wg.Add(numProbes)

		// Substitute a slow fake executor that records peak concurrency.
		original := executeProbe
		defer func() { executeProbe = original }()
		executeProbe = func(_ kubernetes.Interface, _ probe_command.KubesondeCommand) {
			defer wg.Done()
			cur := atomic.AddInt32(&inFlight, 1)
			for {
				old := atomic.LoadInt32(&maxInFlight)
				if cur <= old || atomic.CompareAndSwapInt32(&maxInFlight, old, cur) {
					break
				}
			}
			time.Sleep(probeDuration)
			atomic.AddInt32(&inFlight, -1)
		}

		// Enqueue N distinct probes.
		commands := make([]probe_command.KubesondeCommand, 0, numProbes)
		for i := 0; i < numProbes; i++ {
			commands = append(commands, probe_command.KubesondeCommand{
				Destination:     "test-destination",
				DestinationPort: fmt.Sprintf("%d", 8000+i),
				SourcePodName:   "test-pod",
				ContainerName:   "debugger",
				Namespace:       "default",
				Command:         "sample command",
				Action:          v1.ALLOW,
			})
		}

		client := testclient.NewSimpleClientset()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel() // stop the worker pool when the spec ends
		Run(ctx, client)
		SendToQueue(commands, LOW)

		// Wait for all probes to be executed (generously).
		done := make(chan struct{})
		go func() { wg.Wait(); close(done) }()
		Eventually(done, 10*time.Second).Should(BeClosed())

		Expect(atomic.LoadInt32(&maxInFlight)).To(BeNumerically(">", int32(1)),
			"probes ran serially (max concurrency == 1); dispatcher does not parallelise")
	})
})

/*
var _ = Describe("Runs", func() {
	It("Runs", func() {
		command := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationHostnames: nil,
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}
		commands := []probe_command.KubesondeCommand{command}
		client := testclient.NewSimpleClientset()

		// WHEN
		go Run(client)
		SendToQueue(commands, LOW)

		time.Sleep(2 * time.Second)
		// THEN
		Expect(pq.Len()).To(Equal(0))
	})
})*/
