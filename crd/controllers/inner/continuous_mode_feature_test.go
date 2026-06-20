package inner

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	v1 "kubesonde.io/api/v1"
	"kubesonde.io/controllers/probe_command"
	"kubesonde.io/controllers/state"
)

func TestContinuousMode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ContinuousMode")
}

var _ = Describe("ContinuousMode", func() {
	BeforeEach(func() {
		innerState := v1.ProbeOutput{
			Items:           []v1.ProbeOutputItem{},
			Errors:          []v1.ProbeOutputError{},
			PodNetworking:   []v1.PodNetworkingInfo{},
			PodNetworkingV2: make(v1.PodNetworkingInfoV2),
		}
		state.SetProbeState(&innerState)
	})

	It("Records errors", func() {
		wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)

		// Use gomonkey instead of bou.ke/monkey
		patches := gomonkey.ApplyFunc(time.Now, func() time.Time { return wayback })
		defer patches.Reset()

		state := new(MockedCNIState)
		command := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationHostnames: nil,
			DestinationNamespace: "default",
			DestinationLabels:    "app=dest-pod;type=test",
			SourcePodName:        "test-pod",
			SourceLabels:         "app=source-pod;type=test",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}

		client := fake.NewSimpleClientset()
		p := &v12.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
			Status: v12.PodStatus{
				EphemeralContainerStatuses: []v12.ContainerStatus{
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
				},
			},
			Spec: v12.PodSpec{
				EphemeralContainers: []v12.EphemeralContainer{
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "debugger"}},
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "monitor"}},
				},
			},
		}
		_, err := client.CoreV1().Pods("default").Create(context.TODO(), p, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}

		state.Mock.On("getClient").Return(client)
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("this is an error"))

		output := InspectWithContinuousMode(state, []probe_command.KubesondeCommand{command})

		Expect(output.Errors).To(BeEquivalentTo(
			[]v1.ProbeOutputError{
				{
					Value: v1.ProbeOutputItem{
						Type:            v1.PROBE,
						ExpectedAction:  v1.ALLOW,
						ResultingAction: v1.DENY,
						Source: v1.ProbeEndpointInfo{
							Name:      "test-pod",
							Namespace: "default",
							Labels:    "app=source-pod;type=test",
						},
						Destination: v1.ProbeEndpointInfo{
							Name:      "test-destination",
							Namespace: "default",
							Labels:    "app=dest-pod;type=test",
						},
						Port:      "80",
						Timestamp: wayback.Unix(),
					},
					Reason: "this is an error",
				},
			},
		))
	})

	It("Records error and success", func() {
		state := new(MockedCNIState)
		wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)

		patches := gomonkey.ApplyFunc(time.Now, func() time.Time { return wayback })
		defer patches.Reset()

		errorCommand := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}
		successCommand := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "8080",
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}

		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("this is an error")).Once()
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).Once()

		client := fake.NewSimpleClientset()
		p := &v12.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
			Status: v12.PodStatus{
				EphemeralContainerStatuses: []v12.ContainerStatus{
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
				},
			},
			Spec: v12.PodSpec{
				EphemeralContainers: []v12.EphemeralContainer{
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "debugger"}},
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "monitor"}},
				},
			},
		}
		_, err := client.CoreV1().Pods("default").Create(context.TODO(), p, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}

		state.Mock.On("getClient").Return(client)

		output := InspectWithContinuousMode(state, []probe_command.KubesondeCommand{errorCommand, successCommand})

		Expect(output.Errors).To(BeEquivalentTo(
			[]v1.ProbeOutputError{
				{
					Value: v1.ProbeOutputItem{
						Type:            v1.PROBE,
						ExpectedAction:  v1.ALLOW,
						ResultingAction: v1.DENY,
						Source:          v1.ProbeEndpointInfo{Name: "test-pod", Namespace: "default"},
						Destination:     v1.ProbeEndpointInfo{Name: "test-destination", Namespace: "default"},
						Port:            "80",
						Timestamp:       wayback.Unix(),
					},
					Reason: "this is an error",
				},
			},
		))
		Expect(len(output.Items)).To(BeIdenticalTo(2))
	})

	It("Records a DENY probe in Items when probe fails with error", func() {
		wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)

		patches := gomonkey.ApplyFunc(time.Now, func() time.Time { return wayback })
		defer patches.Reset()

		state := new(MockedCNIState)
		command := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationNamespace: "default",
			DestinationLabels:    "app=dest-pod;type=test",
			SourcePodName:        "test-pod",
			SourceLabels:         "app=source-pod;type=test",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}

		client := fake.NewSimpleClientset()
		p := &v12.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
			Status: v12.PodStatus{
				EphemeralContainerStatuses: []v12.ContainerStatus{
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
				},
			},
			Spec: v12.PodSpec{
				EphemeralContainers: []v12.EphemeralContainer{
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "debugger"}},
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "monitor"}},
				},
			},
		}
		_, err := client.CoreV1().Pods("default").Create(context.TODO(), p, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}

		state.Mock.On("getClient").Return(client)
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("connection refused"))

		output := InspectWithContinuousMode(state, []probe_command.KubesondeCommand{command})

		// Verify error is recorded
		Expect(len(output.Errors)).To(BeIdenticalTo(1))
		Expect(output.Errors[0].Reason).To(ContainSubstring("connection refused"))

		// Verify a DENY probe is also recorded in Items
		Expect(len(output.Items)).To(BeIdenticalTo(1))
		Expect(output.Items[0].Source.Name).To(Equal("test-pod"))
		Expect(output.Items[0].Destination.Name).To(Equal("test-destination"))
		Expect(output.Items[0].Port).To(Equal("80"))
		Expect(output.Items[0].ResultingAction).To(Equal(v1.DENY))
	})

	It("Records both DENY probe and error when probe fails, plus ALLOW when next probe succeeds", func() {
		wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)

		patches := gomonkey.ApplyFunc(time.Now, func() time.Time { return wayback })
		defer patches.Reset()

		errorCommand := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}
		successCommand := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "8080",
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}
		denyCommand := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "9090",
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}

		state := new(MockedCNIState)
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("connection refused")).Once()
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).Once()
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, nil).Once()

		client := fake.NewSimpleClientset()
		p := &v12.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
			Status: v12.PodStatus{
				EphemeralContainerStatuses: []v12.ContainerStatus{
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
				},
			},
			Spec: v12.PodSpec{
				EphemeralContainers: []v12.EphemeralContainer{
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "debugger"}},
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "monitor"}},
				},
			},
		}
		_, err := client.CoreV1().Pods("default").Create(context.TODO(), p, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}

		state.Mock.On("getClient").Return(client)

		output := InspectWithContinuousMode(state, []probe_command.KubesondeCommand{
			errorCommand,
			successCommand,
			denyCommand,
		})

		// Error command: 1 error + 1 DENY probe in Items
		// Success command: 1 ALLOW probe in Items
		// Deny command (result=false, no error): 1 DENY probe in Items
		// Total Items: 3
		Expect(len(output.Errors)).To(BeIdenticalTo(1))

		Expect(len(output.Items)).To(BeIdenticalTo(3))

		// Verify the DENY probe from the error path
		var foundErrorProbe, foundAllow, foundDeny int
		for _, item := range output.Items {
			if item.Port == "80" {
				foundErrorProbe++
				Expect(item.ResultingAction).To(Equal(v1.DENY))
			}
			if item.Port == "8080" {
				foundAllow++
				Expect(item.ResultingAction).To(Equal(v1.ALLOW))
			}
			if item.Port == "9090" {
				foundDeny++
				Expect(item.ResultingAction).To(Equal(v1.DENY))
			}
		}
		Expect(foundErrorProbe).To(BeEquivalentTo(1))
		Expect(foundAllow).To(BeEquivalentTo(1))
		Expect(foundDeny).To(BeEquivalentTo(1))
	})

	It("Deduplicates errors and probes when two identical probes both fail", func() {
		wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)

		patches := gomonkey.ApplyFunc(time.Now, func() time.Time { return wayback })
		defer patches.Reset()

		command1 := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}
		command2 := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}

		state := new(MockedCNIState)
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("connection refused")).Times(2)

		client := fake.NewSimpleClientset()
		p := &v12.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
			Status: v12.PodStatus{
				EphemeralContainerStatuses: []v12.ContainerStatus{
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
				},
			},
			Spec: v12.PodSpec{
				EphemeralContainers: []v12.EphemeralContainer{
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "debugger"}},
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "monitor"}},
				},
			},
		}
		_, err := client.CoreV1().Pods("default").Create(context.TODO(), p, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}

		state.Mock.On("getClient").Return(client)

		output := InspectWithContinuousMode(state, []probe_command.KubesondeCommand{command1, command2})

		// Both probes fail → only ONE error (deduplicated by ComparableProbeOutputItem)
		Expect(len(output.Errors)).To(BeIdenticalTo(1))
		Expect(output.Errors[0].Reason).To(ContainSubstring("connection refused"))

		// Both probes fail → only ONE probe in Items (deduplicated)
		Expect(len(output.Items)).To(BeIdenticalTo(1))
		Expect(output.Items[0].ResultingAction).To(Equal(v1.DENY))
		Expect(output.Items[0].Port).To(Equal("80"))
	})

	It("State is correctly updated when error probe is followed by successful probe on same port", func() {
		wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)

		patches := gomonkey.ApplyFunc(time.Now, func() time.Time { return wayback })
		defer patches.Reset()

		errorCommand := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}
		successCommand := probe_command.KubesondeCommand{
			Destination:          "test-destination",
			DestinationPort:      "80",
			DestinationNamespace: "default",
			SourcePodName:        "test-pod",
			ContainerName:        "debugger",
			Namespace:            "default",
			Command:              "sample command",
			Action:               v1.ALLOW,
		}

		state := new(MockedCNIState)
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, errors.New("connection refused")).Once()
		state.On("runCommand", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(true, nil).Once()

		client := fake.NewSimpleClientset()
		p := &v12.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
			Status: v12.PodStatus{
				EphemeralContainerStatuses: []v12.ContainerStatus{
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
					{State: v12.ContainerState{Running: &v12.ContainerStateRunning{StartedAt: metav1.Time{}}}},
				},
			},
			Spec: v12.PodSpec{
				EphemeralContainers: []v12.EphemeralContainer{
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "debugger"}},
					{EphemeralContainerCommon: v12.EphemeralContainerCommon{Name: "monitor"}},
				},
			},
		}
		_, err := client.CoreV1().Pods("default").Create(context.TODO(), p, metav1.CreateOptions{})
		if err != nil {
			panic(err)
		}

		state.Mock.On("getClient").Return(client)

		output := InspectWithContinuousMode(state, []probe_command.KubesondeCommand{errorCommand, successCommand})

		// Error probe: 1 error + 1 DENY probe in Items
		Expect(len(output.Errors)).To(BeIdenticalTo(1))
		Expect(output.Errors[0].Reason).To(ContainSubstring("connection refused"))

		// Success probe: 1 ALLOW probe in Items
		// The DENY probe from the error path is deduplicated with the ALLOW probe (same source/dest/port)
		// so only the ALLOW probe should remain
		Expect(len(output.Items)).To(BeIdenticalTo(2))

		// Verify both outcomes are present
		var foundDeny, foundAllow int
		for _, item := range output.Items {
			if item.ResultingAction == v1.DENY {
				foundDeny++
			}
			if item.ResultingAction == v1.ALLOW {
				foundAllow++
			}
		}
		Expect(foundDeny).To(BeEquivalentTo(1))
		Expect(foundAllow).To(BeEquivalentTo(1))
	})
})
