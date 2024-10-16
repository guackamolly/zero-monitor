package service

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

// Service for managing nodes that report to master.
type NodeManagerService struct {
	stream    chan (map[string]models.Node)
	connected map[string]models.Node
}

func NewNodeManagerService() *NodeManagerService {
	return &NodeManagerService{
		stream:    make(chan map[string]models.Node),
		connected: map[string]models.Node{},
	}
}

// Joins master node.
func (s *NodeManagerService) Join(node models.Node) error {
	id := node.ID
	if _, ok := s.connected[id]; ok {
		return fmt.Errorf("node %s has already joined, ignorning request", id)
	}

	s.connected[id] = node
	s.updateStream()
	return nil
}

// Updates the state of a node.
func (s *NodeManagerService) Update(node models.Node) error {
	id := node.ID
	if _, ok := s.connected[id]; !ok {
		return fmt.Errorf("node %s hasn't joined yet, ignorning request", id)
	}

	s.connected[id] = node
	s.updateStream()
	return nil
}

func (s NodeManagerService) Connected() map[string]models.Node {
	return s.connected
}

func (s NodeManagerService) Stream() chan (map[string]models.Node) {
	return s.stream
}

func (s *NodeManagerService) updateStream() {
	go func() {
		s.stream <- s.connected
	}()
}
