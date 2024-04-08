package debug_container

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestEphemeralContainers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ephemeral controller")
}

var _ = Describe("EphemeralContainerExists", func() {
	It("Finds duplicate containers", func() {
		pod := v1.Pod{
			Spec: v1.PodSpec{
				EphemeralContainers: []v1.EphemeralContainer{{}, {}},
			},
		}
		Expect(EphemeralContainerExists(&pod)).To(BeFalse())
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
		Expect(EphemeralContainerExists(&pod)).To(BeTrue())
	})
})

/*
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
							Image:                    "jackops93/kubesonde_monitor:latest",
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
*/
func TestInstallContainers(t *testing.T) {
	// Given
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	client := testclient.NewSimpleClientset()
	client.CoreV1().Pods("default").Create(context.TODO(), &pod, metav1.CreateOptions{})
	// When
	installContainers(client, &pod)
	// Then
	updatedPod, err := client.CoreV1().Pods("default").Get(context.TODO(), "test", metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(updatedPod.Spec.EphemeralContainers), updatedPod)
}
