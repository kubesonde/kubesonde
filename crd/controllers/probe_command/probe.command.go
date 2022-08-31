package probe_command

import v1 "kubesonde.io/api/v1"

type KubesondeCommand struct {
	Action               v1.ActionType
	Command              string               `json:"command"`
	ContainerName        string               `json:"ContainerName"`
	ProbeChecker         func(string) bool    `json:"checker"`
	Namespace            string               `json:"sourceNamespace"`
	Protocol             string               `json:"protocol"`
	Destination          string               `json:"destination"`
	DestinationHostnames []string             `json:"destinationHostnames"`
	DestinationIPAddress string               `json:"destinationIPAddress"`
	DestinationNamespace string               `json:"destinationNamespace"`
	DestinationPort      string               `json:"destinationPort"`
	DestinationType      v1.ProbeEndpointType `json:"destinationType"`
	SourcePodName        string               `json:"source"`
	SourceIPAddress      string               `json:"sourceIPAddress"`
	SourceType           v1.ProbeEndpointType `json:"sourceType"`
}
