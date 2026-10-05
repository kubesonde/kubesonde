package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v12 "kubesonde.io/api/v1"
	debug_container "kubesonde.io/controllers/debug-container"
	kubesondeDispatcher "kubesonde.io/controllers/dispatcher"
	eventstorage "kubesonde.io/controllers/event-storage"
	"kubesonde.io/controllers/probe_command"
	"kubesonde.io/controllers/state"
	"kubesonde.io/rest_apis/types"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("Monitor controller")
var MAX_CONNECT_RETRIES = 6

func rebuildServiceProbesForPod(pod v1.Pod) {
	currServices := eventstorage.GetServices()
	serviceProbes := probe_command.BuildCommandsToServices(pod, currServices)

	kubesondeDispatcher.SendToQueue(serviceProbes, kubesondeDispatcher.HIGH)
}

// This function starts an infinite loop
func RunMonitorContainers(ctx context.Context, client kubernetes.Interface, kubesonde v12.Kubesonde) {
	for {
		select {
		case <-ctx.Done():
			log.Info("Monitor stopped")
			return
		default:
		}
		var pods = eventstorage.GetActivePods() /*lo.Filter(GetActivePods(), func(pod v1.Pod, i int) bool {
			return PodWithEphemeralContainer(client, pod)
		})*/
		var currPodsWithNetstat = state.GetNetstatPods()
		sort.Strings(currPodsWithNetstat)
		var currPodsWithNetstatLen = len(currPodsWithNetstat)
		lo.ForEach(pods, func(p v1.Pod, i int) {
			var index = sort.SearchStrings(currPodsWithNetstat, p.Name)
			if index < currPodsWithNetstatLen && currPodsWithNetstat[index] == p.Name {
				return
			}
			fresh_pod, erro := client.CoreV1().Pods(p.Namespace).Get(ctx, p.Name, metav1.GetOptions{})
			if erro != nil {
				log.Info(fmt.Sprintf("Pod %s does not exist, removing from state", p.Name))
				eventstorage.DeleteActivePod(p.Name)
				return
			}
			//	WaitEphemeralContainersToBeRunning(client, p)
			if !debug_container.EphemeralContainerExists(fresh_pod) || !debug_container.EphemeralContainersRunning(fresh_pod) {
				return
			}
			// log.Info(fmt.Sprintf("Running monitor on pod %s", p.Name))

			if !kubesonde.Spec.DisableServiceProbing {
				rebuildServiceProbesForPod(p)
			}

			var stdout, stderr, err = debug_container.RunMonitorContainerProcess(client, p.Namespace, p.Name)
			if err != nil {
				log.Error(err, "Could not run monitor process")
				log.Info(err.Error())
				return
			}
			go ProcessNetInfo(client, stdout, stderr, p.Name)
			state.SetNestatPod(p.Name)

		})
		select {
		case <-ctx.Done():
			log.Info("Monitor stopped")
			return
		case <-time.After(10 * time.Second):
		}
	}
}

func deleteNetstatPodWithLog(podname string, stderr *bytes.Buffer) {
	log.Info(fmt.Sprintf("Restarting monitor container on %s", podname))
	if len(stderr.String()) > 0 {
		log.Info(fmt.Sprintf("Stderr %s", stderr.String()))
	}
	state.DeleteNetstatPod(podname)
}

func eventuallyDecodeNetinfoData(stdout *bytes.Buffer) (types.NestatInfoRequestBody, error) {
	var payload_raw = stdout.String()
	var index = strings.Index(payload_raw, "\n")

	if index < 0 { // Not found
		return nil, errors.New("not found")
	}
	if index == 0 { // First char is \n
		stdout.Next(1)
		payload_raw = stdout.String()
		index = strings.Index(payload_raw, "\n")
		if index < 0 {
			return nil, errors.New("not found")
		}
	}

	var potential_json = stdout.Next(index)
	var payload types.NestatInfoRequestBody
	var err = json.Unmarshal(potential_json, &payload)
	if err != nil {
		log.Error(err, "Could not decode monitor")
		log.Info(string(payload_raw))
		return nil, errors.New("could not decode")
	}
	return payload, nil
}

