package mq

import (
	"fmt"
	"log"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/di"
	"github.com/guackamolly/zero-monitor/internal/service"
)

func (s Socket) RegisterSubscriptions() {
	sc := di.ExtractSubscribeContainer(s.ctx)
	if sc == nil {
		log.Fatalln("subscribe container hasn't been injected")
	}

	go func() {
		for {
			log.Println("waiting for messages...")

			m, err := s.Recv()
			if err != nil {
				log.Printf("failed to receive message, %v\n", err)
				continue
			}

			if l := len(m.Frames); l != 2 {
				log.Printf("received corrupted message, expected 2 frames but got: %d\n", l)
				continue
			}

			mm, err := decode(m.Frames[1])
			mm.Identity = m.Frames[0]
			if err != nil {
				log.Printf("failed to decode message, %v\n", err)
				continue
			}

			handle(s, mm, *sc)
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
	default:
		compose(unknown)
	}
}

func handleJoin(
	s Socket,
	m msg,
	nodeManager *service.NodeManagerService,
	masterConfiguration *service.MasterConfigurationService,
) {
	log.Println("handling node join")
	req, ok := m.Data.(joinNodeRequest)
	if !ok {
		err := fmt.Errorf("couldn't cast data to join node request, got: %v", m.Data)
		s.Reply(m.Identity, compose(xerror, err))
		return
	}

	err := nodeManager.Join(req.Node)
	if err != nil {
		log.Printf("join node call failed, %v\n", err)
	}

	cfg := masterConfiguration.Current()
	s.Reply(m.Identity, compose(reply, joinNodeResponse{StatsPoll: cfg.NodeStatsPolling.Duration()}))
}

func handleUpdate(
	m msg,
	service *service.NodeManagerService,
) {
	log.Println("handling node update")
	node, ok := m.Data.(models.Node)
	if !ok {
		log.Printf("couldn't cast data to Node model, got: %v\n", m.Data)
		return
	}

	err := service.Update(node)
	if err != nil {
		log.Printf("updated node call failed, %v\n", err)
	}
}
