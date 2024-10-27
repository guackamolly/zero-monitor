package service

import (
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
	out, err := event.PublishAndSubscribeFirst[event.QueryNodeConnectionsEventOutput](ev, s.publisher, s.subscriber)
	if err != nil {
		return nil, err
	}

	return out.Connections, nil
}

func (s NodeCommanderService) Processes(id string) ([]models.Process, error) {
	ev := event.NewQueryNodeProcessesEvent(id)
	out, err := event.PublishAndSubscribeFirst[event.QueryNodeProcessesEventOutput](ev, s.publisher, s.subscriber)
	if err != nil {
		return nil, err
	}

	return out.Processes, nil
}
