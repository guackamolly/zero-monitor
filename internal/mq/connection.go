package mq

import (
	"fmt"
	"net"

	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Connects a socket for publishing messages to master node.
func ConnectPublish(s *Socket, host string, port string) error {
	s.Endpoint = fmt.Sprintf("tcp://[%s]:%s", host, port)
	return s.Dial(s.Endpoint)
}

// Connects a socket for subscribing messages from reporting nodes.
func ConnectSubscribe(s *Socket, host string, port string) error {
	s.Endpoint = fmt.Sprintf("tcp://[%s]:%s", lookupHost(host), port)
	return s.Listen(s.Endpoint)
}

func Close(s Socket) error {
	return s.Close()
}

func lookupHost(host string) net.IP {
	if ip := net.ParseIP(host); ip != nil {
		return ip
	}

	ip, err := repositories.PrivateIP()
	if err != nil {
		logging.LogWarning("couldn't find the private ip of the closest network interface card. defaulting to any address, %v", err)
		return net.IPv4(0, 0, 0, 0)
	}

	return net.IP(ip)
}
