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

package probe_command

import (
	"testing"

	"github.com/stretchr/testify/assert"
	kubesondev1 "kubesonde.io/api/v1"
)

func TestKubesondeCommand(t *testing.T) {
	t.Run("Test KubesondeCommand struct fields", func(t *testing.T) {
		// Create a command with all fields
		command := KubesondeCommand{
			Action:               kubesondev1.ALLOW,
			Command:              "test command",
			ContainerName:        "test-container",
			ProbeChecker:         func(string) bool { return true },
			Namespace:            "test-namespace",
			Protocol:             "TCP",
			Destination:          "test-destination",
			DestinationHostnames: []string{"host1", "host2"},
			DestinationIPAddress: "192.168.1.1",
			DestinationNamespace: "dest-namespace",
			DestinationPort:      "8080",
			DestinationLabels:    "label1=value1",
			DestinationType:      kubesondev1.POD,
			SourcePodName:        "source-pod",
			SourceIPAddress:      "192.168.1.2",
			SourceType:           kubesondev1.SERVICE,
			SourceLabels:         "label2=value2",
		}

		// Verify all fields are set correctly
		assert.Equal(t, kubesondev1.ALLOW, command.Action)
		assert.Equal(t, "test command", command.Command)
		assert.Equal(t, "test-container", command.ContainerName)
		assert.Equal(t, "test-namespace", command.Namespace)
		assert.Equal(t, "TCP", command.Protocol)
		assert.Equal(t, "test-destination", command.Destination)
		assert.Equal(t, []string{"host1", "host2"}, command.DestinationHostnames)
		assert.Equal(t, "192.168.1.1", command.DestinationIPAddress)
		assert.Equal(t, "dest-namespace", command.DestinationNamespace)
		assert.Equal(t, "8080", command.DestinationPort)
		assert.Equal(t, "label1=value1", command.DestinationLabels)
		assert.Equal(t, "source-pod", command.SourcePodName)
		assert.Equal(t, "192.168.1.2", command.SourceIPAddress)
		assert.Equal(t, "label2=value2", command.SourceLabels)
	})

	t.Run("Test KubesondeCommand with minimal fields", func(t *testing.T) {
		// Create a command with minimal fields
		command := KubesondeCommand{
			Command:         "minimal command",
			SourcePodName:   "minimal-pod",
			Destination:     "minimal-dest",
			DestinationPort: "80",
		}

		assert.Equal(t, "minimal command", command.Command)
		assert.Equal(t, "minimal-pod", command.SourcePodName)
		assert.Equal(t, "minimal-dest", command.Destination)
		assert.Equal(t, "80", command.DestinationPort)
		assert.Equal(t, kubesondev1.ActionType(""), command.Action) // Default value is empty string
		assert.Equal(t, "", command.ContainerName)                  // Default value
	})
}
