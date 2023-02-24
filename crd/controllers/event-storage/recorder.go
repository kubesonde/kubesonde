package eventstorage

import v1 "k8s.io/api/core/v1"

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

var _activePods = make(map[string]CreatedPodRecord)
var _deletedPods = make(map[string]DeletedPodRecord)
var _services []v1.Service

func AddActivePod(key string, value CreatedPodRecord) {
	_activePods[key] = value
}

func AddService(value v1.Service) {
	_services = append(_services, value)
}

func GetServices() []v1.Service {
	return _services
}
func DeleteActivePod(key string) {
	delete(_activePods, key)
}

func AddDeletedPod(key string, value DeletedPodRecord) {
	_deletedPods[key] = value
}

func GetActivePodByName(key string) CreatedPodRecord {
	return _activePods[key]
}
func GetActivePods() []v1.Pod {
	v := make([]v1.Pod, 0, len(_activePods))

	for _, value := range _activePods {
		v = append(v, value.Pod)
	}
	return v
}

func GetActivePodNames() []string {
	keys := make([]string, len(_activePods))

	i := 0
	for k := range _activePods {
		keys[i] = k
		i++
	}
	return keys
}
func GetDeletedPodNames() []string {
	keys := make([]string, len(_deletedPods))

	i := 0
	for k := range _deletedPods {
		keys[i] = k
		i++
	}
	return keys
}
