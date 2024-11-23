package service

import (
	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/event"
)

// Service for managing network nodes.
type NetworkService struct {
	subscriber event.EventSubscriber

	code *models.JoinNetworkCode
}

func NewNetworkService(
	subscriber event.EventSubscriber,
) *NetworkService {
	s := &NetworkService{
		subscriber: subscriber,
	}

	return s
}

func (s *NetworkService) Code() models.JoinNetworkCode {
	if s.code != nil && !s.code.Expired() {
		return *s.code
	}

	code := models.NewJoinNetworkCode()
	s.code = &code

	return code
}

func (s *NetworkService) Valid(code string) bool {
	if s.code == nil || s.code.Expired() {
		return false
	}

	return s.code.Code == code
}

// todo: cache public key
func (s *NetworkService) PublicKey() ([]byte, error) {
	return s.subscriber.PublicKey()
}

// todo: cache address
func (s *NetworkService) Address() models.Address {
	return s.subscriber.Address()
}
