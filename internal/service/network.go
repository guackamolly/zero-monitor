package service

import "github.com/guackamolly/zero-monitor/internal/data/models"

// Service for managing network nodes.
type NetworkService struct {
	code *models.JoinNetworkCode
}

func NewNetworkService() *NetworkService {
	s := &NetworkService{}

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
