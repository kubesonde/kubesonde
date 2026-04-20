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
	DestinationLabels    string               `json:"destinationLabels"`
	DestinationType      v1.ProbeEndpointType `json:"destinationType"`
	SourcePodName        string               `json:"source"`
	SourceIPAddress      string               `json:"sourceIPAddress"`
	SourceType           v1.ProbeEndpointType `json:"sourceType"`
	SourceLabels         string               `json:"sourceLabels"`
}

type ComparableKubesondeCommand struct {
	Command              string               `json:"command"`
	ContainerName        string               `json:"ContainerName"`
	Namespace            string               `json:"sourceNamespace"`
	Protocol             string               `json:"protocol"`
	Destination          string               `json:"destination"`
	DestinationIPAddress string               `json:"destinationIPAddress"`
	DestinationNamespace string               `json:"destinationNamespace"`
	DestinationPort      string               `json:"destinationPort"`
	DestinationLabels    string               `json:"destinationLabels"`
	DestinationType      v1.ProbeEndpointType `json:"destinationType"`
	SourcePodName        string               `json:"source"`
	SourceIPAddress      string               `json:"sourceIPAddress"`
	SourceType           v1.ProbeEndpointType `json:"sourceType"`
	SourceLabels         string               `json:"sourceLabels"`
}

func (item KubesondeCommand) ToComparableCommand() ComparableKubesondeCommand {
	return ComparableKubesondeCommand{
		Command:              item.Command,
		ContainerName:        item.ContainerName,
		Namespace:            item.Namespace,
		Protocol:             item.Protocol,
		Destination:          item.Destination,
		DestinationIPAddress: item.DestinationIPAddress,
		DestinationNamespace: item.DestinationNamespace,
		DestinationPort:      item.DestinationPort,
		DestinationLabels:    item.DestinationLabels,
		DestinationType:      item.DestinationType,
		SourcePodName:        item.SourcePodName,
		SourceIPAddress:      item.SourceIPAddress,
		SourceType:           item.SourceType,
		SourceLabels:         item.SourceLabels,
	}

}
