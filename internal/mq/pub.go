package mq

import (
	"fmt"
	"log"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/service"
)

func (s Socket) RegisterPublishers() {
	pc := di.ExtractPublishContainer(s.ctx)
	if pc == nil {
		log.Fatalln("publish container hasn't been injected")
	}

	// Join network.
	go func() {
		nr := pc.NodeReporter
		err := joinNetwork(s, nr.Initial())
		if err != nil {
			// TODO: handle join network error gracefully.
			logging.LogFatal("couldn't join network, %v", err)
		}
	}()

	// Handle sub reply messages.
	go func() {
		for {
			m, err := s.Receive()
			if err != nil {
				logging.LogError("failed to receive message from sub socket, %v", err)
				continue
			}

			switch m.Topic {
			case join:
				handleJoinNetworkResponse(s, m, pc.NodeReporter)
				continue
			default:
				logging.LogError("failed to recognize sub reply message, %v", m)
			}
		}
	}()
}

func joinNetwork(
	s Socket,
	node models.Node,
) error {
	m := compose(join, joinRequest{Node: node})

	return s.Publish(m)
}

func updateStats(
	s Socket,
	node models.Node,
) error {
	m := compose(update, updateStatsRequest{Node: node})

	return s.Publish(m)
}

func handleJoinNetworkResponse(
	s Socket,
	m msg,
	nr *service.NodeReporterService,
) error {
	resp, ok := m.Data.(joinResponse)
	if !ok {
		return handleUnknownMessage(m)
	}

	go func() {
		ns := nr.Start(resp.StatsPoll)
		for n := range ns {
			err := updateStats(s, n)
			if err != nil {
				logging.LogError("failed to publish update stats message, %v", err)
			}
		}
	}()

	return nil
}

func handleUnknownMessage(
	m msg,
) error {
	err, ok := m.Data.(error)
	if ok {
		return err
	}

	return fmt.Errorf("couldn't understand message, %v", m)
}
