package service

import (
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/event"
)

// Service for managing network nodes.
type NetworkService struct {
	subscriber event.EventSubscriber
}

func NewNetworkService(
	subscriber event.EventSubscriber,
) *NetworkService {
	s := &NetworkService{
		subscriber: subscriber,
	}

	return s
}

// todo: cache public key
func (s *NetworkService) PublicKey() ([]byte, error) {
	return s.subscriber.PublicKey()
}

// todo: cache address
func (s *NetworkService) Address() models.Address {
	return s.subscriber.Address()
}
