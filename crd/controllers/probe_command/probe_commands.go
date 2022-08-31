package probe_command

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
	v12 "kubesonde.io/api/v1"
)

var (
	generateCurlCommand = func(action v12.ProbingAction) string {
		var curlParams string
		if action.ToPodSelector != "" {
			curlParams = fmt.Sprintf("%s:%s", action.ToPodSelector, action.Port)
		} else if action.Url != "" && action.Endpoint == "" {
			curlParams = action.Url
		} else if action.Url != "" && action.Endpoint != "" {
			curlParams = fmt.Sprintf("%s/%s", action.Url, action.Endpoint)
		}

		curlCMD := fmt.Sprintf(curlCommand, curlParams)
		return curlCMD
	}

	generateDestination = func(action v12.ProbingAction) string {
		var destination string
		if action.ToPodSelector != "" {
			destination = action.ToPodSelector
		} else if action.Url != "" && action.Endpoint == "" {
			destination = action.Url
		} else if action.Url != "" && action.Endpoint != "" {
			destination = fmt.Sprintf("%s/%s", action.Url, action.Endpoint)
		}
		return destination
	}

	generateDestinationPort = func(action v12.ProbingAction) string {
		var destinationPort string
		if action.Port != "" {
			destinationPort = action.Port
		} else if action.Url != "" {
			if strings.HasPrefix(action.Url, "https") {
				destinationPort = "443"
			} else {
				destinationPort = "80"
			}
		} else {
			destinationPort = "80"
		}
		return destinationPort
	}
)

/*
	This Command will make a GET request at the root of the ip:port combination
	and will fetch the return code. Here we suppose that every code <500 is a valid
	response. This is because we do not know if the root of the service is a valid address.
*/
const curlCommand = "curl -s -o /dev/null -I -X GET -w %%{http_code} %s"

// This command scans both UDP and TCP ports and returns only the amount of open ports
const nmapCommand = "nmap --open -sSU -p %d %s"
const nmapUDPCommand = "nmap --open -sU -p %d %s"
const nmapTCPCommand = "nmap --open -sT -p %d %s"
const nmapSCTPCommand = "nmap --open -sY -p %d %s"
const dnsUDPCommand = "nslookup %s %s"

func NslookupSucceded(output string) bool {
	return strings.Contains(output, "Server:")
}

func CurlSucceded(output string) bool {
	statusCode, err := strconv.Atoi(output)
	if err != nil {
		return false
	}
	return statusCode < 500 && statusCode > 0
}

func NmapSucceded(output string) bool {
	const successValue3 = "open|filtered"
	const successValue1 = "1 IP address (1 host up)"
	const successValue2 = "open"
	res := strings.Contains(output, successValue1) && (strings.Contains(output, successValue2) || strings.Contains(output, successValue3))
	return res
}

type PortAndProtocol struct {
	port     int32
	protocol string
}

func getAllPortsAndProtocolsFromPodSelector(pod v1.Pod) []PortAndProtocol {
	var portProto []PortAndProtocol
	// Regular containers
	for _, sourceContainer := range pod.Spec.Containers {
		for _, port := range sourceContainer.Ports {
			var protocol = string(port.Protocol)
			if protocol == "" {
				protocol = "TCP"
			}
			portProto = append(portProto, PortAndProtocol{
				port:     port.ContainerPort,
				protocol: protocol})
		}
	}
	// Init containers
	for _, sourceContainer := range pod.Spec.InitContainers {
		for _, port := range sourceContainer.Ports {
			var protocol = string(port.Protocol)
			if protocol == "" {
				protocol = "TCP"
			}
			portProto = append(portProto, PortAndProtocol{
				port:     port.ContainerPort,
				protocol: protocol})
		}
	}

	if len(portProto) == 0 {
		return []PortAndProtocol{}
	}
	return portProto
}

func generateNmapCommand(cmd string, ip string, port int32) string {
	command := fmt.Sprintf(cmd, port, ip)
	return command
}

