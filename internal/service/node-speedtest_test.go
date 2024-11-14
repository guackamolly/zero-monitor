package service_test

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/guackamolly/zero-monitor/internal/data/models"
	"github.com/guackamolly/zero-monitor/internal/event"
	"github.com/guackamolly/zero-monitor/internal/service"
)

type TestPublishSubscriber struct {
	event.EventPublisher
	event.EventSubscriber
	outputs []event.EventOutput
}

type TestSpeedtestStoreRepository struct {
	speedtests []models.Speedtest
}

func (r *TestSpeedtestStoreRepository) Save(nodeid string, speedtest models.Speedtest) error {
	r.speedtests = append(r.speedtests, speedtest)
	return nil
}

func (r TestSpeedtestStoreRepository) History(nodeid string) ([]models.Speedtest, error) {
	speedtestsCopy := make([]models.Speedtest, len(r.speedtests))
	copy(speedtestsCopy, r.speedtests)

	return speedtestsCopy, nil
}

func (r TestSpeedtestStoreRepository) Lookup(id string) (models.Speedtest, bool, error) {
	for _, st := range r.speedtests {
		if st.ID == id {
			return st, true, nil
		}
	}

	return models.Speedtest{}, false, nil
}

func NewTestSpeedtestStoreRepository(
	speedtests ...models.Speedtest,
) *TestSpeedtestStoreRepository {
	return &TestSpeedtestStoreRepository{
		speedtests: speedtests,
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

func TestSpeedtestLookupOnStoreIfSpeedtestIsNotCached(t *testing.T) {
	speedtest := models.NewSpeedtest("<speedtest-id>", "zero-monitor", "world", "ookla", 5)
	ps := NewTestPublishSubscriber()
	sps := NewTestSpeedtestStoreRepository(speedtest)
	s := service.NewNodeSpeedtestService(ps, ps, sps)

	st, _ := s.Speedtest(speedtest.ID)
	if st != speedtest {
		t.Errorf("expected %v but got %v", speedtest, st)
	}
}

func TestHistoryReturnsStoredSpeedtestsBeforeAnySpeedtestFinishes(t *testing.T) {
	speedtest := models.NewSpeedtest("<speedtest-id>", "zero-monitor", "world", "ookla", 5)
	storeSpeedtests := []models.Speedtest{speedtest}
	ps := NewTestPublishSubscriber()
	sps := NewTestSpeedtestStoreRepository(storeSpeedtests...)
	s := service.NewNodeSpeedtestService(ps, ps, sps)

	nodeid := "node.id"

	hs, _ := s.History(nodeid)
	if !slices.Equal(hs, storeSpeedtests) {
		t.Errorf("expected %v but got %v", storeSpeedtests, hs)
	}
}

func TestHistoryReturnsStoredSpeedtestsAndStartedSpeedtestAfterItFinishes(t *testing.T) {
	speedtest := models.NewSpeedtest("<speedtest-id>", "zero-monitor", "world", "ookla", 5)
	speedtest2 := models.Speedtest{ID: "<speedtest-id-2>", Phase: models.SpeedtestFinish}
	storeSpeedtests := []models.Speedtest{speedtest}
	ps := NewTestPublishSubscriber(
		event.NewNodeSpeedtestEventOutput(nil, speedtest2, nil),
	)
	sps := NewTestSpeedtestStoreRepository(storeSpeedtests...)
	s := service.NewNodeSpeedtestService(ps, ps, sps)
	nodeid := "node.id"
	// 1. load history
	s.History(nodeid)

	// 2. start speedtest
	st, err := s.Start(nodeid)
	if err != nil {
		t.Fatalf("didn't expect Start() to fail, %v", err)
	}
	ch, ok := s.Updates(st.ID)
	if !ok {
		t.Fatal("didn't expect Updates() to fail")
	}

	forceUnbufferedChannelClose(ch)

	// 3. load history again
	hs, _ := s.History(nodeid)
	expected := []models.Speedtest{speedtest2, speedtest}
	if !slices.Equal(hs, expected) {
		t.Errorf("expected %v but got %v", expected, hs)
	}
}

func TestHistoryRespectResultsLimit(t *testing.T) {
	limit := service.SpeedtestHistoryLimit
	storeSpeedtests := []models.Speedtest{}
	for i := 0; i < limit+1; i++ {
		st := models.NewSpeedtest(fmt.Sprintf("<speedtest-id-%d>", i), "zero-monitor", "world", "ookla", 5)
		storeSpeedtests = append(storeSpeedtests, st)
	}
	ps := NewTestPublishSubscriber()
	sps := NewTestSpeedtestStoreRepository(storeSpeedtests...)
	s := service.NewNodeSpeedtestService(ps, ps, sps)

	nodeid := "node.id"
	hs, _ := s.History(nodeid)
	if len(hs) != limit {
		t.Errorf("expected history to return %d results, but got %d", limit, len(hs))
	}
}

func TestHistoryReturnsResultsSortedByTakenDateInAscendingOrder(t *testing.T) {
	st1 := models.NewSpeedtest("<speedtest-id-1>", "zero-monitor", "world", "ookla", 5)
	st2 := models.NewSpeedtest("<speedtest-id-2>", "zero-monitor", "world", "ookla", 5)
	st3 := models.NewSpeedtest("<speedtest-id-3>", "zero-monitor", "world", "ookla", 5)
	st2.TakenAt = st3.TakenAt.Add(1)
	storeSpeedtests := []models.Speedtest{st1, st2, st3}

	ps := NewTestPublishSubscriber()
	sps := NewTestSpeedtestStoreRepository(storeSpeedtests...)
	s := service.NewNodeSpeedtestService(ps, ps, sps)

	nodeid := "node.id"
	hs, _ := s.History(nodeid)
	expected := []models.Speedtest{st2, st3, st1}
	if !slices.Equal(hs, expected) {
		t.Errorf("expected %v but got %v", expected, hs)
	}
}

func TestAfterSpeedtestFinishesAndTheHistoryMapReachesLimitLastHistoryElementIsRemoved(t *testing.T) {
	limit := service.SpeedtestHistoryLimit
	storeSpeedtests := []models.Speedtest{}
	for i := 0; i < limit; i++ {
		st := models.NewSpeedtest(fmt.Sprintf("<speedtest-id-%d>", i), "zero-monitor", "world", "ookla", 5)
		storeSpeedtests = append(storeSpeedtests, st)
	}
	speedtest := models.Speedtest{ID: "<speedtest-id-top>", Phase: models.SpeedtestFinish}

	ps := NewTestPublishSubscriber(
		event.NewNodeSpeedtestEventOutput(nil, speedtest, nil),
	)
	sps := NewTestSpeedtestStoreRepository(storeSpeedtests...)
	s := service.NewNodeSpeedtestService(ps, ps, sps)
	nodeid := "node.id"
	// 1. load history
	s.History(nodeid)

	// 2. start speedtest
	st, err := s.Start(nodeid)
	if err != nil {
		t.Fatalf("didn't expect Start() to fail, %v", err)
	}
	ch, ok := s.Updates(st.ID)
	if !ok {
		t.Fatal("didn't expect Updates() to fail")
	}

	forceUnbufferedChannelClose(ch)

	// 3. load history again
	hs, _ := s.History(nodeid)
	expected := []models.Speedtest{speedtest}
	for i := limit - 1; i >= 1; i-- {
		expected = append(expected, storeSpeedtests[i])
	}

	if !slices.Equal(hs, expected) {
		t.Errorf("expected %v but got %v", expected, hs)
	}
}
