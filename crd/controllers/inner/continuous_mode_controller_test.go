package inner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	v1 "kubesonde.io/api/v1"
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
