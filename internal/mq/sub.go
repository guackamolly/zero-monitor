package mq

import (
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

			if l := len(m.Frames); l != 3 {
				log.Printf("received corrupted message, expected 3 frames but got: %d\n", l)
			}

			msg, err := decode(m.Frames[2])
			if err != nil {
				log.Printf("failed to decode message, %v\n", err)
				continue
			}

			handle(msg, *sc)
		}
	}()
}

func handle(
	m msg,
	serviceContainer di.SubscribeContainer,
) {
	if err, ok := m.Data.(error); ok {
		log.Printf("received err %v for topic %d\n", err, m.Topic)
		return
	}

	switch m.Topic {
	case join:
		handleJoin(m, serviceContainer.NodeManager)
		return
	case update:
		handleUpdate(m, serviceContainer.NodeManager)
		return
	default:
		compose(unknown)
	}
}

func handleJoin(
	m msg,
	service *service.NodeManagerService,
) {
	log.Println("handling node join")
	node, ok := m.Data.(models.Node)
	if !ok {
		log.Printf("couldn't cast data to Node model, got: %v\n", m.Data)
		return
	}

	err := service.Join(node)
	if err != nil {
		log.Printf("join node call failed, %v\n", err)
	}
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
