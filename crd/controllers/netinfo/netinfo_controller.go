package netinfo

import (
	"bytes"
	"context"
	"encoding/json"
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
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("NetInfo controller")

func EventuallyRunNetinfo(client *kubernetes.Clientset) {
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

			log.Info(fmt.Sprintf("Running netinfo on pod %s", p.Name))

			var stdout, stderr, err = debug_container.RunNetinfoProcess(client, p.Namespace, p.Name)
			if err != nil {
				log.Error(err, "Could not run netinfo process")
				log.Info(err.Error())
				return
			}
			go ProcessNetInfo(stdout, stderr, p.Name)
			state.SetNestatPod(p.Name)
			//log.Info(fmt.Sprintf("Pod %s contains netinfo ephemeral container", p.Name))

		})
		time.Sleep(3 * time.Second)
	}
}

func ProcessNetInfo(stdout *bytes.Buffer, stderr *bytes.Buffer, podname string) {
	var counter = 0
	for {
		if stdout.Len() == 0 {
			counter += 1
			if counter >= 6 {
				log.Info(fmt.Sprintf("Restarting netinfo probes on %s", podname))
				log.Info(fmt.Sprintf("Stderr %s", stderr.String()))
				state.DeleteNetstatPod(podname)
				return
			}
			counterAsDuration := time.Duration(counter * 1000)
			time.Sleep(counterAsDuration + time.Second)
			continue
		}
		counter = 0
		var payload_raw = stdout.String()
		var index = strings.Index(payload_raw, "\n")
		if index < 0 { // Not found
			time.Sleep(3 * time.Second)
			continue
		}
		if index == 0 { // First char is \n
			stdout.Next(1)
			payload_raw = stdout.String()
			index = strings.Index(payload_raw, "\n")
			time.Sleep(3 * time.Second)
			continue
		}
		var potential_json = stdout.Next(index)
		var payload types.NestatInfoRequestBody
		var err = json.Unmarshal(potential_json, &payload)
		if err != nil {
			log.Error(err, "Could not decode netinfo")
			log.Info(string(payload_raw))
			time.Sleep(3 * time.Second)
			continue
		}
		PostNestatInfoController(payload, podname)
		time.Sleep(3 * time.Second)
	}
}

func PostNestatInfoController(payload types.NestatInfoRequestBody, podname string) {
	netInfoNotLoopback := findListeningPortsNonInLoopback(payload)
	log.Info(fmt.Sprintf("Received netinfo from %s \n%v", podname, payload))
	// Store only NON loopback listening ports
	state.SetNetInfoV2(podname, &netInfoNotLoopback) // FIXME: we should create a set out of all the responses

	// Should also execute new probes if the port is not already in the storage
	currPods := eventstorage.GetActivePods()
	if len(currPods) <= 1 {
		return
	}

	var initVal = map[int32]string{}
	netinfoMapping := lo.Reduce(netInfoNotLoopback, func(acc map[int32]string, item v12.PodNetworkingItem, i int) map[int32]string {
		intport, _ := strconv.Atoi(item.Port)
		acc[int32(intport)] = item.Protocol
		return acc
	}, initVal)
	netinfoPorts := lo.Map(netInfoNotLoopback, func(pni v12.PodNetworkingItem, i int) string {
		return pni.Port
	})
	if len(netinfoPorts) == 0 {
		return
	}
	clusterConfig := config.GetConfigOrDie()
	apiClient, err := kubernetes.NewForConfig(clusterConfig)
	must(err)
	if err != nil {
		return
	}
	pods, err := apiClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	must(err)
	if err != nil {
		return
	}
	podsFound := lo.Filter(pods.Items, func(pod v1.Pod, i int) bool {
		return pod.Name == podname
	})
	if len(podsFound) == 0 {
		return
	}
	intPorts := lo.Map(netinfoPorts, func(s string, i int) int32 {
		acc, err := strconv.Atoi(s)
		must(err)
		return int32(acc)
	})
	var protocols = lo.Map(intPorts, func(value int32, _ int) string {
		return strings.ToUpper(netinfoMapping[value])
	})
	probes := probe_command.BuildTargetedCommandsToDestination(currPods, podsFound[0], intPorts, protocols)
	if len(probes) == 0 {
		log.Info("No probes could be found")
		return
	}
	filteredProbes := lo.Filter(probes, func(cc probe_command.KubesondeCommand, i int) bool {
		return !eventstorage.ProbeAvailable(cc)
	})
	if len(filteredProbes) > 0 {
		eventstorage.AddProbes(filteredProbes)
		dispatcher.SendToQueue(filteredProbes, dispatcher.HIGH)
	}

}

func findListeningPortsNonInLoopback(payload types.NestatInfoRequestBody) []v12.PodNetworkingItem {

	// 1 TCP
	// 2 UDP
	// https://psutil.readthedocs.io/en/latest/

	netinfo := lo.Map(payload, func(entry types.NestatInfoRequestBodyItem, i int) v12.PodNetworkingItem {
		var protocol string
		if entry.Type == 1 {
			protocol = "TCP"
		} else {
			protocol = "UDP"
			// TODO: HOW ABOUT SCTP?
		}
		return v12.PodNetworkingItem{
			Port:     strconv.Itoa(int(entry.Laddr[1].(float64))),
			IP:       entry.Laddr[0].(string),
			Protocol: protocol,
		}
	})
	netInfoNotLoopback := lo.Filter(netinfo, func(item v12.PodNetworkingItem, i int) bool {
		return item.IP != "127.0.0.1"
	})
	return netInfoNotLoopback
}

func must(err error) {
	if err != nil {
		log.Error(err, "Something went wrong")
	}
}
