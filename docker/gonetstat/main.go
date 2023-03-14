package main

import (
	"encoding/json"
	"fmt"
	"strings"
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
	return NestatInfoRequestBodyItem{
		Fd:     int(data.UID),
		Family: 4,
		Type:   item_type,
		Laddr:  []string{data.LocalAddr.IP.String(), fmt.Sprint(data.LocalAddr.Port)},
		Raddr:  []string{data.RemoteAddr.IP.String(), fmt.Sprint(data.RemoteAddr.Port)},
		Status: "Open",
		Pid:    0,
	}
}

func display_socks() {
	var sockets []NestatInfoRequestBodyItem

	tcp_tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return !strings.Contains(s.LocalAddr.IP.String(), "127.0") && s.LocalAddr.IP.String() != "localhost" && s.LocalAddr.IP.String() != "::1" && s.State == netstat.Listen
	})
	if err != nil {
		return
	}
	for _, e := range tcp_tabs {
		sockets = append(sockets, toNetstatInfoRequestBodyItem(e, TCP_TYPE))
	}

	tcpv6_tabs, err := netstat.TCP6Socks(func(s *netstat.SockTabEntry) bool {
		return !strings.Contains(s.LocalAddr.IP.String(), "127.0") && s.LocalAddr.IP.String() != "localhost" && s.LocalAddr.IP.String() != "::1" && s.State == netstat.Listen
	})
	if err != nil {
		return
	}
	for _, e := range tcpv6_tabs {
		sockets = append(sockets, toNetstatInfoRequestBodyItem(e, TCP_TYPE))
	}

	udp_tabs, err := netstat.UDPSocks(func(s *netstat.SockTabEntry) bool {
		return !strings.Contains(s.LocalAddr.IP.String(), "127.0") && s.LocalAddr.IP.String() != "localhost" && s.LocalAddr.IP.String() != "::1"
	})
	if err != nil {
		return
	}
	for _, e := range udp_tabs {
		sockets = append(sockets, toNetstatInfoRequestBodyItem(e, UDP_TYPE))
	}

	udpv6_tabs, err := netstat.UDP6Socks(func(s *netstat.SockTabEntry) bool {
		return !strings.Contains(s.LocalAddr.IP.String(), "127.0") && s.LocalAddr.IP.String() != "localhost" && s.LocalAddr.IP.String() != "::1"
	})
	if err != nil {
		return
	}
	for _, e := range udpv6_tabs {
		sockets = append(sockets, toNetstatInfoRequestBodyItem(e, UDP_TYPE))
	}

	if len(sockets) == 0 {
		return
	}
	a, _ := json.Marshal(sockets)
	fmt.Printf("%s\n", a)
	return
}
func main() {
	display_socks()
	time.Sleep(3 * time.Second)
	return

}
