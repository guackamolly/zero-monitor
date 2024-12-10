package mq

import (
	"fmt"
	"log"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/domain"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Associated pub clients identity to node models.
var clientNodes = map[string]models.Node{}

func (s Socket) RegisterSubscriptions() {
	sc := ExtractSubscribeContainer(s.ctx)
	if sc == nil {
		logging.LogFatal("subscribe container hasn't been injected")
	}

	go func() {
		for {
			logging.LogDebug("waiting for messages...")
			m, err := s.ReceiveMsg()
			if err != nil {
				logging.LogError("failed to receive message from pub socket, %v", err)
				continue
			}

			handle(s, m, sc)
		}
	}()

	go func() {
		spu := sc.GetNodeStatsPollingDurationUpdates()

		for sp := range spu {
			logging.LogDebug("broadcasting stats polling duration update")
			err := broadcastStatsPollingDurationUpdate(s, sp)
			if err != nil {
				logging.LogError("failed to broadcast stats polling duration update, %v", err)
			}
		}
	}()
}

func handle(
	s Socket,
	m Msg,
	sc *SubscribeContainer,
) {
	logging.LogDebug("(sub) handling topic: %d", m.Topic)
	switch m.Topic {
	case HelloNetwork:
		handleHelloNetworkRequest(s, m)
		return
	case JoinNetwork:
		handleJoinNetworkRequest(s, m, sc.JoinNodesNetwork, sc.RequiresNodesNetworkAuthentication, sc.GetNodeStatsPollingDuration)
		return
	case AuthenticateNetwork:
		handleAuthenticateNetworkRequest(s, m, sc.AuthenticateNodesNetwork)
		return
	case UpdateNodeStats:
		handleUpdateNodeStatsRequest(m, sc.UpdateNodesNetwork)
		return
	default:
		logging.LogWarning("failed to understand message with topic %d", m.Topic)
		return
	}
}

func handleHelloNetworkRequest(
	s Socket,
	m Msg,
) {
	logging.LogDebug("handling hello network request")
	logging.LogInfo("%x wants to participate in the network", m.Identity)

	s.ReplyMsg(m.Identity, Compose(HelloNetwork))
}

func handleJoinNetworkRequest(
	s Socket,
	m Msg,
	join domain.JoinNodesNetwork,
	requiresAuthentication domain.RequiresNodesNetworkAuthentication,
	nodeStatsPollingDuration domain.GetNodeStatsPollingDuration,
) {
	logging.LogDebug("handling join network request")
	req, ok := m.Data.(JoinNetworkRequest)
	if !ok {
		err := fmt.Errorf("couldn't cast data to join network request, got: %v", m.Data)
		s.ReplyMsg(m.Identity, Compose(JoinNetwork, err))
		return
	}

	if requiresAuthentication(req.Node) {
		s.ReplyMsg(m.Identity, Compose(JoinNetwork, RequiresAuthenticationResponse{}))
		return
	}

	s.Clients[req.Node.ID] = m.Identity
	clientNodes[string(m.Identity)] = req.Node

	err := join(req.Node)
	if err != nil {
		logging.LogError("join node call failed, %v", err)
	}

	s.ReplyMsg(m.Identity, Compose(JoinNetwork, JoinNetworkResponse{StatsPoll: nodeStatsPollingDuration()}))
}

func handleAuthenticateNetworkRequest(
	s Socket,
	m Msg,
	authenticate domain.AuthenticateNodesNetwork,
) {
	log.Println("handling authenticate network request...")
	req, ok := m.Data.(AuthenticateNetworkRequest)
	if !ok {
		logging.LogError("couldn't cast data to authenticate network request, got: %v", m.Data)
		return
	}

	err := authenticate(req.Node, req.InviteCode)
	if err != nil {
		s.ReplyMsg(m.Identity, Compose(AuthenticateNetwork).WithError(err))
		return
	}

	s.ReplyMsg(m.Identity, Compose(AuthenticateNetwork, AuthenticateNetworkResponse{}))
}

func handleUpdateNodeStatsRequest(
	m Msg,
	update domain.UpdateNodesNetwork,
) {
	log.Println("handling node update")
	req, ok := m.Data.(UpdateNodeStatsRequest)
	if !ok {
		logging.LogError("couldn't cast data to update stats request, got: %v", m.Data)
		return
	}

	node, ok := clientNodes[string(m.Identity)]
	if !ok {
		logging.LogWarning("received updated node stats request from client that is not registered in the network. possibly impersonating!")
		return
	}

	node = node.WithUpdatedStats(req.Stats)
	clientNodes[string(m.Identity)] = node

	err := update(node)
	if err != nil {
		logging.LogError("updated node call failed, %v", err)
	}
}

func broadcastStatsPollingDurationUpdate(
	s Socket,
	d time.Duration,
) error {
	if len(s.Clients) == 0 {
		logging.LogWarning("skipping stats polling duration update broadcast, no registered pub sockets")
		return nil
	}

	for _, identity := range s.Clients {
		s.ReplyMsg(identity, Compose(UpdateNodeStatsPollDuration, UpdateNodeStatsPollDurationRequest{Duration: d}))
	}

	return nil
}
