package service

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/event"
)

// todo: this shouldn't be here
type SpeedtestUpdates struct {
	Speedtest models.Speedtest
	Updates   chan (models.Speedtest)
}

// Service for executing commands in nodes.
type NodeCommanderService struct {
	publisher  event.EventPublisher
	subscriber event.EventSubscriber
	speedtests map[string]SpeedtestUpdates
}

func NewNodeCommanderService(
	publisher event.EventPublisher,
	subscriber event.EventSubscriber,
) *NodeCommanderService {
	s := &NodeCommanderService{
		publisher:  publisher,
		subscriber: subscriber,
		speedtests: map[string]SpeedtestUpdates{},
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

func (s NodeCommanderService) KillProcess(id string, pid int32) error {
	ev := event.NewKillNodeProcessEvent(id, pid)
	out, err := event.PublishAndSubscribeFirst[event.KillNodeProcessEventOutput](ev, s.publisher, s.subscriber)
	if err != nil {
		return err
	}

	return out.Error()
}

func (s NodeCommanderService) StartSpeedtest(id string) (models.Speedtest, error) {
	ev := event.NewStartNodeSpeedtestEvent(id)
	out, err := event.PublishAndSubscribe[event.NodeSpeedtestEventOutput](ev, s.publisher, s.subscriber)
	if err != nil {
		return models.Speedtest{}, err
	}

	ch := make(chan (models.Speedtest))
	go func() {
		for o := range out {
			if err := o.Error(); err != nil {
				continue
			}

			ch <- o.Speedtest
		}

		close(ch)
	}()

	st, ok := <-ch
	if !ok {
		return st, fmt.Errorf("failed to start speedtest")
	}
	s.speedtests[st.ID] = SpeedtestUpdates{Speedtest: st, Updates: ch}

	return st, nil
}

// todo: this shouldn't be here
func (s NodeCommanderService) SpeedtestUpdates(id string) (chan (models.Speedtest), bool) {
	stup, ok := s.speedtests[id]
	if !ok {
		return nil, false
	}

	return stup.Updates, true
}

// todo: this shouldn't be here
func (s NodeCommanderService) Speedtest(id string) (models.Speedtest, bool) {
	stup, ok := s.speedtests[id]
	if !ok {
		return models.Speedtest{}, false
	}

	return stup.Speedtest, true
}
