package inner

import (
	"context"
	"errors"
	"testing"
	"time"

	"bou.ke/monkey"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	v1 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/probe_command"
	"kubesonde.io/controllers/state"
)

func TestContinuousMode(t *testing.T) {
	RegisterFailHandler(Fail)
	//RunSpecs(t, "ContinuousMode")
}

var _ = Describe("ContinuousMode", func() {
	BeforeEach(func() {
		innerState := v1.ProbeOutput{
			Items:           []v1.ProbeOutputItem{},
			Errors:          []v1.ProbeOutputError{},
			PodNetworking:   []v1.PodNetworkingInfo{},
			PodNetworkingV2: make(v1.PodNetworkingInfoV2),
		}
		state.SetProbeState(&innerState)
	})
	It("Records errors", func() {
		wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
		patch := monkey.Patch(time.Now, func() time.Time { return wayback })
		defer patch.Unpatch()
		state := new(MockedCNIState)
		command := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationHostnames: nil,
			DestinationNamespace: "default",
			DestinationLabels:    "app=dest-pod;type=test",
			SourcePodName:        "test-pod",
			SourceLabels:         "app=source-pod;type=test",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}
		client := fake.NewSimpleClientset()
		p := &v12.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
			Status: v12.PodStatus{EphemeralContainerStatuses: []v12.ContainerStatus{
				v12.ContainerStatus{
					State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}},
				},
				v12.ContainerStatus{
					State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}},
				},
			}},
			Spec: v12.PodSpec{
				EphemeralContainers: []v12.EphemeralContainer{{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "debugger"}}, {EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "monitor"}}}}}
		_, err := client.CoreV1().Pods("default").Create(context.TODO(), p, metav1.CreateOptions{})
		if err != nil {
			log.Info("error injecting pod add: %v", err)
			Panic()
		}

		state.Mock.On("getClient").Return(client)
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, errors.New("this is an error"))
		output := InspectWithContinuousMode(state, []probe_command.KubesondeCommand{command})

		Expect(output.Errors).To(BeEquivalentTo(
			[]v1.ProbeOutputError{
				{
					Value: v1.ProbeOutputItem{
						Type:            v1.PROBE,
						ExpectedAction:  v1.ALLOW,
						ResultingAction: v1.DENY,
						Source: v1.ProbeEndpointInfo{
							Name:           "test-pod",
							Namespace:      "default",
							DeploymentName: "",
							ReplicaSetName: "",
							Labels:         "app=source-pod;type=test"},
						Destination: v1.ProbeEndpointInfo{
							Name:           "test-destination",
							Namespace:      "default",
							DeploymentName: "",
							ReplicaSetName: "",
							Labels:         "app=dest-pod;type=test"},
						Protocol:  "",
						Port:      "80",
						Timestamp: wayback.Unix(),
					},
					Reason: "this is an error",
				},
			},
		))
	})
	It("Records error and success", func() {
		state := new(MockedCNIState)
		wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
		patch := monkey.Patch(time.Now, func() time.Time { return wayback })
		defer patch.Unpatch()
		error_command := probe_command.KubesondeCommand{
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
		success_command := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "8080",
			DestinationHostnames: nil,
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, errors.New("this is an error")).Once()
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
		client := fake.NewSimpleClientset()
		p := &v12.Pod{ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
			Status: v12.PodStatus{EphemeralContainerStatuses: []v12.ContainerStatus{
				v12.ContainerStatus{
					State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}},
				},
				v12.ContainerStatus{
					State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}},
				},
			}},
			Spec: v12.PodSpec{
				EphemeralContainers: []v12.EphemeralContainer{{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "debugger"}}, {EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "monitor"}}}}}
		_, err := client.CoreV1().Pods("default").Create(context.TODO(), p, metav1.CreateOptions{})
		if err != nil {
			log.Info("error injecting pod add: %v", err)
			Panic()
		}

		state.Mock.On("getClient").Return(client)

		output := InspectWithContinuousMode(state, []probe_command.KubesondeCommand{error_command, success_command})

		Expect(output.Errors).To(BeEquivalentTo(
			[]v1.ProbeOutputError{
				{
					Value: v1.ProbeOutputItem{
						Type:            v1.PROBE,
						ExpectedAction:  v1.ALLOW,
						ResultingAction: v1.DENY,
						Source:          v1.ProbeEndpointInfo{Name: "test-pod", Namespace: "default", DeploymentName: "", ReplicaSetName: ""},
						Destination: v1.ProbeEndpointInfo{
							Name:           "test-destination",
							Namespace:      "default",
							DeploymentName: "",
							ReplicaSetName: "",
						},
						Protocol:  "",
						Port:      "80",
						Timestamp: wayback.Unix(),
					},
					Reason: "this is an error",
				},
			},
		))
		Expect(len(output.Items)).To(BeIdenticalTo(1))
	})
})
