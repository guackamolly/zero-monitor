package service_test

import (
	"testing"

	"github.com/guackamolly/zero-monitor/internal/service"
)

func TestNetworkServiceCodeReturnsLastCodeCreated(t *testing.T) {
	s := service.NewNetworkService()

	// create new code
	c1 := s.Code()
	c2 := s.Code()

	if c1 != c2 {
		t.Errorf("expected %v, but got %v", c1, c2)
	}
}
