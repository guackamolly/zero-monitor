package mq

import (
	"fmt"
	"log"
	"time"

	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/logging"
	"github.com/guackamolly/zero-monitor/internal/service"
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
			m, err := s.Receive()
			if err != nil {
				logging.LogError("failed to receive message from pub socket, %v", err)
				continue
			}

			handle(s, m, *sc)
		}
	}()

	go func() {
		mcs := sc.MasterConfiguration
		sp := mcs.Current().NodeStatsPolling
		ch := sc.MasterConfiguration.Stream()

		for cfg := range ch {
			if sp == cfg.NodeStatsPolling {
				continue
			}
			sp = cfg.NodeStatsPolling

			logging.LogInfo("broadcasting stats polling duration update")
			err := broadcastStatsPollingDurationUpdate(s, sp.Duration())
			if err != nil {
				logging.LogError("failed to broadcast stats polling duration update, %v", err)
			}
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

	registeredPubSockets[req.Node.ID] = m.Identity
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

func broadcastStatsPollingDurationUpdate(
	s Socket,
	d time.Duration,
) error {
	if len(registeredPubSockets) == 0 {
		logging.LogInfo("skipping stats polling duration update broadcast, no registered pub sockets")
		return nil
	}

	for _, identity := range registeredPubSockets {
		s.Reply(identity, compose(updateStatsPollDuration, updateStatsPollDurationRequest{Duration: d}))
	}

	return nil
}
