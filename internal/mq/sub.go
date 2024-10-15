package mq

import (
	"fmt"
	"log"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/service"
)

func (s Socket) RegisterSubscriptions() {
	sc := ExtractSubscribeContainer(s.ctx)
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

			msg, err := decode(m.Bytes())
			if err != nil {
				log.Printf("failed to decode message, %v\n", err)
				s.Reply(compose(reply, err))

				continue
			}

			s.Reply(handle(msg, *sc))
		}
	}()
}

func handle(
	m msg,
	serviceContainer SubscribeContainer,
) msg {
	if err, ok := m.data.(error); ok {
		log.Printf("received err %v for topic %d\n", err, m.topic)

		return compose(empty)
	}

	switch m.topic {
	case join:
		return handleJoin(m, serviceContainer.NodeManager)
	case update:
		return handleUpdate(m, serviceContainer.NodeManager)
	default:
		return compose(unknown)
	}
}

func handleJoin(
	m msg,
	service *service.NodeManagerService,
) msg {
	node, ok := m.data.(models.Node)
	if !ok {
		return compose(reply, fmt.Errorf("couldn't cast data to Node model, got: %v", m.data))
	}

	err := service.Join(node)
	if err != nil {
		log.Printf("join node call failed, %v\n", err)
	}

	return compose(reply, err)
}

func handleUpdate(
	m msg,
	service *service.NodeManagerService,
) msg {
	node, ok := m.data.(models.Node)
	if !ok {
		return compose(reply, fmt.Errorf("couldn't cast data to Node model, got: %v", m.data))
	}

	err := service.Join(node)
	if err != nil {
		log.Printf("join node call failed, %v\n", err)
	}

	return compose(reply, err)
}
