package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cakturk/go-netstat/netstat"
)

var (
	TCP_TYPE = 1
	UDP_TYPE = 2
)

type NestatInfoRequestBody []NestatInfoRequestBodyItem
type NestatInfoRequestBodyItem struct {
	Fd     int      `json:"fd"`
	Family int      `json:"family"`
	Type   int      `json:"type"`
	Laddr  []string `json:"laddr"`
	Raddr  []string `json:"raddr"`
	Status string   `json:"status"`
	Pid    int      `json:"pid"`
}

func toNetstatInfoRequestBodyItem(data netstat.SockTabEntry, item_type int) NestatInfoRequestBodyItem {
	pid := 0
	if data.Process != nil {
		pid = data.Process.Pid
	}

	return NestatInfoRequestBodyItem{
		Fd:     int(data.UID),
		Family: 4,
		Type:   item_type,
		Laddr:  []string{data.LocalAddr.IP.String(), fmt.Sprint(data.LocalAddr.Port)},
		Raddr:  []string{data.RemoteAddr.IP.String(), fmt.Sprint(data.RemoteAddr.Port)},
		Status: "Open",
		Pid:    pid,
	}
}

func isListeningTCP(s *netstat.SockTabEntry) bool {
	return !s.LocalAddr.IP.IsLoopback() && s.State == netstat.Listen
}

// UDP has no LISTEN state, so a bound, unconnected socket (no remote port)
// is the closest equivalent to "listening".
func isBoundUDP(s *netstat.SockTabEntry) bool {
	return !s.LocalAddr.IP.IsLoopback() && s.RemoteAddr.Port == 0
}

func display_socks() {
	var sockets []NestatInfoRequestBodyItem

	tables := []struct {
		name string
		fn   func(netstat.AcceptFn) ([]netstat.SockTabEntry, error)
		pred netstat.AcceptFn
		typ  int
	}{
		{"TCP", netstat.TCPSocks, isListeningTCP, TCP_TYPE},
		{"TCP6", netstat.TCP6Socks, isListeningTCP, TCP_TYPE},
		{"UDP", netstat.UDPSocks, isBoundUDP, UDP_TYPE},
		{"UDP6", netstat.UDP6Socks, isBoundUDP, UDP_TYPE},
	}

	for _, t := range tables {
		tabs, err := t.fn(t.pred)
		if err != nil {
			// Don't let one table's failure (e.g. no IPv6 support) drop the
			// results already collected from the others.
			fmt.Fprintf(os.Stderr, "gonetstat: failed to read %s sockets: %v\n", t.name, err)
			continue
		}
		for _, e := range tabs {
			sockets = append(sockets, toNetstatInfoRequestBodyItem(e, t.typ))
		}
	}

	if len(sockets) == 0 {
		return
	}
	a, _ := json.Marshal(sockets)
	fmt.Printf("%s\n", a)
}
func main() {
	display_socks()
	time.Sleep(3 * time.Second)
	return

}
