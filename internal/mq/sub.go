package mq

import (
	"fmt"
	"log"
	"time"

	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/domain"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Associates pub clients zeromq identity with their machine IDs.
// Key: MachineID
// Value: Socket Identity
var registeredPubSockets = map[string][]byte{}

func (s Socket) RegisterSubscriptions() {
	sc := di.ExtractSubscribeContainer(s.ctx)
	if sc == nil {
		logging.LogFatal("subscribe container hasn't been injected")
	}

	go func() {
		for {
			logging.LogInfo("waiting for messages...")
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
			logging.LogInfo("broadcasting stats polling duration update")
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
	sc *di.SubscribeContainer,
) {
	switch m.Topic() {
	case JoinNetwork:
		handleJoinNetworkRequest(s, m, sc.JoinNodesNetwork, sc.GetNodeStatsPollingDuration)
		return
	case UpdateNodeStats:
		handleUpdateNodeStatsRequest(m, sc.UpdateNodesNetwork)
		return
	default:
		logging.LogWarning("failed to understand message with topic %d", m.Topic())
		return
	}
}

func handleJoinNetworkRequest(
	s Socket,
	m Msg,
	join domain.JoinNodesNetwork,
	nodeStatsPollingDuration domain.GetNodeStatsPollingDuration,
) {
	logging.LogInfo("handling join network request")
	req, ok := m.Data().(JoinNetworkRequest)
	if !ok {
		err := fmt.Errorf("couldn't cast data to join network request, got: %v", m.Data)
		s.ReplyMsg(m.Identity(), compose(JoinNetwork, err))
		return
	}

	registeredPubSockets[req.Node.ID] = m.Identity()
	err := join(req.Node)
	if err != nil {
		logging.LogError("join node call failed, %v", err)
	}

	s.ReplyMsg(m.Identity(), compose(JoinNetwork, JoinNetworkResponse{StatsPoll: nodeStatsPollingDuration()}))
}

func handleUpdateNodeStatsRequest(
	m Msg,
	update domain.UpdateNodesNetwork,
) {
	log.Println("handling node update")
	req, ok := m.Data().(UpdateNodeStatsRequest)
	if !ok {
		logging.LogError("couldn't cast data to update stats request, got: %v", m.Data)
		return
	}

	err := update(req.Node)
	if err != nil {
		logging.LogError("updated node call failed, %v", err)
	}
}

func broadcastStatsPollingDurationUpdate(
	s Socket,
	d time.Duration,
) error {
	if len(registeredPubSockets) == 0 {
		logging.LogInfo("skipping stats polling duration update broadcast, no registered pub sockets")
		return nil
	}

	for _, identity := range registeredPubSockets {
		s.ReplyMsg(identity, compose(UpdateNodeStatsPollDuration, UpdateNodeStatsPollDurationRequest{Duration: d}))
	}

	return nil
}
