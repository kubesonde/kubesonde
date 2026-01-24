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
	kubesondev1 "kubesonde.io/api/v1"
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
	// Create a kubesonde object to be used in the test
	kubesonde := &kubesondev1.Kubesonde{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Kubesonde",
			APIVersion: "kubesonde.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-kubesonde",
			Namespace: "default",
		},
		Spec: kubesondev1.KubesondeSpec{},
	}
	client.CoreV1().Pods("default").Create(context.TODO(), &pod, metav1.CreateOptions{})
	installContainers(client, *kubesonde, &pod)
	// Then
	updatedPod, err := client.CoreV1().Pods("default").Get(context.TODO(), "test", metav1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(updatedPod.Spec.EphemeralContainers), updatedPod)
}

func TestEphemeralContainerNamesAndImages(t *testing.T) {
	// Given
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	client := testclient.NewSimpleClientset()
	client.CoreV1().Pods("default").Create(context.TODO(), &pod, metav1.CreateOptions{})

	// Create a kubesonde object with specific images
	kubesonde := &kubesondev1.Kubesonde{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Kubesonde",
			APIVersion: "kubesonde.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-kubesonde",
			Namespace: "default",
		},
		Spec: kubesondev1.KubesondeSpec{
			DebuggerImage: "instrumentisto/nmap:latest",
			MonitorImage:  "jackops93/kubesonde_monitor:latest",
		},
	}

	// When
	installContainers(client, *kubesonde, &pod)

	// Then
	updatedPod, err := client.CoreV1().Pods("default").Get(context.TODO(), "test", metav1.GetOptions{})
	assert.Nil(t, err)

	// Check that we have exactly 2 ephemeral containers
	assert.Equal(t, 2, len(updatedPod.Spec.EphemeralContainers), "Should have exactly 2 ephemeral containers")

	// Check that the containers have the expected names and images
	containerMap := make(map[string]v1.EphemeralContainer)
	for _, container := range updatedPod.Spec.EphemeralContainers {
		containerMap[container.Name] = container
	}

	// Verify debugger container
	debuggerContainer, exists := containerMap["debugger"]
	assert.True(t, exists, "Debugger container should exist")
	assert.Equal(t, "instrumentisto/nmap:latest", debuggerContainer.Image, "Debugger container should have correct image")

	// Verify monitor container
	monitorContainer, exists := containerMap["monitor"]
	assert.True(t, exists, "Monitor container should exist")
	assert.Equal(t, "jackops93/kubesonde_monitor:latest", monitorContainer.Image, "Monitor container should have correct image")

	// Verify no other container names exist
	assert.Len(t, containerMap, 2, "Should have exactly 2 containers")
}
