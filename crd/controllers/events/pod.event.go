package events

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"golang.org/x/sync/semaphore"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v12 "kubesonde.io/api/v1"
	debugcontainer "kubesonde.io/controllers/debug-container"
	"kubesonde.io/controllers/dispatcher"
	. "kubesonde.io/controllers/event-storage"
	"kubesonde.io/controllers/probe_command"
	"kubesonde.io/controllers/state"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

/*
*
Boundaries of this function are:
  - k8s APIs
  - Event Storage
*/
var probe_processing_semaphore = semaphore.NewWeighted(1)
var log = logf.Log.WithName("kubesonde.podEvents")
var notDebuggablePods = []string{"kube-apiserver-minikube", "etcd-minikube", "kube-scheduler-minikube", "kube-controller-manager-minikube"}

func IsProcessingEvent() bool {
	var result = probe_processing_semaphore.TryAcquire(1)
	if result {
		probe_processing_semaphore.Release(1)
	}
	return !result
}

// TODO: Add resilience mechanism to to unlock the resource when pending for too long.
func addPodEvent(client kubernetes.Interface, pod v1.Pod) {
	pods := v1.PodList{
		Items: []v1.Pod{pod},
	}
	_, isnotDebuggable := lo.Find(notDebuggablePods, func(s string) bool {
		return s == pod.Name
	})
	if !isnotDebuggable {
		debugcontainer.InstallEphameralContainers(client, &pods)
	} else {
		log.Info(fmt.Sprintf("Non debuggable pod: %s", pod.Name))
		log.Info(fmt.Sprintf("IP: %s", pod.Status.PodIP))
		log.Info(fmt.Sprintf("Containers: %v", pod.Spec.Containers))

	}

	var timestamp = time.Now().Unix()
	var deployment string
	// TODO: this is a hacky thing.
	if strings.ContainsAny(pod.Name, "-") {
		deployment = strings.Split(pod.Name, "-")[0]
	} else {
		deployment = ""
	}
	/**
	If the active pods list is not empty then build probes
	*/
	var activePods = GetActivePods()
	if len(activePods) > 0 {
		// Build probes
		probes := probe_command.BuildTargetedCommands(pod, activePods)
		probes_from_pods := probe_command.BuildCommandsFromPodSelectors(activePods, "nothing")
		// Current pod probes all services

		AddProbes(probes)
		AddProbes(probes_from_pods)
		dispatcher.SendToQueue(probes, dispatcher.HIGH)
	}
	// TODO: Maybe there should be an event listener on the services to do the same thing.
	curr_services, err := client.CoreV1().Services(pod.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Info("Cannot get services")
	}
	services_probes := probe_command.BuildCommandsToServices(pod, curr_services.Items)
	AddProbes(services_probes)
	other_probes := probe_command.BuildCommandsToOutsideWorld(pod)
	AddProbes(other_probes)

	AddActivePod(pod.Name, CreatedPodRecord{
		Pod:               pod,
		CreationTimestamp: timestamp,
		DeploymentName:    deployment,
	})

	addPodPortsToState(pod)
}
func PodWithEphemeralContainer(client *kubernetes.Clientset, pod v1.Pod) bool {
	ppd, _ := client.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
	return EphemeralContainersRunning(ppd.Spec.EphemeralContainers, pod.Status.EphemeralContainerStatuses)
}
func addPodPortsToState(pod v1.Pod) {
	var name = pod.Name
	var ip = "0.0.0.0"
	var items []v12.PodNetworkingItem
	// Regular containers
	for _, container := range pod.Spec.Containers {
		portMapping := lo.Map(container.Ports, func(port v1.ContainerPort, in int) v12.PodNetworkingItem {
			return v12.PodNetworkingItem{
				Port:     strconv.Itoa(int(port.ContainerPort)),
				IP:       ip, // Unfortunately I do not know if the ip is exposed only on the cluster interface
				Protocol: string(port.Protocol),
			}
		})
		items = append(items, portMapping...)
	}
	// Init containers
	for _, container := range pod.Spec.InitContainers {
		portMapping := lo.Map(container.Ports, func(port v1.ContainerPort, in int) v12.PodNetworkingItem {
			return v12.PodNetworkingItem{
				Port:     strconv.Itoa(int(port.ContainerPort)),
				IP:       ip,
				Protocol: string(port.Protocol),
			}
		})
		// log.V(1).Info(fmt.Sprintf("InitContainers for pod %s: \n\n %v", pod.Name, portMapping))
		items = append(items, portMapping...)
	}
	state.SetConfig(name, &items)
}
func EphemeralContainersRunning(ephc []v1.EphemeralContainer, ephStats []v1.ContainerStatus) bool {
	if len(ephc) <= 1 || ephStats == nil {
		return false
	}
	if len(ephStats) <= 1 {
		return false
	}
	return true
}
func WaitEphemeralContainersToBeRunning(client *kubernetes.Clientset, pod v1.Pod) bool {
	failures := 0
polleph:
	for {
		if failures >= 10 { // More than 20 seconds without containers being ready
			log.Info(fmt.Sprintf("Pod %s seems unresponsive", pod.Name))
			return false
		}
		time.Sleep(2 * time.Second)
		ppd, _ := client.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
		ephc := ppd.Spec.EphemeralContainers
		ephStats := ppd.Status.EphemeralContainerStatuses

		if len(ephc) <= 1 || ephStats == nil {
			failures++
			continue
		}
		if len(ephStats) <= 1 {
			failures++
			continue
		}

		for _, itm := range ephStats {

			if itm.State.Running == nil {
				failures++
				goto polleph
			}
		}
		return true
	}
}
func WaitContainersToBeRunning(client *kubernetes.Clientset, pod v1.Pod) {
polleph:
	for {
		ppd, _ := client.CoreV1().Pods(pod.Namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
		ephStats := ppd.Status.ContainerStatuses

		for _, itm := range ephStats {

			if itm.State.Running == nil {
				goto polleph
			}
		}
		break
	}
}

func deletePodEvent(pod v1.Pod) {
	var deleteTimestamp = time.Now().Unix()
	var activePod = GetActivePodByName(pod.Name)
	AddDeletedPod(pod.Name, DeletedPodRecord{
		Pod:               pod,
		DeploymentName:    activePod.DeploymentName,
		CreationTimestamp: activePod.CreationTimestamp,
		DeletionTimestamp: deleteTimestamp,
	})
	DeleteActivePod(pod.Name)
	state.DeleteNetstatPod(pod.Name)
	errors := []v12.ProbeOutputError{{
		Value: v12.ProbeOutputItem{
			Timestamp: time.Now().Unix(),
			Source: v12.ProbeEndpointInfo{
				Name:      pod.Name,
				Namespace: pod.Namespace,
				IPAddress: pod.Status.PodIP,
			},
		},
		Reason: "PodDeleted",
	}}
	state.AppendErrors(&errors)
	log.Info(fmt.Sprintf("Pod deleted %s", pod.Name))
}

func getActivePodsEvent() []string {
	return GetActivePodNames()
}
