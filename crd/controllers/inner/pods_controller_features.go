package inner

import (
	. "k8s.io/api/core/v1"
	v12 "kubesonde.io/api/v1"
)

var podWithNoOpenPorts = Pod{Spec: PodSpec{
	// For simplicity, we only fill out the required fields.
	Containers: []Container{
		{
			Name:  "test-container",
			Image: "test-image",
			Ports: []ContainerPort{},
		},
		{
			Name:  "test-container-2",
			Image: "test-image-2",
			Ports: []ContainerPort{},
		},
	},
	RestartPolicy: RestartPolicyOnFailure,
}}

var podWithOpenPorts = &Pod{Spec: PodSpec{
	Containers: []Container{
		{
			Name:  "test-container",
			Image: "test-image",
			Ports: []ContainerPort{
				{
					ContainerPort: 111,
				},
				{
					ContainerPort: 112,
				},
			},
		},
		{
			Name:  "test-container-2",
			Image: "test-image-2",
			Ports: []ContainerPort{
				{
					ContainerPort: 221,
				},
				{
					ContainerPort: 222,
				},
			},
		},
	},
	RestartPolicy: RestartPolicyOnFailure,
}}

var probingActions = []v12.ProbingAction{
	{
		Action:          "Allow",
		FromPodSelector: "test-src-pod",
		ToPodSelector:   "test-dest-pod",
		Port:            "123",
	},
	{
		Action:          "Deny",
		FromPodSelector: "test-src-pod",
		Url:             "http://example.website.com",
	},
	{
		Action:          "Deny",
		FromPodSelector: "test-src-pod",
		Url:             "http://example.website.com",
		Endpoint:        "api/healthz",
	},
}

var buildTestPod = func(containers []Container, ip string) Pod {
	return Pod{
		Status: PodStatus{Phase: PodRunning, Conditions: nil, Message: "", Reason: "", NominatedNodeName: "", HostIP: ip, PodIP: ip,
			PodIPs: []PodIP{}, StartTime: nil, InitContainerStatuses: nil, ContainerStatuses: nil, QOSClass: PodQOSBestEffort, EphemeralContainerStatuses: nil},
		Spec: PodSpec{
			Volumes:                       nil,
			InitContainers:                nil,
			Containers:                    containers,
			EphemeralContainers:           nil,
			RestartPolicy:                 "",
			TerminationGracePeriodSeconds: nil,
			ActiveDeadlineSeconds:         nil,
			DNSPolicy:                     "",
			NodeSelector:                  nil,
			ServiceAccountName:            "",
			DeprecatedServiceAccount:      "",
			AutomountServiceAccountToken:  nil,
			NodeName:                      "",
			HostNetwork:                   false,
			HostPID:                       false,
			HostIPC:                       false,
			ShareProcessNamespace:         nil,
			SecurityContext:               nil,
			ImagePullSecrets:              nil,
			Hostname:                      "",
			Subdomain:                     "",
			Affinity:                      nil,
			SchedulerName:                 "",
			Tolerations:                   nil,
			HostAliases:                   nil,
			PriorityClassName:             "",
			Priority:                      nil,
			DNSConfig:                     nil,
			ReadinessGates:                nil,
			RuntimeClassName:              nil,
			EnableServiceLinks:            nil,
			PreemptionPolicy:              nil,
			Overhead:                      nil,
			TopologySpreadConstraints:     nil,
		},
	}
}

var buildContainers = func(ports []int32) Container {
	var containerPorts []ContainerPort
	for _, port := range ports {
		containerPorts = append(containerPorts,
			ContainerPort{
				Name:          "",
				HostPort:      port,
				ContainerPort: port,
				Protocol:      "",
				HostIP:        "",
			})
	}
	return Container{
		Name:                     "",
		Image:                    "",
		Command:                  nil,
		Args:                     nil,
		WorkingDir:               "",
		Ports:                    containerPorts,
		EnvFrom:                  nil,
		Env:                      nil,
		Resources:                ResourceRequirements{},
		VolumeMounts:             nil,
		VolumeDevices:            nil,
		LivenessProbe:            nil,
		ReadinessProbe:           nil,
		StartupProbe:             nil,
		Lifecycle:                nil,
		TerminationMessagePath:   "",
		TerminationMessagePolicy: "",
		ImagePullPolicy:          "",
		SecurityContext:          nil,
		Stdin:                    false,
		StdinOnce:                false,
		TTY:                      false,
	}
}
