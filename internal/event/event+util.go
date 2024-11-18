package event

import (
	"fmt"
	"reflect"

	"github.com/guackamolly/zero-monitor/internal/logging"
)

// Publishes an event in a channel using [p], subscribes it's first value using [s], and closes
// before finishing.
func PublishAndSubscribeFirst[T any](ev Event, p EventPublisher, s EventSubscriber) (T, error) {
	var out T

	err := p.Publish(ev)
	if err != nil {
		return out, err
	}

	ch, close := s.Subscribe(ev)
	if ch == nil {
		return out, fmt.Errorf("couldn't subscribe to event, %v", ev)
	}
	defer close()

	r := <-ch
	err = r.Error()
	if err != nil {
		return out, err
	}

	fmt.Printf("r: %v\n", reflect.TypeOf(r))
	out, ok := r.(T)
	if !ok {
		return out, fmt.Errorf("failed to parse %v to connection slice", r)
	}

	return out, nil
}

func PublishAndSubscribe[T any](ev Event, p EventPublisher, s EventSubscriber) (chan (T), error) {
	var out chan (T)

	err := p.Publish(ev)
	if err != nil {
		return out, err
	}

	ch, sclose := s.Subscribe(ev)
	if ch == nil {
		return out, fmt.Errorf("couldn't subscribe to event, %v", ev)
	}

	out = make(chan (T))
	go func() {
		defer sclose()
		defer close(out)

		for r := range ch {
			tout, ok := r.(T)
			if !ok {
				logging.LogError("failed to parse %v to %v", r, new(T))
				continue
			}

			out <- tout
		}
	}()

	return out, nil
}
