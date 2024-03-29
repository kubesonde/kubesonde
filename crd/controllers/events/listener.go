package events

import (
	"fmt"
	"time"

	kubesondev1 "kubesonde.io/api/v1"
	eventstorage "kubesonde.io/controllers/event-storage"

	v1 "k8s.io/api/core/v1"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"kubesonde.io/controllers/state"
	"kubesonde.io/controllers/utils"
)

func podEventHandler(client kubernetes.Interface, Kubesonde kubesondev1.Kubesonde) cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)

			if utils.SourcePodMatchesKubesondeSpec(Kubesonde, *pod) {

				addPodEvent(client, *pod)
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			if utils.SourcePodMatchesKubesondeSpec(Kubesonde, *pod) {
				return
			}
			deletePodEvent(*pod)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newPod := newObj.(*v1.Pod)
			oldPod := oldObj.(*v1.Pod)
			if !utils.SourcePodMatchesKubesondeSpec(Kubesonde, *oldPod) {
				return
			}
			if oldPod.Status.Phase == v1.PodRunning && newPod.Status.Phase != v1.PodRunning {
				deletePodEvent(*oldPod)
			}
		},
	}
}

func svcEventHandler(Kubesonde kubesondev1.Kubesonde) cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			srv := obj.(*v1.Service)
			if utils.ServiceMatchesKubesondeSpec(Kubesonde, *srv) { // Services are namespaced
				currPods := eventstorage.GetActivePods()
				if len(currPods) == 0 {
					return
				}
				if srv.Name == "kubernetes" || srv.Name == "kube-dns" {
					return
				}
				if srv.Spec.ClusterIP == "None" || srv.Spec.ClusterIP == "" {
					log.V(1).Info(fmt.Sprintf("Service %s does not have a ClusterIP", srv.Name))
					return
				}
				log.V(1).Info(fmt.Sprintf("Service %s Probed", srv.Name))

				srvProbes := getServicesAsProbes(srv.Spec.Ports, currPods, srv.Spec.ClusterIP, srv.Namespace, srv.Name) // This should be an information event that tells that external connections can reach this service
				state.AppendProbes(&srvProbes)

			}

		},
	}
}

// Setup event listener for pods and services. When a new event is received, probes
// and ephemeral containers will be generated
func InitEventListener(client kubernetes.Interface, Kubesonde kubesondev1.Kubesonde) {
	fmt.Printf("Setting up the event listener...")
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(client, time.Second*5)
	podInformer := kubeInformerFactory.Core().V1().Pods().Informer()
	svcInformer := kubeInformerFactory.Core().V1().Services().Informer()

	podInformer.AddEventHandler(podEventHandler(client, Kubesonde))
	svcInformer.AddEventHandler(svcEventHandler(Kubesonde))

	stop := make(chan struct{})
	defer close(stop)
	kubeInformerFactory.Start(stop)
	for {
		activePods := getActivePodsEvent()
		time.Sleep(time.Second * 30)
		fmt.Println("Active pods: ", activePods)
	}
}
