package mq

import (
	"fmt"
	"log"

	"github.com/guackamolly/zero-monitor/internal/domain"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

func (s Socket) RegisterPublishers() {
	pc := ExtractPublishContainer(s.ctx)
	if pc == nil {
		log.Fatalln("publish container hasn't been injected")
	}

	// Join network.
	go func() {
		n := pc.GetCurrentNode()
		err := s.PublishMsg(Compose(JoinNetwork, JoinNetworkRequest{Node: n}))
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
			topic := m.Topic

			logging.LogInfo("(pub) handling topic: %d", topic)
			switch topic {
			case JoinNetwork:
				err = handleJoinNetworkResponse(s, m, pc.StartNodeStatsPolling)
			case UpdateNodeStatsPollDuration:
				err = handleUpdateStatsPollDurationRequest(m, pc.UpdateNodeStatsPolling)
			case NodeConnections:
				err = handleNodeConnectionsRequest(s, m, pc.GetCurrentNodeConnections)
			case NodeProcesses:
				err = handleNodeProcessesRequest(s, m, pc.GetCurrentNodeProcesses)
			case KillNodeProcess:
				err = handleKillNodeProcessRequest(s, m, pc.KillNodeProcess)
			default:
				err = fmt.Errorf("failed to recognize sub reply message, %v", m)
			}

			if err != nil {
				logging.LogError("(pub) failed to handle topic %d, %v", topic, err)
				err = nil
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
			err := s.PublishMsg(Compose(UpdateNodeStats, UpdateNodeStatsRequest{Node: n}))
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
	m Msg,
	connections domain.GetCurrentNodeConnections,
) error {
	conns, err := connections()
	if err != nil {
		return s.PublishMsg(m.WithError(err))
	}

	return s.PublishMsg(m.WithData(NodeConnectionsResponse{Connections: conns}))
}

func handleNodeProcessesRequest(
	s Socket,
	m Msg,
	processes domain.GetCurrentNodeProcesses,
) error {
	procs, err := processes()
	if err != nil {
		return s.PublishMsg(m.WithError(err))
	}

	return s.PublishMsg(m.WithData(NodeProcessesResponse{Processes: procs}))
}

func handleKillNodeProcessRequest(
	s Socket,
	m Msg,
	killNodeProcess domain.KillNodeProcess,
) error {
	req, ok := m.Data.(KillNodeProcessRequest)
	if !ok {
		return fmt.Errorf("couldn't cast data to KillNodeProcessRequest, %v", m.Data)
	}

	err := killNodeProcess(req.PID)
	if err != nil {
		return s.PublishMsg(m.WithError(err))
	}

	return s.PublishMsg(m)
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