func ProcessNetInfo(apiClient kubernetes.Interface, stdout *bytes.Buffer, stderr *bytes.Buffer, podname string) {
	var counter = 0
	for {
		if stdout.Len() == 0 {
			counter += 1
			if counter >= MAX_CONNECT_RETRIES {
				deleteNetstatPodWithLog(podname, stderr)
				counter = 0
				return
			}
			counterAsDuration := time.Duration(counter * 1000)
			time.Sleep(counterAsDuration + time.Second)
			continue
		}
		counter = 0

		payload, err := eventuallyDecodeNetinfoData(stdout)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		PostNestatInfoController(apiClient, payload, podname)

	}
}

func PostNestatInfoController(apiClient kubernetes.Interface, payload types.NestatInfoRequestBody, podname string) {
	filteredProbes := buildProbesFromMonitorContainer(apiClient, payload, podname)
	if len(filteredProbes) > 0 {
		eventstorage.AddProbes(filteredProbes)
		kubesondeDispatcher.SendToQueue(filteredProbes, kubesondeDispatcher.HIGH)
	}

}

func buildProbesFromMonitorContainer(apiClient kubernetes.Interface, payload types.NestatInfoRequestBody, podname string) []probe_command.KubesondeCommand {
	netInfoNotLoopback := findListeningPortsNonInLoopback(payload)
	// log.Info(fmt.Sprintf("Received monitor from %s \n%v", podname, payload))
	// Store only NON loopback listening ports
	state.AppendNetInfo(podname, &netInfoNotLoopback)

	// Should also execute new probes if the port is not already in the storage
	currPods := eventstorage.GetActivePods()
	if len(currPods) <= 1 {
		return []probe_command.KubesondeCommand{}
	}

	// Key by (port, protocol) rather than just port: a port used by both TCP
	// and UDP (e.g. statsd 9125, alertmanager 9094) must keep both entries -
	// a map[int32]string keyed on port alone would silently drop one.
	type portProtocol struct {
		port     int32
		protocol string
	}
	seenPortProtocols := map[portProtocol]bool{}
	var uniquePortProtocols []portProtocol
	for _, item := range netInfoNotLoopback {
		intport, err := strconv.ParseInt(item.Port, 10, 32)
		if err != nil {
			log.Error(err, "Invalid port number", "port", item.Port)
			continue
		}
		pp := portProtocol{port: int32(intport), protocol: strings.ToUpper(item.Protocol)}
		if !seenPortProtocols[pp] {
			seenPortProtocols[pp] = true
			uniquePortProtocols = append(uniquePortProtocols, pp)
		}
	}
	if len(uniquePortProtocols) == 0 {
		return []probe_command.KubesondeCommand{}
	}
	// FIXME: get namespace from declaration
	pod, err := apiClient.CoreV1().Pods(currPods[0].Namespace).Get(context.TODO(), podname, metav1.GetOptions{})
	if err != nil {
		return []probe_command.KubesondeCommand{}
	}

	intPorts := lo.Map(uniquePortProtocols, func(pp portProtocol, i int) int32 {
		return pp.port
	})
	protocols := lo.Map(uniquePortProtocols, func(pp portProtocol, i int) string {
		return pp.protocol
	})
	probes := probe_command.BuildTargetedCommandsToDestination(currPods, *pod, intPorts, protocols)
	if len(probes) == 0 {
		log.Info("No probes could be found")
		return []probe_command.KubesondeCommand{}
	}
	filteredProbes := lo.Filter(probes, func(cc probe_command.KubesondeCommand, i int) bool {
		return !eventstorage.ProbeAvailable(cc)
	})
	return filteredProbes

}

func findListeningPortsNonInLoopback(payload types.NestatInfoRequestBody) []v12.PodNetworkingItem {

	// 1 TCP
	// 2 UDP

	monitor := lo.Map(payload, func(entry types.NestatInfoRequestBodyItem, i int) v12.PodNetworkingItem {
		var protocol string
		if entry.Type == 1 {
			protocol = "TCP"
		} else {
			protocol = "UDP"
			// TODO: HOW ABOUT SCTP?
		}
		return v12.PodNetworkingItem{
			Port:     entry.Laddr[1],
			IP:       entry.Laddr[0],
			Protocol: protocol,
		}
	})
	netInfoNotLoopback := lo.Filter(monitor, func(item v12.PodNetworkingItem, i int) bool {
		return item.IP != "127.0.0.1"
	})
	return netInfoNotLoopback
}
