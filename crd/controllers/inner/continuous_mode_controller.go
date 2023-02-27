package inner

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Client    kubernetes.Interface
	Kubesonde v12.Kubesonde
}

func (state *KubesondeContinuousState) getClient() kubernetes.Interface {
	return state.Client
}

func (state *KubesondeContinuousState) runCommand(client kubernetes.Interface, namespace string, command probe_command.KubesondeCommand, checker func(string) bool) (bool, error) {
	return runRemoteCommandWithErrorHandler(client, namespace, command, checker)
}

func (state *KubesondeContinuousState) runGenericCommand(client kubernetes.Interface, namespace string, command probe_command.KubesondeCommand) (string, error) {
	return runGenericCommand(client, namespace, command)
}

func ephemeralContainerExists(pod v1.Pod) bool {
	// FIXME: check if container with `debugger` name exists
	ephCont := pod.Spec.EphemeralContainers
	var ephNames = lo.Map(ephCont, func(ec v1.EphemeralContainer, i int) string {
		return ec.Name
	})
	_, ok1 := lo.Find(ephNames, func(s string) bool {
		return s == "debugger"
	})
	_, ok2 := lo.Find(ephNames, func(s string) bool {
		return s == "monitor"
	})
	if len(ephNames) != 2 {
		log.Info(fmt.Sprintf("Pod %s has %v ephemeral containers", pod.Name, ephNames))
	}
	return ok1 && ok2
}

func toProbeError(kubesondeCommand probe_command.KubesondeCommand, err error) v12.ProbeOutputError {
	return v12.ProbeOutputError{
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
	}
}

func toProbeItem(kubesondeCommand probe_command.KubesondeCommand, result v12.ActionType) v12.ProbeOutputItem {
	return v12.ProbeOutputItem{
		Type:                 v12.PROBE,
		ExpectedAction:       kubesondeCommand.Action,
		DestinationHostnames: kubesondeCommand.DestinationHostnames,
		ResultingAction:      result,
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
	}
}

func withDeploymentInformation(client kubernetes.Interface, output v12.ProbeOutputItem) v12.ProbeOutputItem {
	// Source is always a pod
	source_pod, err_source := client.CoreV1().Pods(output.Source.Namespace).Get(context.TODO(), output.Source.Name, metav1.GetOptions{})
	dest_pod, err_dest := client.CoreV1().Pods(output.Destination.Namespace).Get(context.TODO(), output.Destination.Name, metav1.GetOptions{})
	if err_source == nil {
		s_replica, s_deployment := utils.GetReplicaAndDeployment(client, *source_pod)
		output.Source.ReplicaSetName = s_replica
		output.Source.DeploymentName = s_deployment
	}
	if err_dest == nil {
		d_replica, d_deployment := utils.GetReplicaAndDeployment(client, *dest_pod)
		output.Destination.ReplicaSetName = d_replica
		output.Destination.DeploymentName = d_deployment
	}

	return output
}
func InspectWithContinuousMode(mode KubesondeMode, commands []probe_command.KubesondeCommand) v12.ProbeOutput {
	// runCommand, client := state.runCommand, state.getClient()
	client := mode.getClient()

	for _, kubesondeCommand := range commands {
		pod, err := client.CoreV1().Pods(kubesondeCommand.Namespace).Get(context.TODO(), kubesondeCommand.SourcePodName, metav1.GetOptions{})
		if err != nil {
			continue
		}

		if !ephemeralContainerExists(*pod) {
			log.Info("Ephemeral containers are not ready")
			continue
		}

		result, err := mode.runCommand(client, kubesondeCommand.Namespace, kubesondeCommand, kubesondeCommand.ProbeChecker)

		if err != nil {
			errors := []v12.ProbeOutputError{toProbeError(kubesondeCommand, err)}
			state.AppendErrors(&errors)
			log.Info("Error when Probing...")
		} else if err == nil && result {
			probes := []v12.ProbeOutputItem{withDeploymentInformation(client, toProbeItem(kubesondeCommand, v12.ALLOW))}
			state.AppendProbes(&probes)
		} else if err == nil && !result {
			probes := []v12.ProbeOutputItem{withDeploymentInformation(client, toProbeItem(kubesondeCommand, v12.DENY))}
			state.AppendProbes(&probes)
		}
	}

	return state.GetProbeState()
}

func InspectAndStoreResult(client kubernetes.Interface, probes []probe_command.KubesondeCommand) {
	// log.Info("Probing...")
	probestate := new(KubesondeContinuousState)
	probestate.Client = client
	probeOutput := InspectWithContinuousMode(probestate, probes)

	state.AppendProbes(&probeOutput.Items)
	state.AppendErrors(&probeOutput.Errors)
}
