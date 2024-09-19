package eventstorage

import (
	"fmt"

	"github.com/samber/lo"
	"kubesonde.io/controllers/probe_command"
)

var (
	commands = make(map[string]probe_command.KubesondeCommand)
)

func AddProbe(command probe_command.KubesondeCommand) {
	key := fmt.Sprintf("%s-%s-%s-%s-%s", command.SourcePodName, command.Command, command.DestinationIPAddress, command.DestinationPort, command.Protocol)
	_, ok := commands[key]
	if ok {
		return
	} else {
		commands[key] = command
	}
}

func AddProbes(probes []probe_command.KubesondeCommand) {

	for _, probe := range probes {
		AddProbe(probe)
	}
}
func GetProbes() []probe_command.KubesondeCommand {
	return lo.Values(commands)
}

func ProbeAvailable(command probe_command.KubesondeCommand) bool {
	key := fmt.Sprintf("%s-%s-%s-%s-%s", command.SourcePodName, command.Command, command.DestinationIPAddress, command.DestinationPort, command.Protocol)
	_, ok := commands[key]
	return ok
}
