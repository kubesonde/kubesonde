package inner

import (
	"strconv"

	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
	v12 "kubesonde.io/api/v1"
	kubesondemetrics "kubesonde.io/controllers/metrics"
	"kubesonde.io/controllers/probe_command"
)

type KubesondeTargetedState struct {
	Log       logr.Logger
	Client    *kubernetes.Clientset
	Kubesonde v12.Kubesonde
}

func (state *KubesondeTargetedState) getClient() *kubernetes.Clientset {
	return state.Client
}
func (state *KubesondeTargetedState) logError(err error, message string) {
	state.Log.Error(err, message)

}
func (state *KubesondeTargetedState) logInfo(value string) {
	state.Log.Info(value)

}

func (state *KubesondeTargetedState) getKubesonde() v12.Kubesonde {
	return state.Kubesonde

}
func (state *KubesondeTargetedState) runCommand(client *kubernetes.Clientset, namespace string, command probe_command.KubesondeCommand, checker func(string) bool) (bool, error) {
	return runRemoteCommand(client, namespace, command, checker), nil // TODO: Fixme

}

func (state *KubesondeTargetedState) runGenericCommand(client *kubernetes.Clientset, namespace string, command probe_command.KubesondeCommand) (string, error) {
	return runGenericCommand(client, namespace, command)

}

func (state *KubesondeTargetedState) AddEdge(command probe_command.KubesondeCommand, value bool) {
	// from string, to string, label string, value bool
	var expectedAction = false

	if command.Action == v12.ALLOW {
		expectedAction = true
	}

	kubesondemetrics.TargetedMetricsSummary.With(map[string]string{
		"from":        command.SourcePodName,
		"to":          command.Destination,
		"label":       command.DestinationPort,
		"exists":      strconv.FormatBool(value),
		"shouldExist": strconv.FormatBool(expectedAction),
	})

}
