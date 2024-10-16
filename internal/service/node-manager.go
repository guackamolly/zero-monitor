package service

import (
	"fmt"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Established duration to check if any of the network nodes
// are currently offline or not.
var offlineTimeout = time.Second * 10

// Service for managing nodes that report to master.
type NodeManagerService struct {
	stream  chan ([]models.Node)
	network []models.Node
}

func NewNodeManagerService() *NodeManagerService {
	s := &NodeManagerService{
		stream:  make(chan []models.Node),
		network: []models.Node{},
	}

	go func() {
		for {
			time.Sleep(offlineTimeout)
			t := time.Now()
			for _, n := range s.network {
				if n.LastSeen.Sub(t).Abs() < offlineTimeout {
					continue
				}

				err := s.Update(n.SetOffline())
				if err != nil {
					logging.LogError("very strange error when notifying network that node is offline, %v", err)
				}
			}
		}
	}()

	return s
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
