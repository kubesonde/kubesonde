package inner

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
	"kubesonde.io/controllers/probe_command"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type KubesondeMode interface {
	getClient() kubernetes.Interface
	runCommand(client kubernetes.Interface, namespace string, command probe_command.KubesondeCommand, checker func(string) bool) (bool, error)
	runGenericCommand(client kubernetes.Interface, namespace string, command probe_command.KubesondeCommand) (string, error)
}

func runGenericCommand(client kubernetes.Interface, namespace string, command probe_command.KubesondeCommand) (string, error) {
	req := client.
		CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Namespace(namespace).
		Name(command.SourcePodName).
		SubResource("exec").Timeout(time.Second * 5)

	scheme := runtime.NewScheme()
	if err := v1.AddToScheme(scheme); err != nil {
		log.Error(err, "error adding to scheme")
		return err.Error(), err
	}
	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&v1.PodExecOptions{
		Command:   strings.Fields(command.Command),
		Container: command.ContainerName,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config.GetConfigOrDie(), "POST", req.URL())
	if err != nil {
		log.Error(err, "Remote Command failed")
		return err.Error(), err
	}
	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		/*log.Info(fmt.Sprintf(`
		Namespace: %s,
		Endpoint: %s
		Error: %v
		When running command: %s
		Source Pod: %s
		Destination : %s:%s,
		Stdout: %s
		Stderr : %s
		`, namespace, req.URL().String(), err, command.Command, command.SourcePodName, command.Destination, command.DestinationPort, &stdout, &stderr))*/
		return err.Error(), err
	}
	// log.Info(fmt.Sprintf("Output for command: %s\nSource Pod: %s\nDestination : %s:%s\nStdout:\n%s\n---------\nStderr:\n%s\n",
	//	command.Command, command.SourcePodName, command.Destination, command.DestinationPort, stdout.String(), stderr.String()))
	return fmt.Sprintf("%s andstderr %s", stdout.String(), stderr.String()), nil
	//return stdout.String(), nil
}

func runRemoteCommandWithErrorHandler(client kubernetes.Interface, namespace string, command probe_command.KubesondeCommand, checker func(string) bool) (bool, error) {
	req := client.
		CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Namespace(namespace).
		Name(command.SourcePodName).
		SubResource("exec").Timeout(time.Second * 5)

	scheme := runtime.NewScheme()
	if err := v1.AddToScheme(scheme); err != nil {
		log.Error(err, "error adding to scheme")
		return false, err
	}
	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&v1.PodExecOptions{
		Command:   strings.Fields(command.Command),
		Container: command.ContainerName,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config.GetConfigOrDie(), "POST", req.URL())
	if err != nil {
		log.Error(err, "Remote Command failed")
		return false, err
	}
	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		/*log.Info(fmt.Sprintf(`
		Namespace: %s,
		Endpoint: %s
		Error: %v
		When running command: %s
		Source Pod: %s
		Destination : %s:%s,
		Stdout: %s
		Stderr : %s
		`, namespace, req.URL().String(), err, command.Command, command.SourcePodName, command.Destination, command.DestinationPort, &stdout, &stderr))*/
		return false, err
	}
	// log.Info(fmt.Sprintf("Output for command: %s\nSource Pod: %s\nDestination : %s:%s\nStdout:\n%s\n---------\nStderr:\n%s\n",
	//	command.Command, command.SourcePodName, command.Destination, command.DestinationPort, stdout.String(), stderr.String()))

	return checker(stdout.String()), nil
}
