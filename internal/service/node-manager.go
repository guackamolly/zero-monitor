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
	code    *models.JoinNetworkCode
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
	if !s.IsAuthenticated(node) {
		return fmt.Errorf("node %s hasn't authenticated yet", node.ID)
	}

	s.updateStream()
	return nil
}

// Updates the state of a node.
func (s *NodeManagerService) Update(node models.Node) error {
	if !s.IsAuthenticated(node) {
		return fmt.Errorf("node %s hasn't authenticated yet", node.ID)
	}

	s.network[s.nodeIdx(node)] = node
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

// Removes a node from the network.
func (s *NodeManagerService) Remove(node models.Node) error {
	if !s.IsAuthenticated(node) {
		return fmt.Errorf("node %s does not exist on the netwrok", node.ID)
	}

	network := []models.Node{}
	s.lock.Lock()
	for _, n := range s.network {
		if n.ID != node.ID {
			network = append(network, n)
		}
	}

	s.network = network
	s.lock.Unlock()

	s.updateStream()
	return nil
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

func (s *NodeManagerService) Code() models.JoinNetworkCode {
	if s.code != nil && !s.code.Expired() {
		return *s.code
	}

	code := models.NewJoinNetworkCode()
	s.code = &code

	return code
}

func (s *NodeManagerService) Valid(code string) bool {
	if s.code == nil || s.code.Expired() {
		return false
	}

	return s.code.Code == code
}

func (s *NodeManagerService) Authenticate(node models.Node, code string) error {
	if s.IsAuthenticated(node) {
		return fmt.Errorf("already authenticated")
	}

	if !s.Valid(code) {
		return fmt.Errorf("invalid code")

	}

	s.network = append(s.network, node)
	return nil
}

func (s *NodeManagerService) IsAuthenticated(node models.Node) bool {
	return s.nodeIdx(node) >= 0
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
