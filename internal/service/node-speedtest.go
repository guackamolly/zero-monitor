package service

import (
	"fmt"
	"slices"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/event"
	"github.com/guackamolly/zero-monitor/internal/logging"
)

const (
	SpeedtestHistoryLimit = 25
)

// Service for managing nodes speedtests.
type NodeSpeedtestService struct {
	// NodeID -> Speedtests
	history map[string][]models.Speedtest
	// NodeID -> bool
	loadedHistory map[string]bool
	// Speedtest ID -> Speedtest
	cache map[string]models.Speedtest
	// Speedtest ID -> chan(Speedtest)
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
		cache:         map[string]models.Speedtest{},
		updates:       map[string]chan models.Speedtest{},
		history:       map[string][]models.Speedtest{},
		loadedHistory: map[string]bool{},
		publisher:     publisher,
		subscriber:    subscriber,
		store:         store,
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
			s.cache[st.ID] = st
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

		if hs := s.history[nodeid]; hs != nil {
			hs = slices.Insert(hs, 0, st)
			if len(hs) > SpeedtestHistoryLimit {
				hs = hs[0:SpeedtestHistoryLimit]
			}
			s.history[nodeid] = hs
		}
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
	sts, ok := s.cache[id]
	if ok {
		return sts, true
	}

	sts, ok, err := s.store.Lookup(id)
	if err != nil {
		logging.LogError("failed to lookup speedtest on store, %v", err)
	}

	return sts, ok
}

func (s NodeSpeedtestService) History(nodeid string) ([]models.Speedtest, bool) {
	s.loadHistory(nodeid)

	sts, ok := s.history[nodeid]
	return sts, ok
}

func (s *NodeSpeedtestService) loadHistory(nodeid string) {
	if s.loadedHistory[nodeid] {
		return
	}

	hs, err := s.store.History(nodeid)
	if err != nil {
		logging.LogError("failed to load history for node %s, %v", nodeid, err)
		return
	}

	slices.SortFunc(hs, func(x, y models.Speedtest) int {
		return -x.TakenAt.Compare(y.TakenAt)
	})

	s.loadedHistory[nodeid] = true
	lhs := []models.Speedtest{}
	for i, st := range hs {
		if i < SpeedtestHistoryLimit {
			lhs = append(lhs, st)
		}

		s.cache[st.ID] = st
	}
	s.history[nodeid] = lhs
}
