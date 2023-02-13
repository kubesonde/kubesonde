package dispatcher

import (
	"container/heap"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		result := heap.Pop(&pq).(*Item).value
		Expect(result).To(Equal(command))
		Expect(pq.Len()).To(Equal(0))
	})
})
