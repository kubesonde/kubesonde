package inner

import (
	"time"

	"github.com/go-logr/logr"
	"k8s.io/client-go/kubernetes"
	v12 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/probe_command"
	"kubesonde.io/controllers/state"
	"kubesonde.io/controllers/utils"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("Kubesonde Runner")

type KubesondeContinuousState struct {
	Log       logr.Logger
	Client    *kubernetes.Clientset
	Kubesonde v12.Kubesonde
}

func (state *KubesondeContinuousState) AddEdge(command probe_command.KubesondeCommand, value bool) {
	panic("implement me")
}

func (state *KubesondeContinuousState) logInfo(s string) {
	panic("implement me")
}

func (state *KubesondeContinuousState) getKubesonde() v12.Kubesonde {
	panic("implement me")
}

func (state *KubesondeContinuousState) logError(err error, message string) {
	panic("implement me")
}

func (state *KubesondeContinuousState) getClient() *kubernetes.Clientset {
	return state.Client
}

func (state *KubesondeContinuousState) runCommand(client *kubernetes.Clientset, namespace string, command probe_command.KubesondeCommand, checker func(string) bool) (bool, error) {
	return runRemoteCommandWithErrorHandler(client, namespace, command, checker)
}

func (state *KubesondeContinuousState) runGenericCommand(client *kubernetes.Clientset, namespace string, command probe_command.KubesondeCommand) (string, error) {
	return runGenericCommand(client, namespace, command)
}

func InspectWithContinuousMode(mode KubesondeMode, commands []probe_command.KubesondeCommand) v12.ProbeOutput {
	//runCommand, client := state.runCommand, state.getClient()
	client := mode.getClient()

	for _, kubesondeCommand := range commands {

		result, err := mode.runCommand(client, kubesondeCommand.Namespace, kubesondeCommand, kubesondeCommand.ProbeChecker)

		if err != nil {
			errors := []v12.ProbeOutputError{
				{
					Value: v12.ProbeOutputItem{
						Type:            v12.PROBE,
						ExpectedAction:  kubesondeCommand.Action,
						ResultingAction: v12.DENY,
						Source: v12.ProbeEndpointInfo{
							Type:      kubesondeCommand.SourceType,
							Name:      kubesondeCommand.SourcePodName,
							Namespace: kubesondeCommand.Namespace,
							IPAddress: kubesondeCommand.SourceIPAddress,
						},
						Destination: v12.ProbeEndpointInfo{
							Type:      kubesondeCommand.DestinationType,
							Name:      kubesondeCommand.Destination,
							Namespace: kubesondeCommand.Namespace,
							IPAddress: kubesondeCommand.DestinationIPAddress,
						},
						DestinationHostnames: kubesondeCommand.DestinationHostnames,
						Protocol:             kubesondeCommand.Protocol,
						Port:                 kubesondeCommand.DestinationPort,
						Timestamp:            time.Now().Unix(),
					},
					Reason: err.Error(),
				}}
			state.AppendErrors(&errors)
			log.Info("Error when Probing...")
		} else if err == nil && result == true {
			probes := []v12.ProbeOutputItem{{
				Type:                 v12.PROBE,
				ExpectedAction:       kubesondeCommand.Action,
				DestinationHostnames: kubesondeCommand.DestinationHostnames,
				ResultingAction:      v12.ALLOW,
				Source: v12.ProbeEndpointInfo{
					Type:      kubesondeCommand.SourceType,
					Name:      kubesondeCommand.SourcePodName,
					Namespace: kubesondeCommand.Namespace,
					IPAddress: kubesondeCommand.SourceIPAddress,
				},
				Destination: v12.ProbeEndpointInfo{
					Type:      kubesondeCommand.DestinationType,
					Name:      kubesondeCommand.Destination,
					Namespace: kubesondeCommand.Namespace,
					IPAddress: kubesondeCommand.DestinationIPAddress,
				},
				Port:      kubesondeCommand.DestinationPort,
				Protocol:  kubesondeCommand.Protocol,
				Timestamp: time.Now().Unix(),
			}}
			state.AppendProbes(&probes)
		} else if err == nil && result == false {
			probes := []v12.ProbeOutputItem{{
				Type:                 v12.PROBE,
				ExpectedAction:       kubesondeCommand.Action,
				DestinationHostnames: kubesondeCommand.DestinationHostnames,
				ResultingAction:      v12.DENY,
				Source: v12.ProbeEndpointInfo{
					Type:      kubesondeCommand.SourceType,
					Name:      kubesondeCommand.SourcePodName,
					Namespace: kubesondeCommand.Namespace,
					IPAddress: kubesondeCommand.SourceIPAddress,
				},
				Destination: v12.ProbeEndpointInfo{
					Type:      kubesondeCommand.DestinationType,
					Name:      kubesondeCommand.Destination,
					Namespace: kubesondeCommand.Namespace,
					IPAddress: kubesondeCommand.DestinationIPAddress,
				},
				Port:      kubesondeCommand.DestinationPort,
				Protocol:  kubesondeCommand.Protocol,
				Timestamp: time.Now().Unix(),
			}}
			state.AppendProbes(&probes)
		}
	}

	return state.GetProbeState()
}

func InspectAndStoreResult(client *kubernetes.Clientset, probes []probe_command.KubesondeCommand) {
	log.Info("Probing...")
	probestate := new(KubesondeContinuousState)
	probestate.Client = client
	probeOutput := InspectWithContinuousMode(probestate, probes)
	//state.AppendNetInfo(&probeOutput.PodNetworking)
	deployments := utils.GetDeploymentNamesInNamespace(client, probes[0].Namespace)
	replicas := utils.GetReplicaSetsNamesInNamespace(client, probes[0].Namespace)
	enriched_state := state.EnrichState(&probeOutput, replicas, deployments)
	state.AppendProbes(&enriched_state.Items)
	state.AppendErrors(&enriched_state.Errors)
}
