package service

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/event"
)

// Service for executing commands in nodes.
type NodeCommanderService struct {
	publisher  event.EventPublisher
	subscriber event.EventSubscriber
}

func NewNodeCommanderService(
	publisher event.EventPublisher,
	subscriber event.EventSubscriber,
) *NodeCommanderService {
	s := &NodeCommanderService{
		publisher:  publisher,
		subscriber: subscriber,
	}

	return s
}

func (s NodeCommanderService) Connections(id string) ([]models.Connection, error) {
	ev := event.NewQueryNodeConnectionsEvent(id)
	err := s.publisher.Publish(ev)
	if err != nil {
		return nil, err
	}

	ch, close := s.subscriber.Subscribe(ev)

	if ch == nil {
		return nil, fmt.Errorf("coudln't subscribe to event, %v", ev)
	}
	defer close()

	r := <-ch
	err = r.Error()
	if err != nil {
		return nil, err
	}

	out, ok := r.(event.QueryNodeConnectionsEventOutput)
	if !ok {
		return nil, fmt.Errorf("failed to parse %v to connection slice", r)
	}

	return out.Connections, nil
}
