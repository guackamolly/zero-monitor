package event

import (
	"github.com/guackamolly/zero-monitor/internal/data/models"
)

type CloseSubscription func()

// An event is a message that is published under a channel, which
// one or more clients can subscribe to.
type Event interface {
	ID() string
}

// Represents the output of an event.
type EventOutput interface {
	Origin() Event
	Error() error
}

// Represents metadata of an event stream.
type EventStream interface {
	// If any, returns the public key used by publishers to authenticate
	// events.
	PublicKey() ([]byte, error)
	// Returns the address that the event stream is open on.
	Address() models.Address
}

// An interface for marking clients that can publish events to a channel.
type EventPublisher interface {
	EventStream
	Publish(Event) error
}

// An interface for marking clients that can subscribe the output of events from a channel.
type EventSubscriber interface {
	EventStream
	Subscribe(Event) (chan (EventOutput), CloseSubscription)
}
