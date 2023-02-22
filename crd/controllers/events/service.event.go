package events

import (
	"context"
	"fmt"
	"strconv"
	"time"

	. "kubesonde.io/controllers/event-storage"

	"golang.org/x/sync/semaphore"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	v12 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/inner"
	"kubesonde.io/controllers/probe_command"
)

var service_processing_semaphore = semaphore.NewWeighted(1)

// TODO: Not sure if this is correct
func getServicesAsProbes(servicePorts []v1.ServicePort, pods []v1.Pod, serviceClusterIP string, srvNamespace string, srvName string) []v12.ProbeOutputItem {
	var newItems []v12.ProbeOutputItem
	for _, servicePort := range servicePorts {

		// TODO: extract to function
		// TODO: we should also check the host port mapping
		// FIXME: this is a hack
		forwardedPort, dst := getDestinationAndPort(pods, servicePort, serviceClusterIP, srvNamespace, srvName)
		log.Info(fmt.Sprintf("Service %s, Port %s, TargetPort %d/%s",
			srvName, strconv.FormatInt(int64(servicePort.Port), 10), servicePort.TargetPort.IntValue(), servicePort.TargetPort.String()))
		newItems = append(newItems, v12.ProbeOutputItem{
			Type:            v12.INFO,
			ExpectedAction:  v12.DENY,
			ResultingAction: v12.ALLOW,
			Destination:     dst,
			Source: v12.ProbeEndpointInfo{
				Type: v12.INTERNET,
				Name: "Internet",
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
func getDestinationAndPort(pods []v1.Pod, servicePort v1.ServicePort, serviceClusterIP string, srvNamespace string, srvName string) (string, v12.ProbeEndpointInfo) {
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
						IPAddress: serviceClusterIP,
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
			IPAddress: serviceClusterIP,
		}
	}
	return forwardedPort, dst
}

func addServiceEvent(client *kubernetes.Clientset, service v1.Service) {
	log.Info(fmt.Sprintf("Acquire lock for service %s", service.Name))
	err := service_processing_semaphore.Acquire(context.Background(), 1)
	if err != nil {
		log.Error(err, "Could nod lock semaphore for service %s", service.Name)
	}
	/**
	If the active pods list is not empty then build probes
	*/
	if len(getActivePodsEvent()) > 0 {
		// Build probes
		probes := probe_command.BuildCommandFromService(GetActivePods(), service)
		inner.InspectAndStoreResult(client, probes)
	}

	// TODO: unlock the event storage here
	service_processing_semaphore.Release(1)
	log.Info(fmt.Sprintf("Release lock for service %s", service.Name))

}
