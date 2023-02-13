package debug_container

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

func TestEphemeralContainers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ephemeral controller")
}

var _ = Describe("ephemeralContainerExists", func() {
	It("Finds duplicate containers", func() {
		pod := v1.Pod{
			Spec: v1.PodSpec{
				EphemeralContainers: []v1.EphemeralContainer{{}, {}},
			},
		}
		Expect(ephemeralContainerExists(pod)).To(BeFalse())
	})
	It("Returns true if all containers have been installed", func() {
		pod := v1.Pod{
			Spec: v1.PodSpec{
				EphemeralContainers: []v1.EphemeralContainer{{
					EphemeralContainerCommon: v1.EphemeralContainerCommon{Name: "debugger"},
				}, {
					EphemeralContainerCommon: v1.EphemeralContainerCommon{Name: "monitor"},
				}},
			},
		}
		Expect(ephemeralContainerExists(pod)).To(BeTrue())
	})
})
