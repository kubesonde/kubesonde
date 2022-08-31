package state

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"golang.org/x/sync/semaphore"
	v1 "kubesonde.io/api/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	innerState v1.ProbeOutput = v1.ProbeOutput{
		Items:                      []v1.ProbeOutputItem{},
		Errors:                     []v1.ProbeOutputError{},
		PodNetworking:              []v1.PodNetworkingInfo{},
		PodNetworkingV2:            make(v1.PodNetworkingInfoV2),
		PodConfigurationNetworking: make(v1.PodNetworkingInfoV2),
	}
	sem             = semaphore.NewWeighted(1)
	podsWithNetstat = []string{}
	log             = logf.Log.WithName("controllers.state")
)

func SetNestatPod(pod string) {
	podsWithNetstat = append(podsWithNetstat, pod)
}
func DeleteNetstatPod(pod string) {
	podsWithNetstat = lo.Filter(podsWithNetstat, func(s string, i int) bool {
		return s != pod
	})
}

func GetNetstatPods() []string {
	return podsWithNetstat
}

func SetProbeState(probes *v1.ProbeOutput) {
	log.Info(fmt.Sprintf("Setting probe state. Items: %d  - Errors %d", len(probes.Items), len(probes.Errors)))
	must(sem.Acquire(context.Background(), 1))
	innerState = *probes
	sem.Release(1)
	return
}

func GetProbeState() v1.ProbeOutput {

	return innerState
}

func AppendProbes(items *[]v1.ProbeOutputItem) {
	must(sem.Acquire(context.Background(), 1))
	var newItems = append(innerState.Items, *items...)
	innerState.Items = lo.UniqBy(newItems, func(poi v1.ProbeOutputItem) v1.ComparableProbeOutputItem {
		return poi.ToComparableProbe()
	})
	sem.Release(1)
}

func AppendErrors(items *[]v1.ProbeOutputError) {
	must(sem.Acquire(context.Background(), 1))
	var newItems = append(innerState.Errors, *items...)
	innerState.Errors = lo.UniqBy(newItems, func(poe v1.ProbeOutputError) v1.ComparableProbeOutputItem {
		return poe.Value.ToComparableProbe()
	})
	sem.Release(1)
}

func AppendNetInfo(items *[]v1.PodNetworkingInfo) {
	must(sem.Acquire(context.Background(), 1))
	innerState.PodNetworking = append(innerState.PodNetworking, *items...)
	sem.Release(1)
}

func AppendNetInfoV2(key string, items *[]v1.PodNetworkingItem) {
	must(sem.Acquire(context.Background(), 1))
	innerState.PodNetworkingV2[key] = append(innerState.PodNetworkingV2[key], *items...)
	sem.Release(1)
}
func AppendConfig(key string, items *[]v1.PodNetworkingItem) {
	must(sem.Acquire(context.Background(), 1))
	innerState.PodConfigurationNetworking[key] = append(innerState.PodNetworkingV2[key], *items...)
	sem.Release(1)
}
func SetConfig(key string, items *[]v1.PodNetworkingItem) {
	must(sem.Acquire(context.Background(), 1))
	innerState.PodConfigurationNetworking[key] = *items
	sem.Release(1)
}

func SetNetInfoV2(key string, items *[]v1.PodNetworkingItem) {
	must(sem.Acquire(context.Background(), 1))
	innerState.PodNetworkingV2[key] = *items
	sem.Release(1)
}

func must(err error) {
	if err != nil {
		log.Error(err, "Could not acquire lock")
		panic(err)
	}
}
