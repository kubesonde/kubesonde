package events

import (
	"fmt"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	v12 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/utils"
)

func getServicesAsProbes(service v1.Service, pods []v1.Pod) []v12.ProbeOutputItem {
	srvNamespace := service.Namespace
	srvName := service.Name
	servicePorts := service.Spec.Ports
	var newItems []v12.ProbeOutputItem
	for _, servicePort := range servicePorts {

		// TODO: extract to function
		// TODO: we should also check the host port mapping
		// FIXME: this is a hack
		forwardedPort, dst := getDestinationAndPort(service, pods, servicePort)
		log.Info(fmt.Sprintf("Service %s, Port %s, TargetPort %d/%s",
			srvName, strconv.FormatInt(int64(servicePort.Port), 10), servicePort.TargetPort.IntValue(), servicePort.TargetPort.String()))
		newItems = append(newItems, v12.ProbeOutputItem{
			Type:            v12.INFO,
			ExpectedAction:  v12.DENY,
			ResultingAction: v12.ALLOW,
			Destination:     dst,
			Source: v12.ProbeEndpointInfo{
				Type:      v12.INTERNET,
				Name:      "Internet",
				Namespace: srvNamespace,
				Labels:    utils.MapToString(service.Labels),
			},
			Protocol:      string(servicePort.Protocol),
			Port:          strconv.FormatInt(int64(servicePort.Port), 10),
			ForwardedPort: forwardedPort,
			Timestamp:     time.Now().Unix(),
		})
	}
	return newItems
}
func isSamePort(podPort v1.ContainerPort, servicePort v1.ServicePort) bool {
	if servicePort.TargetPort.StrVal == "" {
		return podPort.ContainerPort == servicePort.Port
	}
	if podPort.Name == servicePort.TargetPort.String() {
		return true
	}
	if podPort.ContainerPort == int32(servicePort.TargetPort.IntValue()) {
		return true
	}
	return false
}
func getDestinationAndPort(service v1.Service, pods []v1.Pod, servicePort v1.ServicePort) (string, v12.ProbeEndpointInfo) {
	srvNamespace := service.Namespace
	srvName := service.Name
	var dst v12.ProbeEndpointInfo
	var forwardedPort string
	var found = -1
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				// TODO: Maybe we use labels
				if isSamePort(port, servicePort) {
					forwardedPort = strconv.FormatInt(int64(port.ContainerPort), 10)
					found = 0
					dst = v12.ProbeEndpointInfo{
						Type:      v12.SERVICE,
						Name:      pod.Name,
						Namespace: pod.Namespace,
						IPAddress: service.Spec.ClusterIP,
						Labels:    utils.MapToString(service.Labels),
					}
				}
			}
		}

	}

	if found == -1 {
		dst = v12.ProbeEndpointInfo{
			Type:      v12.SERVICE,
			Name:      fmt.Sprintf("Unknown - %s", srvName),
			Namespace: srvNamespace,
			IPAddress: service.Spec.ClusterIP,
			Labels:    utils.MapToString(service.Labels),
		}
	}
	return forwardedPort, dst
}
