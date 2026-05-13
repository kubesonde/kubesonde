package eventstorage

import (
	"fmt"
	"sync"

	"github.com/samber/lo"
	"kubesonde.io/controllers/probe_command"
)

var (
	commands   = make(map[string]probe_command.KubesondeCommand)
	commandsMu sync.RWMutex
)

func AddProbe(command probe_command.KubesondeCommand) {
	key := fmt.Sprintf("%s-%s-%s-%s-%s", command.SourcePodName, command.Command, command.DestinationIPAddress, command.DestinationPort, command.Protocol)
	commandsMu.RLock()
	_, ok := commands[key]
	commandsMu.RUnlock()
	if ok {
		return
	}
	commandsMu.Lock()
	commands[key] = command
	commandsMu.Unlock()
}

func AddProbes(probes []probe_command.KubesondeCommand) {
	for _, probe := range probes {
		AddProbe(probe)
	}
}

func GetProbes() []probe_command.KubesondeCommand {
	commandsMu.RLock()
	defer commandsMu.RUnlock()
	return lo.Values(commands)
}

func ProbeAvailable(command probe_command.KubesondeCommand) bool {
	key := fmt.Sprintf("%s-%s-%s-%s-%s", command.SourcePodName, command.Command, command.DestinationIPAddress, command.DestinationPort, command.Protocol)
	commandsMu.RLock()
	defer commandsMu.RUnlock()
	_, ok := commands[key]
	return ok
}
