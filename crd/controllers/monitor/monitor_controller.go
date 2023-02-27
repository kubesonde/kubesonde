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
	"kubesonde.io/controllers/dispatcher"
	eventstorage "kubesonde.io/controllers/event-storage"
	"kubesonde.io/controllers/probe_command"
	"kubesonde.io/controllers/state"
	"kubesonde.io/rest_apis/types"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("Monitor controller")
var MAX_CONNECT_RETRIES = 6

// This function starts an infinite loop
func RunMonitorContainers(client kubernetes.Interface) {
	for {
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
			//	WaitEphemeralContainersToBeRunning(client, p)

			// log.Info(fmt.Sprintf("Running monitor on pod %s", p.Name))

			var stdout, stderr, err = debug_container.RunMonitorContainerProcess(client, p.Namespace, p.Name)
			if err != nil {
				log.Error(err, "Could not run monitor process")
				log.Info(err.Error())
				return
			}
			go ProcessNetInfo(client, stdout, stderr, p.Name)
			state.SetNestatPod(p.Name)
			// log.Info(fmt.Sprintf("Pod %s contains monitor ephemeral container", p.Name))

		})
		time.Sleep(10 * time.Second)
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
		dispatcher.SendToQueue(filteredProbes, dispatcher.HIGH)
	}

}

func buildProbesFromMonitorContainer(apiClient kubernetes.Interface, payload types.NestatInfoRequestBody, podname string) []probe_command.KubesondeCommand {
	netInfoNotLoopback := findListeningPortsNonInLoopback(payload)
	// log.Info(fmt.Sprintf("Received monitor from %s \n%v", podname, payload))
	// Store only NON loopback listening ports
	state.SetNetInfoV2(podname, &netInfoNotLoopback) // FIXME: we should create a set out of all the responses

	// Should also execute new probes if the port is not already in the storage
	currPods := eventstorage.GetActivePods()
	if len(currPods) <= 1 {
		return []probe_command.KubesondeCommand{}
	}

	var initVal = map[int32]string{}
	monitorMapping := lo.Reduce(netInfoNotLoopback, func(acc map[int32]string, item v12.PodNetworkingItem, i int) map[int32]string {
		intport := lo.Must1(strconv.ParseInt(item.Port, 10, 32))
		acc[int32(intport)] = item.Protocol
		return acc
	}, initVal)
	monitorPorts := lo.Map(netInfoNotLoopback, func(pni v12.PodNetworkingItem, i int) string {
		return pni.Port
	})
	if len(monitorPorts) == 0 {
		return []probe_command.KubesondeCommand{}
	}
	pods, err := apiClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	must(err)
	if err != nil {
		return []probe_command.KubesondeCommand{}
	}
	podsFound := lo.Filter(pods.Items, func(pod v1.Pod, i int) bool {
		return pod.Name == podname
	})
	if len(podsFound) == 0 {
		return []probe_command.KubesondeCommand{}
	}
	intPorts := lo.Map(monitorPorts, func(s string, i int) int32 {
		acc, err := strconv.ParseInt(s, 10, 32)
		must(err)
		return int32(acc)
	})
	var protocols = lo.Map(intPorts, func(value int32, _ int) string {
		return strings.ToUpper(monitorMapping[value])
	})
	probes := probe_command.BuildTargetedCommandsToDestination(currPods, podsFound[0], intPorts, protocols)
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

func must(err error) {
	if err != nil {
		log.Error(err, "Something went wrong")
	}
}
