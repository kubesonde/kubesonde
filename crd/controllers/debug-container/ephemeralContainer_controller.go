package debug_container

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/samber/lo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	scheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("controllers.state")

func InstallEphameralContainers(client kubernetes.Interface, pods *v1.PodList) {

	for _, pod := range pods.Items {
		if !ephemeralContainerExists(pod) {
			installContainers(client, pod)
			log.V(1).Info(fmt.Sprintf("ephemeral container installed in %s pod", pod.Name))
		}

	}
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
	ok := ok1 && ok2
	if len(ephNames) != 2 {
		log.Info(fmt.Sprintf("Pod %s has %v ephemeral containers", pod.Name, ephNames))
	}
	return ok
}

func installContainers(client kubernetes.Interface, pod v1.Pod) {
	podJS, err := json.Marshal(pod)
	if err != nil {
		log.Error(err, "error creating JSON for pod: %s", pod.Name)
	}
	debugPod, err := generateDebugContainers(&pod)
	if err != nil {
		log.Error(err, "something went wrong")
	}
	debugJS, err := json.Marshal(debugPod)
	if err != nil {
		log.Error(err, "error creating JSON for debug container")
	}
	patch, err := strategicpatch.CreateTwoWayMergePatch(podJS, debugJS, pod)
	if err != nil {
		log.Error(err, "error creating patch to add debug container: %v")
	}

	pods := client.CoreV1().Pods(pod.Namespace)
	_, err = pods.Patch(context.TODO(), pod.Name, types.StrategicMergePatchType, patch, metav1.PatchOptions{}, "ephemeralcontainers")
	if err != nil {
		log.Error(err, fmt.Sprintf("Error while setting up the ephemeral container in pod %s", pod.Name))
		return
	}

}

func generateDebugContainers(pod *v1.Pod) (*v1.Pod, error) {
	privileged := true
	ec1 := &v1.EphemeralContainer{
		EphemeralContainerCommon: v1.EphemeralContainerCommon{
			Name:                     "debugger",
			Image:                    "instrumentisto/nmap:latest",
			ImagePullPolicy:          v1.PullIfNotPresent,
			Stdin:                    true,
			TerminationMessagePolicy: v1.TerminationMessageReadFile,
			TTY:                      true,
			Command:                  []string{"sh"},
			SecurityContext: &v1.SecurityContext{
				Privileged: &privileged,
			},
		},
	}
	ec2 := &v1.EphemeralContainer{
		EphemeralContainerCommon: v1.EphemeralContainerCommon{
			Name:                     "monitor",
			Image:                    "registry.cs.aalto.fi/kubesonde/monitor:latest",
			ImagePullPolicy:          v1.PullIfNotPresent,
			Stdin:                    true,
			TerminationMessagePolicy: v1.TerminationMessageReadFile,
			TTY:                      true,
			Command:                  []string{"sh"},
		}}

	copied := pod.DeepCopy()
	copied.Spec.EphemeralContainers = append(copied.Spec.EphemeralContainers, *ec1, *ec2)
	return copied, nil
}

func RunMonitorContainerProcess(client kubernetes.Interface, namespace string, sourcePodName string) (*bytes.Buffer, *bytes.Buffer, error) {

	pod, err := client.CoreV1().Pods(namespace).Get(context.TODO(), sourcePodName, metav1.GetOptions{})
	if err != nil {
		return bytes.NewBuffer(nil), bytes.NewBuffer(nil), err
	}

	if !ephemeralContainerExists(*pod) {
		log.Info("Ephemeral containers are not ready")
		return bytes.NewBuffer(nil), bytes.NewBuffer(nil), err
	}
	req := client.
		CoreV1().
		RESTClient().
		Post().
		Resource("pods").
		Namespace(namespace).
		Name(sourcePodName).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command:   []string{"/workspace/main"},
			Container: "monitor",
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(config.GetConfigOrDie(), "POST", req.URL())

	if err != nil {
		log.Error(err, "Remote Command failed")
		return bytes.NewBuffer(nil), bytes.NewBuffer(nil), err
	}
	var stdout, stderr bytes.Buffer
	go func() { // On the background try to enstablish again connection
		for {
			err = exec.Stream(remotecommand.StreamOptions{
				Stdin:  nil,
				Stdout: &stdout,
				Stderr: &stderr,
			})
			if err != nil {
				log.Info("Monitor container not found in Pod %s", sourcePodName)
				time.Sleep(3 * time.Second)
			}
		}
	}()
	// log.Info(fmt.Sprintf("Output for command: %s\nSource Pod: %s\nDestination : %s:%s\nStdout:\n%s\n---------\nStderr:\n%s\n",
	// command.Command, command.SourcePodName, command.Destination, command.DestinationPort, stdout.String(), stderr.String()))
	return &stdout, &stderr, nil
}
