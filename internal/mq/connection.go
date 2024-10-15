package mq

import (
	"fmt"
	"net"
)

// Connects a socket for publishing messages to master node.
func ConnectPublish(s Socket, ip net.IP, port int) error {
	addr := fmt.Sprintf("tcp://[%s]:%d", ip, port)

	return s.Dial(addr)
}

// Connects a socket for subscribing messages from reporting nodes.
func ConnectSubscribe(s Socket, ip net.IP, port int) error {
	addr := fmt.Sprintf("tcp://[%s]:%d", ip, port)

	return s.Listen(addr)
}

func Close(s Socket) error {
	return s.Close()
}
