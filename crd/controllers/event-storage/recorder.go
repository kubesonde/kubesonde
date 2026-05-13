package eventstorage

import (
	"sync"

	v1 "k8s.io/api/core/v1"
)

type CreatedPodRecord struct {
	Pod               v1.Pod
	DeploymentName    string
	ReplicaSetName    string
	CreationTimestamp int64
}

type DeletedPodRecord struct {
	Pod               v1.Pod
	DeploymentName    string
	CreationTimestamp int64
	DeletionTimestamp int64
}

var (
	_activePods  = make(map[string]CreatedPodRecord)
	_deletedPods = make(map[string]DeletedPodRecord)
	_services    []v1.Service
	storageMu    sync.RWMutex
)

func AddActivePod(key string, value CreatedPodRecord) {
	storageMu.Lock()
	_activePods[key] = value
	storageMu.Unlock()
}

func AddService(value v1.Service) {
	storageMu.Lock()
	_services = append(_services, value)
	storageMu.Unlock()
}

func GetServices() []v1.Service {
	storageMu.RLock()
	defer storageMu.RUnlock()
	return _services
}

func DeleteActivePod(key string) {
	storageMu.Lock()
	delete(_activePods, key)
	storageMu.Unlock()
}

func AddDeletedPod(key string, value DeletedPodRecord) {
	storageMu.Lock()
	_deletedPods[key] = value
	storageMu.Unlock()
}

func GetActivePodByName(key string) CreatedPodRecord {
	storageMu.RLock()
	defer storageMu.RUnlock()
	return _activePods[key]
}

func GetActivePods() []v1.Pod {
	storageMu.RLock()
	defer storageMu.RUnlock()
	v := make([]v1.Pod, 0, len(_activePods))
	for _, value := range _activePods {
		v = append(v, value.Pod)
	}
	return v
}

func GetActivePodNames() []string {
	storageMu.RLock()
	defer storageMu.RUnlock()
	keys := make([]string, len(_activePods))
	i := 0
	for k := range _activePods {
		keys[i] = k
		i++
	}
	return keys
}

func GetDeletedPodNames() []string {
	storageMu.RLock()
	defer storageMu.RUnlock()
	keys := make([]string, len(_deletedPods))
	i := 0
	for k := range _deletedPods {
		keys[i] = k
		i++
	}
	return keys
}

func ClearEventStorage() {
	storageMu.Lock()
	for k := range _activePods {
		delete(_activePods, k)
	}
	for k := range _deletedPods {
		delete(_deletedPods, k)
	}
	_services = nil
	storageMu.Unlock()
}
