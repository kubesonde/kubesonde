package recursiveprobing

import (
	"context"
	"fmt"
	"time"

	kubesondev1 "kubesonde.io/api/v1"
	kubesondeDispatcher "kubesonde.io/controllers/dispatcher"
	eventstorage "kubesonde.io/controllers/event-storage"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("Recursive probing")

// This function starts an infinite loop that runs all the probes at regular
// intervals until ctx is cancelled (e.g. when the Kubesonde resource is deleted).
func RecursiveProbing(ctx context.Context, Kubesonde kubesondev1.Kubesonde, when time.Duration) {
	if ctx.Err() != nil {
		return
	}
	var task = func() {
		if ctx.Err() != nil {
			log.Info("Recursive probing stopped")
			return
		}
		size := kubesondeDispatcher.QueueSize()
		log.Info(fmt.Sprintf("Probe queue size: %d", size))
		if size == 0 {
			go RunProbing()
		}

		RecursiveProbing(ctx, Kubesonde, when)
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
