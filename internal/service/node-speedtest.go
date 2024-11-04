package service

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/event"
)

// Service for managing nodes speedtests.
type NodeSpeedtestService struct {
	speedtests map[string]models.Speedtest
	updates    map[string]chan (models.Speedtest)
	publisher  event.EventPublisher
	subscriber event.EventSubscriber
}

func NewNodeSpeedtestService(
	publisher event.EventPublisher,
	subscriber event.EventSubscriber,
) *NodeSpeedtestService {
	s := &NodeSpeedtestService{
		speedtests: map[string]models.Speedtest{},
		updates:    map[string]chan models.Speedtest{},
		publisher:  publisher,
		subscriber: subscriber,
	}

	return s
}

func (s NodeSpeedtestService) Start(nodeid string) (models.Speedtest, error) {
	ev := event.NewStartNodeSpeedtestEvent(nodeid)
	out, err := event.PublishAndSubscribe[event.NodeSpeedtestEventOutput](ev, s.publisher, s.subscriber)
	if err != nil {
		return models.Speedtest{}, err
	}

	ch := make(chan (models.Speedtest))
	go func() {
		var sid string
		for o := range out {
			if err := o.Error(); err != nil {
				continue
			}

			sid = o.Speedtest.ID

			s.speedtests[sid] = o.Speedtest
			ch <- o.Speedtest
		}

		close(ch)
		delete(s.updates, sid)
	}()

	st, ok := <-ch
	if !ok {
		return st, fmt.Errorf("failed to start speedtest")
	}
	s.updates[st.ID] = ch

	return st, nil
}

func (s NodeSpeedtestService) Updates(id string) (chan (models.Speedtest), bool) {
	ch, ok := s.updates[id]
	if !ok {
		return nil, false
	}

	return ch, true
}

func (s NodeSpeedtestService) Speedtest(id string) (models.Speedtest, bool) {
	sts, ok := s.speedtests[id]
	if !ok {
		return models.Speedtest{}, false
	}

	return sts, true
}
