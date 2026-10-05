//go:build linux

package main

import (
	"net"
	"testing"

	"github.com/cakturk/go-netstat/netstat"
)

func entry(ip string, localPort, remotePort uint16, state netstat.SkState) *netstat.SockTabEntry {
	return &netstat.SockTabEntry{
		LocalAddr:  &netstat.SockAddr{IP: net.ParseIP(ip), Port: localPort},
		RemoteAddr: &netstat.SockAddr{IP: net.ParseIP("0.0.0.0"), Port: remotePort},
		State:      state,
	}
}

func TestIsBoundUDP(t *testing.T) {
	cases := []struct {
		name string
		e    *netstat.SockTabEntry
		want bool
	}{
		{"bound non-loopback socket is a listener", entry("0.0.0.0", 9125, 0, 0), true},
		{"connected socket (has a remote port) is not a listener", entry("0.0.0.0", 9125, 5000, 0), false},
		{"loopback socket is excluded", entry("127.0.0.1", 9125, 0, 0), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isBoundUDP(c.e); got != c.want {
				t.Errorf("isBoundUDP() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestIsListeningTCP(t *testing.T) {
	cases := []struct {
		name string
		e    *netstat.SockTabEntry
		want bool
	}{
		{"listening non-loopback socket", entry("0.0.0.0", 80, 0, netstat.Listen), true},
		{"established socket is not listening", entry("0.0.0.0", 80, 0, 1), false},
		{"loopback listener is excluded", entry("127.0.0.1", 80, 0, netstat.Listen), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := isListeningTCP(c.e); got != c.want {
				t.Errorf("isListeningTCP() = %v, want %v", got, c.want)
			}
		})
	}
}
