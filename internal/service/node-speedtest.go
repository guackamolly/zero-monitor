package service

import (
	"fmt"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/event"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Service for managing nodes speedtests.
type NodeSpeedtestService struct {
	speedtests map[string]models.Speedtest
	updates    map[string]chan (models.Speedtest)
	publisher  event.EventPublisher
	subscriber event.EventSubscriber
	store      repositories.SpeedtestStoreRepository
}

func NewNodeSpeedtestService(
	publisher event.EventPublisher,
	subscriber event.EventSubscriber,
	store repositories.SpeedtestStoreRepository,
) *NodeSpeedtestService {
	s := &NodeSpeedtestService{
		speedtests: map[string]models.Speedtest{},
		updates:    map[string]chan models.Speedtest{},
		publisher:  publisher,
		subscriber: subscriber,
		store:      store,
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
		var st models.Speedtest
		for o := range out {
			if err := o.Error(); err != nil {
				continue
			}

			st = o.Speedtest
			s.speedtests[st.ID] = st
			ch <- st

			// todo: this is a workaround to break out of channel, since upstream is not closing channel
			if st.Finished() {
				break
			}
		}

		err = s.store.Save(nodeid, st)
		if err != nil {
			logging.LogError("couldn't store speedtest, %v", err)
		}

		close(ch)
		delete(s.updates, st.ID)
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
