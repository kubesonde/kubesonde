package inner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	kubesondev1 "kubesonde.io/api/v1"
	v1 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/probe_command"
	"kubesonde.io/controllers/state"
)

func TestWithDeploymentInformationFast(t *testing.T) {

	innerState := v1.ProbeOutput{
		Items: []v1.ProbeOutputItem{
			{
				Type: v1.PROBE,
				Source: v1.ProbeEndpointInfo{
					Type:           v1.POD,
					Name:           "FirstPod-replica-pod",
					DeploymentName: "FirstPod",
					ReplicaSetName: "FirstPod-replica",
					Labels:         "key=val;k3y=v4l",
				},
				Destination: v1.ProbeEndpointInfo{
					Type:           v1.POD,
					Name:           "SecondPod-replica-pod",
					DeploymentName: "SecondPod",
					ReplicaSetName: "SecondPod-replica",
					Labels:         "anotherkey=val2",
				},
				Protocol:        "TCP",
				Port:            "80",
				ResultingAction: v1.ALLOW,
			},
		},
		Errors:          []v1.ProbeOutputError{},
		PodNetworking:   []v1.PodNetworkingInfo{},
		PodNetworkingV2: make(v1.PodNetworkingInfoV2),
	}
	state.SetProbeState(&innerState)

	output := v1.ProbeOutputItem{
		Type: v1.PROBE,
		Source: v1.ProbeEndpointInfo{
			Type:           v1.POD,
			Name:           "SecondPod-replica-pod",
			DeploymentName: "SecondPod",
			ReplicaSetName: "SecondPod-replica",
		},
		Destination: v1.ProbeEndpointInfo{
			Type:           v1.POD,
			Name:           "FirstPod-replica-pod",
			DeploymentName: "FirstPod",
			ReplicaSetName: "FirstPod-replica",
		},
		Protocol:        "TCP",
		Port:            "80",
		ResultingAction: v1.ALLOW,
	}
	client := fake.NewSimpleClientset()
	updated := withDeploymentInformation(client, output)

	assert.Equal(t, "SecondPod-replica-pod", updated.Source.Name)
	assert.Equal(t, "SecondPod", updated.Source.DeploymentName)
	assert.Equal(t, "SecondPod-replica", updated.Source.ReplicaSetName)
	assert.Equal(t, "FirstPod-replica-pod", updated.Destination.Name)
	assert.Equal(t, "FirstPod", updated.Destination.DeploymentName)
	assert.Equal(t, "FirstPod-replica", updated.Destination.ReplicaSetName)
	assert.Equal(t, "key=val;k3y=v4l", updated.Destination.Labels)
	assert.Equal(t, "anotherkey=val2", updated.Source.Labels)

}

func TestCommandToProbe(t *testing.T) {
	// Create a command with all fields
	command := probe_command.KubesondeCommand{
		Action:               kubesondev1.ALLOW,
		Command:              "test command",
		ContainerName:        "test-container",
		ProbeChecker:         func(string) bool { return true },
		Namespace:            "test-namespace",
		Protocol:             "TCP",
		Destination:          "test-destination",
		DestinationHostnames: []string{"host1", "host2"},
		DestinationIPAddress: "192.168.1.1",
		DestinationNamespace: "dest-namespace",
		DestinationPort:      "8080",
		DestinationLabels:    "label1=value1",
		DestinationType:      kubesondev1.SERVICE,
		DestinationSelector:  "app=web",
		SourcePodName:        "source-pod",
		SourceIPAddress:      "192.168.1.2",
		SourceType:           kubesondev1.POD,
		SourceLabels:         "label2=value2",
	}

	output := toProbeItem(command, kubesondev1.ALLOW)

	assert.Equal(t, kubesondev1.ALLOW, output.ExpectedAction)
	assert.Equal(t, "test-destination", output.Destination.Name)
	assert.Equal(t, "dest-namespace", output.Destination.Namespace)
	assert.Equal(t, "label1=value1", output.Destination.Labels)
	assert.Equal(t, "app=web", output.Destination.Selector)
	assert.Equal(t, "source-pod", output.Source.Name)
	assert.Equal(t, "test-namespace", output.Source.Namespace)
	assert.Equal(t, "label2=value2", output.Source.Labels)
}
