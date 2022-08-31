package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func TestProbeStateController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils")
}

var _ = Describe("Pod filtering", func() {
	It("Filters pod by status", func() {
		podList := v1.PodList{
			Items: []v1.Pod{
				{
					Status: v1.PodStatus{
						Phase: "target",
					},
				},
				{
					Status: v1.PodStatus{
						Phase: "target",
					},
				},
				{
					Status: v1.PodStatus{
						Phase: "another",
					},
				},
			},
		}
		result := FilterPodsByStatus(&podList, "target")
		Expect(len(result.Items)).To(Equal(2))
	})
})