func buildCommand(source v1.Pod, dest v1.Pod, port int32, protocol string, destType v12.ProbeEndpointType, srcType v12.ProbeEndpointType) KubesondeCommand {
	var namespace = source.Namespace
	addresses, err := net.LookupAddr(dest.Status.PodIP)
	if err != nil {
		addresses = []string{}
	}

	var cmd string
	if protocol == "TCP" {
		cmd = generateNmapCommand(nmapTCPCommand, dest.Status.PodIP, port)
	} else if protocol == "UDP" {
		cmd = generateNmapCommand(nmapUDPCommand, dest.Status.PodIP, port)
	} else if protocol == "SCTP" {
		cmd = generateNmapCommand(nmapSCTPCommand, dest.Status.PodIP, port)
	} else {
		cmd = generateNmapCommand(nmapCommand, dest.Status.PodIP, port)
	}
	return KubesondeCommand{
		Action:               v12.DENY,
		ContainerName:        "debugger",
		Namespace:            namespace,
		Command:              cmd,
		Protocol:             protocol,
		Destination:          dest.Name,
		DestinationPort:      strconv.Itoa(int(port)),
		DestinationHostnames: addresses,
		DestinationNamespace: dest.Namespace,
		DestinationIPAddress: dest.Status.PodIP,
		DestinationType:      destType,
		SourcePodName:        source.Name,
		SourceIPAddress:      source.Status.PodIP,
		SourceType:           srcType,
		ProbeChecker:         NmapSucceded,
	}
}

func buildCommandBase(source v1.Pod,
	dest string,
	destNamespace string,
	destIP string,
	destPort int32,
	protocol string,
	destType v12.ProbeEndpointType,
	srcType v12.ProbeEndpointType) KubesondeCommand {
	var namespace = source.Namespace
	addresses, err := net.LookupAddr(destIP)
	if err != nil {
		addresses = []string{}
	}

	var cmd string
	if protocol == "TCP" {
		cmd = generateNmapCommand(nmapTCPCommand, destIP, destPort)
	} else if protocol == "UDP" {
		cmd = generateNmapCommand(nmapUDPCommand, destIP, destPort)
	} else if protocol == "SCTP" {
		cmd = generateNmapCommand(nmapSCTPCommand, destIP, destPort)
	} else {
		cmd = generateNmapCommand(nmapCommand, destIP, destPort)
	}
	return KubesondeCommand{
		Action:               v12.DENY,
		ContainerName:        "debugger",
		Namespace:            namespace,
		Command:              cmd,
		Protocol:             protocol,
		Destination:          dest,
		DestinationPort:      strconv.Itoa(int(destPort)),
		DestinationHostnames: addresses,
		DestinationNamespace: destNamespace,
		DestinationIPAddress: destIP,
		DestinationType:      destType,
		SourcePodName:        source.Name,
		SourceIPAddress:      source.Status.PodIP,
		SourceType:           srcType,
		ProbeChecker:         NmapSucceded,
	}
}

// TODO: refactor-me!
func BuildCommandsFromSpec(actions []v12.ProbingAction, namespace string) []KubesondeCommand {
	var commands []KubesondeCommand
	for _, action := range actions {

		command := KubesondeCommand{
			Action:          action.Action,
			SourcePodName:   action.FromPodSelector,
			ContainerName:   "debugger",
			Namespace:       namespace,
			Command:         generateCurlCommand(action),
			Destination:     generateDestination(action),
			DestinationPort: generateDestinationPort(action),
			ProbeChecker:    NmapSucceded,
		}
		commands = append(commands, command)
	}

	return commands
}

func BuildCommandsFromPodSelectors(pods []v1.Pod, _ string) []KubesondeCommand {

	var commands []KubesondeCommand
	for _, source := range pods {
		for _, destination := range pods {
			if cmp.Equal(destination, source) == false {
				for _, portProto := range getAllPortsAndProtocolsFromPodSelector(destination) {
					commands = append(commands, buildCommand(source, destination, portProto.port, portProto.protocol, v12.POD, v12.POD))

				}

			}
		}
		other_commands := BuildCommandsToOutsideWorld(source)
		commands = append(commands, other_commands...)
	}

	return commands
}

