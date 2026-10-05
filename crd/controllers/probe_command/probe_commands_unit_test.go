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

var _ = Describe("NmapTCPSucceeded", func() {
	It("Allows a confirmed open TCP port", func() {
		output := "Nmap scan report for 10.0.0.1\n80/tcp open  http\n1 IP address (1 host up) scanned"
		Expect(NmapTCPSucceeded(output)).To(BeTrue())
	})

	It("Denies when the port isn't reported open", func() {
		output := "Nmap scan report for 10.0.0.1\n1 IP address (1 host up) scanned"
		Expect(NmapTCPSucceeded(output)).To(BeFalse())
	})

	It("Denies when the host never responded", func() {
		output := "0 hosts up"
		Expect(NmapTCPSucceeded(output)).To(BeFalse())
	})
})

var _ = Describe("NmapUDPSucceeded", func() {
	It("Allows a confirmed open UDP port", func() {
		output := "Nmap scan report for 10.0.0.1\n53/udp open  domain\n1 IP address (1 host up) scanned"
		Expect(NmapUDPSucceeded(output)).To(BeTrue())
	})

	It("Denies open|filtered: a non-response is indistinguishable from a block", func() {
		output := "Nmap scan report for 10.0.0.1\n53/udp open|filtered  domain\n1 IP address (1 host up) scanned"
		Expect(NmapUDPSucceeded(output)).To(BeFalse())
	})

	It("Denies when the port is closed", func() {
		output := "Nmap scan report for 10.0.0.1\n1 IP address (1 host up) scanned"
		Expect(NmapUDPSucceeded(output)).To(BeFalse())
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
		Expect(BuildCommandsFromPodSelectors([]Pod{}, "", false, false)).To(BeNil())
	})
	It("Creates correct commands", func() {
		var ports = []int32{80, 443}
		container := buildContainers(ports)
		podA := buildTestPod([]Container{container}, "10.0.0.1")
		podB := buildTestPod([]Container{container}, "10.0.0.2")

		output := BuildCommandsFromPodSelectors([]Pod{podA, podB}, "", false, false)
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
	It("Drops the 2 outside-world probes per pod when both flags are disabled", func() {
		var ports = []int32{80, 443}
		container := buildContainers(ports)
		podA := buildTestPod([]Container{container}, "10.0.0.1")
		podB := buildTestPod([]Container{container}, "10.0.0.2")

		output := BuildCommandsFromPodSelectors([]Pod{podA, podB}, "", true, true)
		// Only the 4 pod-to-pod commands remain (2 per pod), no outside-world probes.
		Expect(len(output)).To(Equal(4))
	})
})

var _ = Describe("BuildCommandsToOutsideWorld", func() {
	target := buildTestPod([]Container{buildContainers([]int32{80})}, "10.0.0.1")

	It("Creates both Internet and kube-dns probes when nothing is disabled", func() {
		output := BuildCommandsToOutsideWorld(target, false, false)
		Expect(len(output)).To(Equal(6)) // 4 Google (DNS TCP/UDP, HTTP, HTTPS) + 2 kube-dns (TCP/UDP)
	})

	It("Skips Google probes when Internet probing is disabled", func() {
		output := BuildCommandsToOutsideWorld(target, true, false)
		Expect(len(output)).To(Equal(2))
		for _, cmd := range output {
			Expect(cmd.Destination).To(Equal("KUBE DNS"))
		}
	})

	It("Skips kube-dns probes when Service probing is disabled", func() {
		output := BuildCommandsToOutsideWorld(target, false, true)
		Expect(len(output)).To(Equal(4))
		for _, cmd := range output {
			Expect(cmd.Destination).ToNot(Equal("KUBE DNS"))
		}
	})

	It("Creates no probes when both Internet and Service probing are disabled", func() {
		output := BuildCommandsToOutsideWorld(target, true, true)
		Expect(output).To(BeNil())
	})

	It("Checks the KUBE DNS TCP probe with NmapSucceded, since nmap output never contains \"Server:\"", func() {
		output := BuildCommandsToOutsideWorld(target, false, false)
		kubeDNSTCP, found := lo.Find(output, func(cmd KubesondeCommand) bool {
			return cmd.Destination == "KUBE DNS" && cmd.Protocol == "TCP"
		})
		Expect(found).To(BeTrue())

		nmapOpenOutput := "Nmap scan report for kube-dns.kube-system.svc.cluster.local\nHost is up.\n53/tcp open  domain\n1 IP address (1 host up) scanned"
		Expect(kubeDNSTCP.ProbeChecker(nmapOpenOutput)).To(BeTrue())
	})
})

var _ = Describe("Build targeted commands from pod", func() {
	It("Creates empty commands", func() {
		Expect(BuildCommandsFromPodSelectors([]Pod{}, "", false, false)).To(BeNil())
	})
	It("Creates correct commands", func() {
		var ports = []int32{80, 443}
		container := buildContainers(ports)
		targetContainer := buildContainers([]int32{8080})
		target := buildTestPod([]Container{targetContainer}, "10.0.0.1")

		available := []Pod{buildTestPod([]Container{container}, "10.0.0.2")}

		output := BuildTargetedCommands(target, available, false, false)
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
	It("Drops outside-world probes when both flags are disabled", func() {
		var ports = []int32{80, 443}
		container := buildContainers(ports)
		targetContainer := buildContainers([]int32{8080})
		target := buildTestPod([]Container{targetContainer}, "10.0.0.1")

		available := []Pod{buildTestPod([]Container{container}, "10.0.0.2")}

		output := BuildTargetedCommands(target, available, true, true)
		// target -> available 80, target -> available 443, available -> target 8080
		Expect(len(output)).To(Equal(3))
	})
})
