package service_test

import (
	"testing"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/data/repositories"
	"github.com/guackamolly/zero-monitor/internal/event"
	"github.com/guackamolly/zero-monitor/internal/service"
)

type TestPublishSubscriber struct {
	event.EventPublisher
	event.EventSubscriber
	outputs []event.EventOutput
}

type TestSpeedtestStoreRepository struct {
	repositories.SpeedtestStoreRepository
	speedtests []models.Speedtest
}

func (r *TestSpeedtestStoreRepository) Save(nodeid string, speedtest models.Speedtest) error {
	r.speedtests = append(r.speedtests, speedtest)
	return nil
}

func (r TestSpeedtestStoreRepository) History(nodeid string) ([]models.Speedtest, error) {
	return r.speedtests, nil
}

func NewTestSpeedtestStoreRepository() *TestSpeedtestStoreRepository {
	return &TestSpeedtestStoreRepository{
		speedtests: []models.Speedtest{},
	}
}

func NewTestPublishSubscriber(
	outputs ...event.EventOutput,
) TestPublishSubscriber {
	return TestPublishSubscriber{outputs: outputs}
}

func (ps TestPublishSubscriber) Publish(event.Event) error {
	return nil
}

func (ps TestPublishSubscriber) Subscribe(ev event.Event) (chan (event.EventOutput), event.CloseSubscription) {
	outs := ps.outputs

	ch := make(chan (event.EventOutput))
	go func() {
		for _, out := range outs {
			ch <- out
		}

		close(ch)
	}()

	return ch, func() {}
}

func forceUnbufferedChannelClose[T any](ch chan (T)) {
	for range ch {
	}
	// todo: use wait group...
	time.Sleep(1 * time.Millisecond)
}

func TestSpeedtestsUpdatesAreNotAvailableAfterSpeedtestFinishes(t *testing.T) {
	ps := NewTestPublishSubscriber(
		event.NewNodeSpeedtestEventOutput(nil, models.NewSpeedtest("<speedtest-id>", "zero-monitor", "world", "ookla", 5), nil),
	)
	sps := NewTestSpeedtestStoreRepository()

	s := service.NewNodeSpeedtestService(ps, ps, sps)
	nid := "<node-id>"

	st, err := s.Start(nid)
	if err != nil {
		t.Fatalf("was not expecting Start() to fail, %v", err)
	}

	ch, ok := s.Updates(st.ID)
	if !ok {
		t.Fatal("was not expected Updates() to not return channel")
	}
	forceUnbufferedChannelClose(ch)

	_, ok = s.Updates(st.ID)
	if ok {
		t.Error("expected updates to not return channel")
	}
}

func TestSpeedtestsUpdatesChannelIsClosedAfterSpeedtestFinishes(t *testing.T) {
	ps := NewTestPublishSubscriber(
		event.NewNodeSpeedtestEventOutput(nil, models.NewSpeedtest("<speedtest-id>", "zero-monitor", "world", "ookla", 0), nil),
	)
	sps := NewTestSpeedtestStoreRepository()

	s := service.NewNodeSpeedtestService(ps, ps, sps)
	nid := "<node-id>"

	st, err := s.Start(nid)
	if err != nil {
		t.Fatalf("was not expecting Start() to fail, %v", err)
	}

	ch, ok := s.Updates(st.ID)
	if !ok {
		t.Fatal("was not expected Updates() to not return channel")
	}
	forceUnbufferedChannelClose(ch)
	_, ok = <-ch
	if ok {
		t.Error("expected updates channel to be closed")
	}
}

func TestSpeedtestIsUpdatedWheneverUpdatesChannelEmitsNewValues(t *testing.T) {
	ps := NewTestPublishSubscriber(
		event.NewNodeSpeedtestEventOutput(nil, models.NewSpeedtest("<speedtest-id>", "zero-monitor", "world", "ookla", 5), nil),
		event.NewNodeSpeedtestEventOutput(nil, models.NewSpeedtest("<speedtest-id>", "zero-monitor", "world", "ookla", 2), nil),
	)
	sps := NewTestSpeedtestStoreRepository()

	s := service.NewNodeSpeedtestService(ps, ps, sps)
	nid := "<node-id>"

	st, err := s.Start(nid)
	if err != nil {
		t.Fatalf("was not expecting Start() to fail, %v", err)
	}

	ch, ok := s.Updates(st.ID)
	if !ok {
		t.Fatal("was not expected Updates() to not return channel")
	}

	st = <-ch
	ust, _ := s.Speedtest(st.ID)
	if ust != st {
		t.Errorf("expected %v but got %v", st, ust)
	}
}