func BuildCommandsToOutsideWorld(target v1.Pod) []KubesondeCommand {
	var commands []KubesondeCommand

	googleDNSTCP := KubesondeCommand{

		Action:               v12.DENY,
		SourcePodName:        target.Name,
		ContainerName:        "debugger",
		Namespace:            target.Namespace,
		Command:              fmt.Sprintf(nmapTCPCommand, 53, "8.8.8.8"),
		Destination:          "Google DNS",
		DestinationPort:      "53",
		DestinationIPAddress: "8.8.8.8",
		SourceIPAddress:      target.Status.PodIP,
		Protocol:             "TCP",
		SourceType:           v12.POD,
		DestinationType:      v12.INTERNET,
		ProbeChecker:         NmapSucceded,
	}
	googleDNSUDP := KubesondeCommand{

		Action:               v12.DENY,
		SourcePodName:        target.Name,
		ContainerName:        "debugger",
		Namespace:            target.Namespace,
		Command:              fmt.Sprintf(dnsUDPCommand, "google.com", "8.8.8.8"),
		Destination:          "Google DNS",
		DestinationPort:      "53",
		DestinationIPAddress: "8.8.8.8",
		SourceIPAddress:      target.Status.PodIP,
		Protocol:             "UDP",
		SourceType:           v12.POD,
		DestinationType:      v12.INTERNET,
		ProbeChecker:         NslookupSucceded,
	}
	googleHTTP := KubesondeCommand{

		Action:               v12.DENY,
		SourcePodName:        target.Name,
		ContainerName:        "debugger",
		Namespace:            target.Namespace,
		Command:              fmt.Sprintf(nmapTCPCommand, 80, "google.com"),
		Destination:          "Google",
		DestinationPort:      "80",
		DestinationIPAddress: "google.com",
		Protocol:             "TCP",
		SourceIPAddress:      target.Status.PodIP,
		SourceType:           v12.POD,
		DestinationType:      v12.INTERNET,
		ProbeChecker:         NmapSucceded,
	}
	googleHTTPS := KubesondeCommand{

		Action:               v12.DENY,
		SourcePodName:        target.Name,
		ContainerName:        "debugger",
		Namespace:            target.Namespace,
		Command:              fmt.Sprintf(nmapTCPCommand, 443, "google.com"),
		Destination:          "Google",
		DestinationPort:      "443",
		DestinationIPAddress: "google.com",
		SourceIPAddress:      target.Status.PodIP,
		Protocol:             "TCP",
		SourceType:           v12.POD,
		DestinationType:      v12.INTERNET,
		ProbeChecker:         NmapSucceded,
	}
	commands = append(commands /*cloudflareDNS,*/, googleDNSTCP, googleDNSUDP, googleHTTPS, googleHTTP)

	return commands
}

// Creates probe commends where target is the source of the probe and each available pod
// is the destination
func BuildTargetedCommands(target v1.Pod, availablePods []v1.Pod) []KubesondeCommand {
	var commands []KubesondeCommand

	targetPortsProto := getAllPortsAndProtocolsFromPodSelector(target)

	for _, source := range availablePods {
		for _, sourcePortProto := range getAllPortsAndProtocolsFromPodSelector(source) {
			commands = append(commands, buildCommand(target, source, sourcePortProto.port, sourcePortProto.protocol, v12.POD, v12.POD))
		}
		for _, targetPortProto := range targetPortsProto {
			commands = append(commands, buildCommand(source, target, targetPortProto.port, targetPortProto.protocol, v12.POD, v12.POD))
		}
		other_commands := BuildCommandsToOutsideWorld(source)
		commands = append(commands, other_commands...)
	}

	other_commands := BuildCommandsToOutsideWorld(target)
	commands = append(commands, other_commands...)

	return commands
}

// Creates probe commends where target is the target of the probe and each available pod
// is the source. The specified ports will be used as destination ports
func BuildTargetedCommandsToDestination(availablePods []v1.Pod, probeDestination v1.Pod, probeDestinationPorts []int32, protocol []string) []KubesondeCommand {
	var commands []KubesondeCommand

	for _, source := range availablePods {
		for idx := range probeDestinationPorts {
			commands = append(commands, buildCommand(source, probeDestination, probeDestinationPorts[idx], protocol[idx], v12.POD, v12.POD))

		}
	}

	return commands
}

// Creates probes from service if available
func BuildCommandFromService(availablePods []v1.Pod, probeDestination v1.Service) []KubesondeCommand {
	var commands []KubesondeCommand

	for _, source := range availablePods {
		for _, port := range probeDestination.Spec.Ports {
			commands = append(commands, buildCommandBase(source, probeDestination.Name, probeDestination.Namespace, probeDestination.Spec.ClusterIP, port.Port, "TCP", v12.SERVICE, v12.POD))
			commands = append(commands, buildCommandBase(source, probeDestination.Name, probeDestination.Namespace, probeDestination.Spec.ClusterIP, port.Port, "UDP", v12.SERVICE, v12.POD))
			commands = append(commands, buildCommandBase(source, probeDestination.Name, probeDestination.Namespace, probeDestination.Spec.ClusterIP, port.Port, "SCTP", v12.SERVICE, v12.POD))
		}
	}

	return commands
}
