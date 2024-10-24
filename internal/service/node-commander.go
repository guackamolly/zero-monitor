package service

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type Command byte

const (
	connections Command = iota + 1
)

// Service for executing commands in nodes.
type NodeCommanderService struct {
	execute func(id string, command Command) (any, error)
}

func NewNodeCommanderService(
	execute func(id string, command Command) (any, error),
) *NodeCommanderService {
	s := &NodeCommanderService{
		execute: execute,
	}

	return s
}

func (s NodeCommanderService) Connections(id string) ([]models.Connection, error) {
	r, err := s.execute(id, connections)
	if err != nil {
		return nil, err
	}

	conns, ok := r.([]models.Connection)
	if !ok {
		return nil, fmt.Errorf("failed to parse %v to connection slice", r)
	}

	return conns, nil
}
