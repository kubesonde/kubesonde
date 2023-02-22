package recursiveprobing

import (
	"time"

	securityv1 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/dispatcher"
	eventstorage "kubesonde.io/controllers/event-storage"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("Recursive probing")

func RecursiveProbing(Kubesonde securityv1.Kubesonde, when time.Duration) {
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
	dispatcher.SendToQueue(probes, dispatcher.LOW)

}
