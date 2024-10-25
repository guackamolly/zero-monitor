package event

import "fmt"

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

	out, ok := r.(T)
	if !ok {
		return out, fmt.Errorf("failed to parse %v to connection slice", r)
	}

	return out, nil
}
