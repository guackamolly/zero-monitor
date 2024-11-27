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
				err = handleJoinNetworkResponse(s, m, pc.StartNodeStatsPolling, pc.GetCurrentNode)
			case AuthenticateNetwork:
				err = handleAuthenticateNetworkResponse(s, m, pc.GetCurrentNode)
			case UpdateNodeStatsPollDuration:
				err = handleUpdateStatsPollDurationRequest(m, pc.UpdateNodeStatsPolling)
			case NodeConnections:
				err = handleNodeConnectionsRequest(s, m, pc.GetCurrentNodeConnections)
			case NodeProcesses:
				err = handleNodeProcessesRequest(s, m, pc.GetCurrentNodeProcesses)
			case NodePackages:
				err = handleNodePackagesRequest(s, m, pc.GetCurrentNodePackages)
			case KillNodeProcess:
				err = handleKillNodeProcessRequest(s, m, pc.KillNodeProcess)
			case StartNodeSpeedtest:
				err = handleStartNodeSpeedtestRequest(s, m, pc.StartNodeSpeedtest)
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
	currentNode domain.GetCurrentNode,
) error {
	if _, ok := m.Data.(RequiresAuthenticationResponse); ok {
		return handleRequiresAuthenticationResponse(s, currentNode)
	}

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

func handleAuthenticateNetworkResponse(
	s Socket,
	m Msg,
	currentNode domain.GetCurrentNode,
) error {
	if _, ok := m.Data.(AuthenticateNetworkResponse); !ok {
		logging.LogFatal("failed to authenticate within nodes network")
	}

	return s.PublishMsg(Compose(JoinNetwork, JoinNetworkRequest{currentNode()}))
}

func handleRequiresAuthenticationResponse(
	s Socket,
	currentNode domain.GetCurrentNode,
) error {
	// disallow handshaking more than once, otherwise both master and node will enter in a race condition like state
	if handshaked {
		logging.LogFatal("invalid state: already authenticated but master replied with <requires authentication>")
	}

	handshaked = true
	return s.PublishMsg(Compose(AuthenticateNetwork, AuthenticateNetworkRequest{InviteCode: InviteCode(), Node: currentNode()}))
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

func handleNodePackagesRequest(
	s Socket,
	m Msg,
	packages domain.GetCurrentNodePackages,
) error {
	pkgs, err := packages()
	if err != nil {
		return s.PublishMsg(m.WithError(err))
	}

	return s.PublishMsg(m.WithData(NodePackagesResponse{Packages: pkgs}))
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

	return s.PublishMsg(m.WithData(nil))
}

func handleStartNodeSpeedtestRequest(
	s Socket,
	m Msg,
	startNodeSpeedtest domain.StartNodeSpeedtest,
) error {
	ch, err := startNodeSpeedtest()
	if err != nil {
		return s.PublishMsg(m.WithError(err))
	}

	go func() {
		for st := range ch {
			err = s.PublishMsg(m.WithData(NodeSpeedtestResponse{Speedtest: st}))
			if err != nil {
				logging.LogError("failed to publish node speed test response, %v", err)
			}
		}
	}()

	return nil
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
