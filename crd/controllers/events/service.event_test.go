/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	v12 "kubesonde.io/api/v1"
)

func TestGetDestinationAndPort_ServiceFound(t *testing.T) {
	t.Run("Selector is populated when pod with matching port is found", func(t *testing.T) {
		service := v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-svc",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "10.0.0.1",
				Selector: map[string]string{
					"app": "myapp",
				},
				Ports: []v1.ServicePort{
					{
						Port:       80,
						Protocol:   v1.ProtocolTCP,
						TargetPort: intstr.FromInt(8080),
					},
				},
			},
		}

		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Ports: []v1.ContainerPort{
							{ContainerPort: 8080},
						},
					},
				},
			},
		}

		_, dst := getDestinationAndPort(service, []v1.Pod{pod}, service.Spec.Ports[0])

		assert.Equal(t, v12.SERVICE, dst.Type)
		assert.Equal(t, "test-pod", dst.Name)
		assert.Equal(t, "10.0.0.1", dst.IPAddress)
		assert.Equal(t, "app=myapp", dst.Selector)
	})
}

func TestGetDestinationAndPort_ServiceNotFound(t *testing.T) {
	t.Run("Selector is populated even when no pod matches", func(t *testing.T) {
		service := v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-svc",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "10.0.0.2",
				Selector: map[string]string{
					"app":     "myapp",
					"version": "v1",
				},
				Ports: []v1.ServicePort{
					{
						Port:       443,
						Protocol:   v1.ProtocolTCP,
						TargetPort: intstr.FromInt(8443),
					},
				},
			},
		}

		_, dst := getDestinationAndPort(service, []v1.Pod{}, service.Spec.Ports[0])

		assert.Equal(t, v12.SERVICE, dst.Type)
		assert.Contains(t, dst.Name, "Unknown - test-svc")
		assert.Equal(t, "10.0.0.2", dst.IPAddress)
		assert.Equal(t, "app=myapp;version=v1", dst.Selector)
	})
}

func TestGetDestinationAndPort_EmptySelector(t *testing.T) {
	t.Run("Selector is empty string when service has no selector", func(t *testing.T) {
		service := v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-svc",
				Namespace: "default",
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "10.0.0.3",
				Ports: []v1.ServicePort{
					{
						Port:       80,
						Protocol:   v1.ProtocolTCP,
						TargetPort: intstr.FromInt(80),
					},
				},
			},
		}

		pod := v1.Pod{
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Ports: []v1.ContainerPort{{ContainerPort: 80}}},
				},
			},
		}

		_, dst := getDestinationAndPort(service, []v1.Pod{pod}, service.Spec.Ports[0])

		assert.Equal(t, v12.SERVICE, dst.Type)
		assert.Equal(t, "", dst.Selector)
	})
}

func TestGetServicesAsProbes_ServiceType(t *testing.T) {
	t.Run("Service probes include selector on destination", func(t *testing.T) {
		service := v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "web-svc",
				Namespace: "production",
				Labels:    map[string]string{"tier": "frontend"},
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "10.0.1.1",
				Selector:  map[string]string{"app": "web"},
				Ports: []v1.ServicePort{
					{Port: 80, Protocol: v1.ProtocolTCP, TargetPort: intstr.FromInt(8080)},
				},
			},
		}

		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "web-pod", Namespace: "production"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Ports: []v1.ContainerPort{{ContainerPort: 8080}}},
				},
			},
		}

		items := getServicesAsProbes(service, []v1.Pod{pod})

		assert.Len(t, items, 1)
		assert.Equal(t, items[0].Destination.Type, v12.SERVICE)
		assert.Contains(t, items[0].Destination.Selector, "app=web")
		assert.Equal(t, "Internet", items[0].Source.Name)
	})

	t.Run("Service probes include both labels and selector on destination", func(t *testing.T) {
		service := v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "api-svc",
				Namespace: "staging",
				Labels: map[string]string{
					"tier": "api",
					"env":  "staging",
					"team": "backend",
				},
			},
			Spec: v1.ServiceSpec{
				ClusterIP: "10.0.2.1",
				Selector:  map[string]string{"app": "api", "version": "v2"},
				Ports: []v1.ServicePort{
					{Port: 8443, Protocol: v1.ProtocolTCP, TargetPort: intstr.FromInt(8443)},
				},
			},
		}

		pod := v1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: "api-pod", Namespace: "staging"},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{Ports: []v1.ContainerPort{{ContainerPort: 8443}}},
				},
			},
		}

		items := getServicesAsProbes(service, []v1.Pod{pod})

		assert.Len(t, items, 1)
		assert.Equal(t, v12.SERVICE, items[0].Destination.Type)
		assert.Equal(t, items[0].Destination.Selector, "app=api;version=v2")
		assert.Equal(t, items[0].Destination.Labels, "env=staging;team=backend;tier=api")

	})
}
