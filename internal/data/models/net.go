package models

import (
	"fmt"
	"net"

	nett "github.com/shirou/gopsutil/net"
)

type ConnectionKind uint32
type ConnectionState string
type IP []byte
type Address struct {
	IP   IP
	Port uint16
}

type Connection struct {
	Kind          ConnectionKind
	State         ConnectionState
	LocalAddress  Address
	RemoteAddress Address
}

func NewAddress(
	addr nett.Addr,
) Address {
	return Address{
		IP:   IP(net.ParseIP(addr.IP)),
		Port: uint16(addr.Port),
	}
}

func NewConnection(
	kind uint32,
	state string,
	localaddr nett.Addr,
	remoteaddr nett.Addr,
) Connection {
	return Connection{
		Kind:          ConnectionKind(kind),
		State:         ConnectionState(state),
		LocalAddress:  NewAddress(localaddr),
		RemoteAddress: NewAddress(remoteaddr),
	}
}

func (v IP) String() string {
	if len(v) < 4 {
		return "-"
	}

	return net.IP(v).String()
}

func (v ConnectionKind) String() string {
	switch v {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return "UNKNOWN"
	}
}

func (v Address) String() string {
	println(v.IP.String())
	return fmt.Sprintf("%s:%d", v.IP, v.Port)
}
