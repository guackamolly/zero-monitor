package mq

import (
	"fmt"
	"log"

	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/service"
)

func (s Socket) RegisterSubscriptions() {
	sc := di.ExtractSubscribeContainer(s.ctx)
	if sc == nil {
		logging.LogFatal("subscribe container hasn't been injected")
	}

	go func() {
		for {
			logging.LogInfo("waiting for messages...")
			m, err := s.Receive()
			if err != nil {
				logging.LogError("failed to receive message from pub socket, %v", err)
				continue
			}

			handle(s, m, *sc)
		}
	}()
}

func handle(
	s Socket,
	m msg,
	serviceContainer di.SubscribeContainer,
) {
	if err, ok := m.Data.(error); ok {
		log.Printf("received err %v for topic %d\n", err, m.Topic)
		return
	}

	switch m.Topic {
	case join:
		handleJoin(s, m, serviceContainer.NodeManager, serviceContainer.MasterConfiguration)
		return
	case update:
		handleUpdate(m, serviceContainer.NodeManager)
		return
	}
}

func handleJoin(
	s Socket,
	m msg,
	nodeManager *service.NodeManagerService,
	masterConfiguration *service.MasterConfigurationService,
) {
	logging.LogInfo("handling node join")
	req, ok := m.Data.(joinRequest)
	if !ok {
		err := fmt.Errorf("couldn't cast data to join request, got: %v", m.Data)
		s.Reply(m.Identity, compose(join, err))
		return
	}

	err := nodeManager.Join(req.Node)
	if err != nil {
		log.Printf("join node call failed, %v\n", err)
	}

	cfg := masterConfiguration.Current()
	s.Reply(m.Identity, compose(join, joinResponse{StatsPoll: cfg.NodeStatsPolling.Duration()}))
}

func handleUpdate(
	m msg,
	nodeManager *service.NodeManagerService,
) {
	log.Println("handling node update")
	req, ok := m.Data.(updateStatsRequest)
	if !ok {
		logging.LogError("couldn't cast data to update stats request, got: %v", m.Data)
		return
	}

	err := nodeManager.Update(req.Node)
	if err != nil {
		logging.LogError("updated node call failed, %v", err)
	}
}
