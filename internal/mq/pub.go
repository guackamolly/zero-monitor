package mq

import (
	"fmt"
	"log"

	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/domain"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

func (s Socket) RegisterPublishers() {
	pc := di.ExtractPublishContainer(s.ctx)
	if pc == nil {
		log.Fatalln("publish container hasn't been injected")
	}

	// Join network.
	go func() {
		n := pc.GetCurrentNode()
		err := s.PublishMsg(compose(JoinNetwork, JoinNetworkRequest{Node: n}))
		if err != nil {
			// TODO: handle join network error gracefully.
			logging.LogFatal("couldn't join network, %v", err)
		}
	}()

	// Handle sub reply messages.
	go func() {
		for {
			m, err := s.ReceiveMsg()
			if err != nil {
				logging.LogError("failed to receive message from sub socket, %v", err)
				continue
			}

			switch m.Topic {
			case JoinNetwork:
				handleJoinNetworkResponse(s, m, pc.StartNodeStatsPolling)
				continue
			case UpdateNodeStatsPollDuration:
				handleUpdateStatsPollDurationRequest(m, pc.UpdateNodeStatsPolling)
				continue
			case NodeConnections:
				handleNodeConnectionsRequest(s, pc.GetCurrentNodeConnections)
				continue
			default:
				logging.LogError("failed to recognize sub reply message, %v", m)
			}
		}
	}()
}

func handleJoinNetworkResponse(
	s Socket,
	m Msg,
	start domain.StartNodeStatsPolling,
) error {
	resp, ok := m.Data.(JoinNetworkResponse)
	if !ok {
		return handleUnknownMessage(m)
	}

	go func() {
		ns := start(resp.StatsPoll)
		for n := range ns {
			err := s.PublishMsg(compose(UpdateNodeStats, UpdateNodeStatsRequest{Node: n}))
			if err != nil {
				logging.LogError("failed to publish update stats message, %v", err)
			}
		}
	}()

	return nil
}

func handleUpdateStatsPollDurationRequest(
	m Msg,
	update domain.UpdateNodeStatsPolling,
) error {
	req, ok := m.Data.(UpdateNodeStatsPollDurationRequest)
	if !ok {
		return handleUnknownMessage(m)
	}

	update(req.Duration)
	return nil
}

func handleNodeConnectionsRequest(
	s Socket,
	connections domain.GetCurrentNodeConnections,
) error {
	conns, err := connections()
	if err != nil {
		return s.PublishMsg(compose(NodeConnections, err))
	}

	return s.PublishMsg(compose(NodeConnections, conns))
}

func handleUnknownMessage(
	m Msg,
) error {
	err, ok := m.Data.(error)
	if ok {
		return err
	}

	return fmt.Errorf("couldn't understand message, %v", m)
}
