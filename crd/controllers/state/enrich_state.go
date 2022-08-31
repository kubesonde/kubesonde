package state

import (
	"strings"

	v12 "kubesonde.io/api/v1"
)

func EnrichState(state *v12.ProbeOutput, replicas []string, deployments []string) *v12.ProbeOutput {
	for idx, item := range state.Items {
		for _, replica := range replicas {
			if strings.Contains(item.Source.Name, replica) {
				state.Items[idx].Source.ReplicaSetName = replica

			}
			if strings.Contains(item.Destination.Name, replica) {
				state.Items[idx].Destination.ReplicaSetName = replica

			}
		}
		for _, deployment := range deployments {
			if strings.Contains(item.Source.Name, deployment) {
				state.Items[idx].Source.DeploymentName = deployment

			}
			if strings.Contains(item.Destination.Name, deployment) {
				state.Items[idx].Destination.DeploymentName = deployment

			}
		}
	}
	for idx, item := range state.Errors {
		for _, replica := range replicas {
			if strings.Contains(item.Value.Source.Name, replica) {
				state.Errors[idx].Value.Source.ReplicaSetName = replica

			}
			if strings.Contains(item.Value.Destination.Name, replica) {
				state.Errors[idx].Value.Destination.ReplicaSetName = replica

			}
		}
		for _, deployment := range deployments {
			if strings.Contains(item.Value.Source.Name, deployment) {
				state.Errors[idx].Value.Source.DeploymentName = deployment

			}
			if strings.Contains(item.Value.Destination.Name, deployment) {
				state.Errors[idx].Value.Destination.DeploymentName = deployment

			}
		}
	}
	return state
}
