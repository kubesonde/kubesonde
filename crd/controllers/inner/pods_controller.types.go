package inner

import (
	"strconv"

	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
	v12 "kubesonde.io/api/v1"
	kubesondemetrics "kubesonde.io/controllers/metrics"
	"kubesonde.io/controllers/probe_command"
)

type KubesondeState struct {
	Log       logr.Logger
	Client    *kubernetes.Clientset
	Kubesonde v12.Kubesonde
}

func (state *KubesondeState) getClient() *kubernetes.Clientset {
	return state.Client
}
func (state *KubesondeState) logError(err error, message string) {
	state.Log.Error(err, message)

}
func (state *KubesondeState) logInfo(value string) {
	state.Log.Info(value)

}

func (state *KubesondeState) getKubesonde() v12.Kubesonde {
	return state.Kubesonde

}

func (state *KubesondeState) runGenericCommand(client *kubernetes.Clientset, namespace string, command probe_command.KubesondeCommand) (string, error) {
	return runGenericCommand(client, namespace, command)

}
func (state *KubesondeState) runCommand(client *kubernetes.Clientset, namespace string, command probe_command.KubesondeCommand, checker func(string) bool) (bool, error) {
	return runRemoteCommand(client, namespace, command, checker), nil // TODO: FIX this

}

func (state *KubesondeState) AddEdge(command probe_command.KubesondeCommand, value bool) {
	// from string, to string, label string, value bool

	kubesondemetrics.MetricsSummary.With(map[string]string{
		"from":   command.SourcePodName,
		"to":     command.Destination,
		"label":  command.DestinationPort,
		"exists": strconv.FormatBool(value),
	})

}
