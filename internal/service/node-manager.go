package service

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

// Service for managing nodes that report to master.
type NodeManagerService struct {
	stream  chan ([]models.Node)
	network []models.Node
}

func NewNodeManagerService() *NodeManagerService {
	return &NodeManagerService{
		stream:  make(chan []models.Node),
		network: []models.Node{},
	}
}

// Joins master node.
func (s *NodeManagerService) Join(node models.Node) error {
	if i := s.nodeIdx(node); i >= 0 {
		return fmt.Errorf("node %s has already joined, ignorning request", node.ID)
	}

	s.network = append(s.network, node)
	s.updateStream()
	return nil
}

// Updates the state of a node.
func (s *NodeManagerService) Update(node models.Node) error {
	var i int
	if i = s.nodeIdx(node); i < 0 {
		return fmt.Errorf("node %s hasn't joined yet, ignorning request", node.ID)
	}

	s.network[i] = node
	s.updateStream()
	return nil
}

func (s NodeManagerService) Network() []models.Node {
	return s.network
}

func (s NodeManagerService) Stream() chan ([]models.Node) {
	return s.stream
}

func (s NodeManagerService) nodeIdx(node models.Node) int {
	for i, n := range s.network {
		if n.ID == node.ID {
			return i
		}
	}

	return -1
}

func (s *NodeManagerService) updateStream() {
	go func() {
		s.stream <- s.network
	}()
}
