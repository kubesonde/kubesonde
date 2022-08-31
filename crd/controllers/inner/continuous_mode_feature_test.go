package inner

import (
	"errors"
	"testing"
	"time"

	"bou.ke/monkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	v1 "kubesonde.io/api/v1"
	v12 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/probe_command"
	"kubesonde.io/controllers/state"
)

func TestContinuousMode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ContinuousMode")
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
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(false, errors.New("this is an error"))
		state.On("getClient").Times(1)

		output := InspectWithContinuousMode(state, []probe_command.KubesondeCommand{command})

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
		state.On("getClient").Times(1)

		output := InspectWithContinuousMode(state, []probe_command.KubesondeCommand{error_command, success_command})

		Expect(output.Errors).To(BeEquivalentTo(
			[]v1.ProbeOutputError{
				{
					Value: v1.ProbeOutputItem{
						Type:            v12.PROBE,
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
