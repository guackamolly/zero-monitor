package service

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

// Service for managing nodes that report to master.
type NodeManagerService struct {
	connected map[string]models.Node
}

// Joins master node.
func (s *NodeManagerService) Join(node models.Node) error {
	id := node.ID
	if _, ok := s.connected[id]; ok {
		return fmt.Errorf("node %s has already joined, ignorning request", id)
	}

	s.connected[id] = node
	return nil
}

// Updates the state of a node.
func (s *NodeManagerService) Update(node models.Node) error {
	id := node.ID
	if _, ok := s.connected[id]; !ok {
		return fmt.Errorf("node %s hasn't joined yet, ignorning request", id)
	}

	s.connected[id] = node
	return nil
}
