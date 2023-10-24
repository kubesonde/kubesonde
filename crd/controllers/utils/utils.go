package utils

import (
	"context"
	"errors"
	"strings"

	lo "github.com/samber/lo"
	v1 "k8s.io/api/apps/v1"
	k8sAPI "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kubesondev1 "kubesonde.io/api/v1"
)

var ALLOW_ALL = "all"

func FilterPodsByStatus(pods *k8sAPI.PodList, status k8sAPI.PodPhase) k8sAPI.PodList {
	filteredPods := lo.Filter(pods.Items, func(p k8sAPI.Pod, i int) bool {
		return p.Status.Phase == status
	})

	return k8sAPI.PodList{
		TypeMeta: pods.TypeMeta,
		ListMeta: pods.ListMeta,
		Items:    filteredPods,
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func InNamespace(configurationNamespace string, podNamespace string) bool {
	var acceptedGlobNamespace = []string{"", ALLOW_ALL}
	if contains(acceptedGlobNamespace, configurationNamespace) {
		return true
	} else {
		return configurationNamespace == podNamespace
	}

}

func SourcePodMatchesKubesondeSpec(Kubesonde kubesondev1.Kubesonde, pod k8sAPI.Pod) bool {

	// If within include, allow
	for _, item := range Kubesonde.Spec.Include {
		if item.FromPodSelector == ALLOW_ALL { // Keyword
			return true
		}
		// Split key-value
		kv := strings.Split(item.FromPodSelector, "=") // TODO: Error handling
		key := kv[0]
		value := kv[1]
		podValue, ok := pod.ObjectMeta.Labels[key]
		if ok && podValue == value { // If key and value match
			return true
		}
	}
	// Exclude does not give information about source Pod.

	// If allow by default check in namespace
	if Kubesonde.Spec.Probe == ALLOW_ALL {
		return InNamespace(Kubesonde.Spec.Namespace, pod.Namespace)
	}
	// Otherwise we only accepts the included items.
	return false
}

func GetDeployment(replica v1.ReplicaSet) (string, error) {

	refs := replica.OwnerReferences
	depRefs := lo.Filter(refs, func(ref metav1.OwnerReference, idx int) bool {
		return ref.Kind == "Deployment"
	})
	depName := lo.Map(depRefs, func(ref metav1.OwnerReference, idx int) string {
		return ref.Name
	})
	if len(depName) == 1 {
		return depName[0], nil
	}
	return "", errors.New("no deployment")
}

func GetReplicaAndDeployment(client kubernetes.Interface, pod k8sAPI.Pod) (string, string) {
	if replicaSetName, err := GetReplicaSet(pod); err != nil {
		return "", ""
	} else {
		replicaSet, err := client.AppsV1().ReplicaSets(pod.Namespace).Get(context.TODO(), replicaSetName, metav1.GetOptions{})
		if err != nil {
			return replicaSetName, ""
		}
		if deployment, err2 := GetDeployment(*replicaSet); err2 != nil {
			return replicaSetName, ""
		} else {
			return replicaSetName, deployment
		}
	}
}
func GetReplicaSet(pod k8sAPI.Pod) (string, error) {

	refs := pod.OwnerReferences
	depRefs := lo.Filter(refs, func(ref metav1.OwnerReference, idx int) bool {
		return ref.Kind == "ReplicaSet"
	})
	depName := lo.Map(depRefs, func(ref metav1.OwnerReference, idx int) string {
		return ref.Name
	})
	if len(depName) == 1 {
		return depName[0], nil
	}
	return "", errors.New("no replicas")
}

func GetDeploymentNamesInNamespace(client kubernetes.Interface, namespace string) []string {
	deployments, _ := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	names := lo.Map(deployments.Items, func(d v1.Deployment, i int) string {
		return d.Name
	})
	return names
}
func GetReplicaSetsNamesInNamespace(client kubernetes.Interface, namespace string) []string {
	replicas, _ := client.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
	names := lo.Map(replicas.Items, func(d v1.ReplicaSet, i int) string {
		return d.Name
	})
	return names
}
