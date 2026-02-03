package probe_command

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	. "k8s.io/api/core/v1"
)

func TestPodsController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "inner controllers unit tests")
}

var _ = Describe("curlSucceded", func() {

	It("Should correctly recognize valid curl output", func() {
		Expect(CurlSucceded("200")).To(Equal(true))
	})

	It("Should correctly recognize invalid curl output", func() {
		Expect(CurlSucceded("000")).To(Equal(false))
	})

})

var _ = Describe("getAllPortsAndProtocolsFromPodSelector", func() {

	It("Returns all the ports in a pod", func() {

		Expect(getAllPortsAndProtocolsFromPodSelector(*podWithOpenPorts)).To(Equal([]PortAndProtocol{{port: 111, protocol: "TCP"}, {port: 112, protocol: "TCP"}, {port: 221, protocol: "TCP"}, {port: 222, protocol: "TCP"}}))
	})

	It("Returns empty array if pod has no open ports", func() {
		Expect(getAllPortsAndProtocolsFromPodSelector(podWithNoOpenPorts)).To(Equal([]PortAndProtocol{}))
	})
})

var _ = Describe("Build commands from spec", func() {
	It("Creates correct commands", func() {
		Skip("FIXME")

		expectedCommands := []KubesondeCommand{
			{
				Action:          "Allow",
				DestinationPort: "123",
				SourcePodName:   "test-src-pod",
				Destination:     "test-dest-pod",
				ContainerName:   "debugger",
				Namespace:       "test-namespace",
				Command:         "curl -s -o /dev/null -I -X GET -w %{http_code} test-dest-pod:123",
				ProbeChecker:    NmapSucceded,
			},
			{
				Destination:     "http://example.website.com",
				Action:          "Deny",
				SourcePodName:   "test-src-pod",
				DestinationPort: "80",
				ContainerName:   "debugger",
				Namespace:       "test-namespace",
				Command:         "curl -s -o /dev/null -I -X GET -w %{http_code} http://example.website.com",
				ProbeChecker:    NmapSucceded,
			},
			{
				Destination:     "http://example.website.com/api/healthz",
				Action:          "Deny",
				SourcePodName:   "test-src-pod",
				ContainerName:   "debugger",
				DestinationPort: "80",
				Namespace:       "test-namespace",
				Command:         "curl -s -o /dev/null -I -X GET -w %{http_code} http://example.website.com/api/healthz",
				ProbeChecker:    NmapSucceded,
			},
		}
		// FIXME
		result := BuildCommandsFromSpec(probingActions, "test-namespace")
		result_json := lo.Must1(json.Marshal(result))
		origin_json := lo.Must1(json.Marshal(expectedCommands))

		Expect(result_json).To(BeEquivalentTo(origin_json))
	})
})

var _ = Describe("Build commands from pod", func() {
	It("Creates empty commands", func() {
		Expect(BuildCommandsFromPodSelectors([]Pod{}, "")).To(BeNil())
	})
	It("Creates correct commands", func() {
		var ports = []int32{80, 443}
		container := buildContainers(ports)
		podA := buildTestPod([]Container{container}, "10.0.0.1")
		podB := buildTestPod([]Container{container}, "10.0.0.2")

		output := BuildCommandsFromPodSelectors([]Pod{podA, podB}, "")
		// 1. PodA -> PodB:80
		// 2. PodA -> PodB:443
		// 3. PodA -> Internet:443
		// 4. PodA -> Internet:80
		// 5  PodA -> DNS
		// 6. PodB -> PodA:80
		// 7. PodB -> PodA:443
		// 8. PodB -> Internet:443
		// 9. PodB -> Internet:80
		// 10. PodB -> DNS
		Expect(len(output)).To(Equal(16))
	})
})

var _ = Describe("Build targeted commands from pod", func() {
	It("Creates empty commands", func() {
		Expect(BuildCommandsFromPodSelectors([]Pod{}, "")).To(BeNil())
	})
	It("Creates correct commands", func() {
		var ports = []int32{80, 443}
		container := buildContainers(ports)
		targetContainer := buildContainers([]int32{8080})
		target := buildTestPod([]Container{targetContainer}, "10.0.0.1")

		available := []Pod{buildTestPod([]Container{container}, "10.0.0.2")}

		output := BuildTargetedCommands(target, available)
		/*
			target -> available 80
			target -> available 443
			available -> target 80
			Google DNS
			Google HTTP
			Google HTTPS
		*/
		Expect(len(output)).To(Equal(15))
	})
})
