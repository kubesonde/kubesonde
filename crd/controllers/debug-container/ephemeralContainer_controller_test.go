package debug_container

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

func TestEphemeralContainers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ephemeral controller")
}

var _ = Describe("Checking for already existing containers", func() {
	It("Finds duplicate containers", func() {
		pod := v1.Pod{
			Spec: v1.PodSpec{
				EphemeralContainers: []v1.EphemeralContainer{{}, {}},
			},
		}
		Expect(ephemeralContainerExists(pod)).To(BeFalse())
	})
})
