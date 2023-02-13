package utils

import (
	"context"

	lo "github.com/samber/lo"
	v1 "k8s.io/api/apps/v1"
	k8sAPI "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

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
	var acceptedGlobNamespace = []string{"", "all"}
	if contains(acceptedGlobNamespace, configurationNamespace) {
		return true
	} else {
		return configurationNamespace == podNamespace
	}

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
