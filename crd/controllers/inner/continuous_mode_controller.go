package inner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/samber/lo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v12 "kubesonde.io/api/v1"
	debug_container "kubesonde.io/controllers/debug-container"
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

func withDeploymentInformationSlow(client kubernetes.Interface, output v12.ProbeOutputItem) v12.ProbeOutputItem {
	//log.Info("Getting information about deployment, slowly...")
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
func withDeploymentInformation(client kubernetes.Interface, output v12.ProbeOutputItem) v12.ProbeOutputItem {
	// Source is always a pod
	curr_state := state.GetProbeState().Items
	src, source_in_state := lo.Find(curr_state, func(item v12.ProbeOutputItem) bool { return item.Source.Name == output.Source.Name })
	src_2, source_in_state_2 := lo.Find(curr_state, func(item v12.ProbeOutputItem) bool { return item.Destination.Name == output.Source.Name })
	if source_in_state {
		output.Source.ReplicaSetName = src.Source.ReplicaSetName
		output.Source.DeploymentName = src.Source.DeploymentName
	} else if source_in_state_2 {
		output.Source.ReplicaSetName = src_2.Destination.ReplicaSetName
		output.Source.DeploymentName = src_2.Destination.DeploymentName
	}

	dst, dst_in_state := lo.Find(curr_state, func(item v12.ProbeOutputItem) bool { return item.Source.Name == output.Destination.Name })
	dst_2, dst_in_state_2 := lo.Find(curr_state, func(item v12.ProbeOutputItem) bool { return item.Destination.Name == output.Destination.Name })
	if dst_in_state {
		output.Source.ReplicaSetName = dst.Source.ReplicaSetName
		output.Source.DeploymentName = dst.Source.DeploymentName
	} else if dst_in_state_2 {
		output.Destination.ReplicaSetName = dst_2.Destination.ReplicaSetName
		output.Destination.DeploymentName = dst_2.Destination.DeploymentName
	}

	if (source_in_state || source_in_state_2) && (dst_in_state || dst_in_state_2) {
		return output
	}

	return withDeploymentInformationSlow(client, output)
}

func fixOutput(s string) string {
	splits := strings.Split(s, "andstderr")
	if len(splits) == 1 {
		return limitStringLength(s)
	}
	return fmt.Sprintf("%s -- %s", limitStringLength(splits[0]), limitStringLength(splits[1]))
}

func limitStringLength(s string) string {
	if len(s) > 200 {
		return s[:200]
	}
	return s
}

func InspectWithContinuousMode(mode KubesondeMode, commands []probe_command.KubesondeCommand) v12.ProbeOutput {
	// runCommand, client := state.runCommand, state.getClient()
	client := mode.getClient()
	// FIXME: here I should return only the current probes.
	for _, kubesondeCommand := range commands {
		_, source_has_netinfo := lo.Find(state.GetNetstatPods(), func(item string) bool {
			return item == kubesondeCommand.SourcePodName
		})
		if kubesondeCommand.SourceType == v12.POD && !source_has_netinfo {
			pod, err := client.CoreV1().Pods(kubesondeCommand.Namespace).Get(context.TODO(), kubesondeCommand.SourcePodName, metav1.GetOptions{})
			if err != nil {
				continue
			}

			if !debug_container.EphemeralContainerExists(pod) || !debug_container.EphemeralContainersRunning(pod) {
				continue
			}
		}
		result, err := mode.runCommand(client, kubesondeCommand.Namespace, kubesondeCommand, kubesondeCommand.ProbeChecker)
		command := fmt.Sprintf("wget --server-response --timeout=3 -O- http://%s:%s", kubesondeCommand.DestinationIPAddress, kubesondeCommand.DestinationPort)
		debug_info := fmt.Sprintf("From: %s - Command: %s", kubesondeCommand.SourcePodName, command)

		var output string

		if kubesondeCommand.Protocol == "TCP" && kubesondeCommand.DestinationPort != "53" && kubesondeCommand.DestinationType != v12.INTERNET {
			genericCommand := kubesondeCommand
			genericCommand.Command = command
			output, _ = mode.runGenericCommand(client, kubesondeCommand.Namespace, genericCommand)
		} else {
			output = "SKIP"
		}
		if err != nil {
			errors := []v12.ProbeOutputError{toProbeError(kubesondeCommand, err)}
			state.AppendErrors(&errors)
			log.Info("Error when Probing...")
		} else if result {
			probe_output := withDeploymentInformation(client, toProbeItem(kubesondeCommand, v12.ALLOW))
			probe_output.DebugOutput = fixOutput(fmt.Sprintf("%s %s", debug_info, output))
			probes := []v12.ProbeOutputItem{probe_output}
			state.AppendProbes(&probes)
		} else if !result {
			probe_output := withDeploymentInformation(client, toProbeItem(kubesondeCommand, v12.DENY))
			probe_output.DebugOutput = fixOutput(fmt.Sprintf("%s %s", debug_info, output))
			probes := []v12.ProbeOutputItem{probe_output}
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
