package v1

type ProbeOutputItemType string

const (
	PROBE ProbeOutputItemType = "Probe"
	INFO  ProbeOutputItemType = "Information"
)

type ProbeEndpointType string

const (
	POD      ProbeEndpointType = "Pod"
	SERVICE  ProbeEndpointType = "Service"
	INTERNET ProbeEndpointType = "Internet"
)

type ComparableProbeOutputItem struct {
	Type ProbeOutputItemType `json:"type"`
	// ExpectedAction is the expected outcome of the probe. It might have values "allow" or "deny"
	ExpectedAction ActionType `json:"expectedAction,omitempty"`
	// ResultingAction is the resulted outcome of the probe. It might have values "allow" or "deny"
	ResultingAction ActionType `json:"resultingAction,omitempty"`
	// Source is a selector for the origin Pod or a set of pods
	Source ProbeEndpointInfo `json:"source,omitempty"`
	// Destination is a selector for the destination Pod or a set of pods
	Destination ProbeEndpointInfo `json:"destination,omitempty"`
	Protocol    string            `json:"protocol,omitempty"`
	// Port is the probing port for ToPodSelector defaults to 80
	Port          string `json:"port,omitempty"`
	ForwardedPort string `json:"forwardedPort,omitempty"`
}

func (item ProbeOutputItem) ToComparableProbe() ComparableProbeOutputItem {
	return ComparableProbeOutputItem{
		Type:            item.Type,
		ExpectedAction:  item.ExpectedAction,
		ResultingAction: item.ResultingAction,
		Source:          item.Source,
		Destination:     item.Destination,
		Protocol:        item.Protocol,
		Port:            item.Port,
		ForwardedPort:   item.ForwardedPort,
	}

}

type ProbeOutputItem struct {
	Type ProbeOutputItemType `json:"type"`
	// ExpectedAction is the expected outcome of the probe. It might have values "allow" or "deny"
	ExpectedAction ActionType `json:"expectedAction,omitempty"`
	// ResultingAction is the resulted outcome of the probe. It might have values "allow" or "deny"
	ResultingAction ActionType `json:"resultingAction,omitempty"`
	// Source is a selector for the origin Pod or a set of pods
	Source ProbeEndpointInfo `json:"source,omitempty"`
	// Destination is a selector for the destination Pod or a set of pods
	Destination          ProbeEndpointInfo `json:"destination,omitempty"`
	DestinationHostnames []string          `json:"destinationHostnames,omitempty"`
	Protocol             string            `json:"protocol,omitempty"`
	// Port is the probing port for ToPodSelector defaults to 80
	Port          string `json:"port,omitempty"`
	ForwardedPort string `json:"forwardedPort,omitempty"`
	Timestamp     int64  `json:"timestamp,omitempty"`
	// DebugOutput returns the http code of the request (assuming is TCP)
	DebugOutput string `json:"debugOutput,omitempty"`
}

type ProbeEndpointInfo struct {
	Type      ProbeEndpointType `json:"type,omitempty"`
	IPAddress string            `json:"IPAddress,omitempty"`
	// Name is the name of the endpoint
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	// DeploymentName is the name of the deployment belonging to the source Pod
	// +optional
	DeploymentName string `json:"deploymentName,omitempty"`
	// ReplicaSetName is the protocol to use when probing ToPodSelector
	// +optional
	ReplicaSetName string `json:"replicaSetName,omitempty"`
}

type ProbeOutputError struct {
	Value  ProbeOutputItem `json:"value,omitempty"`
	Reason string          `json:"reason,omitempty"`
}

type PodNetworkingInfo struct {
	PodName string `json:"podName"`
	Netstat string `json:"netstat"`
}
type ProbeOutput struct {
	Items                      []ProbeOutputItem   `json:"items"`
	Errors                     []ProbeOutputError  `json:"errors"`
	PodNetworking              []PodNetworkingInfo `json:"podNetworking"` // TODO: Delete this
	PodNetworkingV2            PodNetworkingInfoV2 `json:"podNetworkingv2"`
	PodConfigurationNetworking PodNetworkingInfoV2 `json:"podConfigurationNetworking"`
	Start                      string              `json:"start,omitempty"` // TODO: This may be unuseful
	End                        string              `json:"end,omitempty"`   // TODO: This may be unuseful
}

type PodNetworkingItem struct {
	Port     string `json:"port"`
	IP       string `json:"ip"`
	Protocol string `json:"protocol"`
}

type PodNetworkingInfoV2 map[string][]PodNetworkingItem
