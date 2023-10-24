package recursiveprobing

import (
	"time"

	kubesondev1 "kubesonde.io/api/v1"
	kubesondeDispatcher "kubesonde.io/controllers/dispatcher"
	eventstorage "kubesonde.io/controllers/event-storage"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("Recursive probing")

// This function starts an infinite loop that runs all the probes at regular
// intervals
func RecursiveProbing(Kubesonde kubesondev1.Kubesonde, when time.Duration) {
	var task = func() {
		go RunProbing()
		RecursiveProbing(Kubesonde, when)
	}
	time.AfterFunc(when, task)
}

func RunProbing() {
	var probes = eventstorage.GetProbes()

	if len(probes) <= 1 {
		log.Info("Not enough probes")
		return
	}

	log.Info("Running all probes again")
	kubesondeDispatcher.SendToQueue(probes, kubesondeDispatcher.LOW)

}
