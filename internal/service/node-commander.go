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

func (s NodeCommanderService) Packages(id string) ([]models.Package, error) {
	ev := event.NewQueryNodePackagesEvent(id)
	out, err := event.PublishAndSubscribeFirst[event.QueryNodePackagesEventOutput](ev, s.publisher, s.subscriber)
	if err != nil {
		return nil, err
	}

	return out.Packages, nil
}

func (s NodeCommanderService) Processes(id string) ([]models.Process, error) {
	ev := event.NewQueryNodeProcessesEvent(id)
	out, err := event.PublishAndSubscribeFirst[event.QueryNodeProcessesEventOutput](ev, s.publisher, s.subscriber)
	if err != nil {
		return nil, err
	}

	return out.Processes, nil
}

func (s NodeCommanderService) KillProcess(id string, pid int32) error {
	ev := event.NewKillNodeProcessEvent(id, pid)
	out, err := event.PublishAndSubscribeFirst[event.KillNodeProcessEventOutput](ev, s.publisher, s.subscriber)
	if err != nil {
		return err
	}

	return out.Error()
}

func (s NodeCommanderService) Disconnect(id string) error {
	ev := event.NewDisconnectNodeEvent(id)
	out, err := event.PublishAndSubscribeFirst[event.DisconnectNodeEventOutput](ev, s.publisher, s.subscriber)
	if err != nil {
		return err
	}

	return out.Error()
}
