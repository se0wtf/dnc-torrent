package internal

import "net"

type Peer struct {
	IP net.IP
	Port uint16
}