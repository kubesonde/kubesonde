package debug_container

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
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

func TestGenerateDebugContainers(t *testing.T) {
	// GIVEN
	pod := v1.Pod{}
	privileged := true
	expected_pod := v1.Pod{
		Spec: v1.PodSpec{
			EphemeralContainers: []v1.EphemeralContainer{
				{
					EphemeralContainerCommon: v1.EphemeralContainerCommon{
						Name:                     "debugger",
						Image:                    "instrumentisto/nmap:latest",
						ImagePullPolicy:          v1.PullIfNotPresent,
						Stdin:                    true,
						TerminationMessagePolicy: v1.TerminationMessageReadFile,
						TTY:                      true,
						Command:                  []string{"sh"},
						SecurityContext: &v1.SecurityContext{
							Privileged: &privileged,
						},
					},
				},
				{
					EphemeralContainerCommon: v1.EphemeralContainerCommon{
						Name:                     "monitor",
						Image:                    "registry.cs.aalto.fi/kubesonde/monitor:latest",
						ImagePullPolicy:          v1.PullIfNotPresent,
						Stdin:                    true,
						TerminationMessagePolicy: v1.TerminationMessageReadFile,
						TTY:                      true,
						Command:                  []string{"sh"},
					},
				},
			},
		},
	}
	result := lo.Must1(generateDebugContainers(&pod))
	assert.Equal(t, &expected_pod, result)
}
