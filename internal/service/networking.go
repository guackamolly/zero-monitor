package service

import (
	"net"

	"github.com/guackamolly/zero-monitor/internal/data/repositories"
)

// Service for networking operations.
type NetworkingService struct {
}

func NewNetworkingService() *NetworkingService {
	s := &NetworkingService{}
	return s
}

func (s NetworkingService) PublicIP() (net.IP, error) {
	return repositories.PublicIP()
}

func (s NetworkingService) PrivateIP() (net.IP, error) {
	ip, err := repositories.PrivateIP()
	if err != nil {
		ip, err = repositories.InterfaceIP()
	}

	return ip, err
}
