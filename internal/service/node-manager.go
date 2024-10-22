package service

import (
	"fmt"
	"sync"

	"github.com/guackamolly/zero-monitor/internal/data/models"
)

// Service for managing nodes that report to master.
type NodeManagerService struct {
	streams []chan ([]models.Node)
	network []models.Node
	lock    *sync.RWMutex
}

// Creates a service for managing network nodes. Allows passing
// previously saved nodes as varadiac param.
func NewNodeManagerService(nodes ...models.Node) *NodeManagerService {
	s := &NodeManagerService{
		streams: []chan []models.Node{},
		network: []models.Node(nodes),
		lock:    &sync.RWMutex{},
	}

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

func (s NodeManagerService) Node(id string) (models.Node, bool) {
	for _, n := range s.network {
		if n.ID == id {
			return n, true
		}
	}

	return models.Node{}, false
}

func (s NodeManagerService) Network() []models.Node {
	return s.network
}

func (s *NodeManagerService) Stream() chan ([]models.Node) {
	stream := make(chan ([]models.Node))

	s.lock.Lock()
	s.streams = append(s.streams, stream)
	s.lock.Unlock()

	return stream
}

// Notifies the manager that a stream should be released.
func (s *NodeManagerService) Release(stream chan ([]models.Node)) {
	s.lock.Lock()

	l := len(s.streams)
	j := 0
	sc := make([]chan ([]models.Node), l-1)
	for i := 0; i < l; i++ {
		if s.streams[i] == stream {
			continue
		}

		sc[j] = s.streams[i]
		j++
	}

	s.streams = sc
	s.lock.Unlock()

	close(stream)
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
		for _, st := range s.streams {
			st <- s.network
		}
	}()
}
