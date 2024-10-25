package event

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

// An interface for marking clients that can publish events to a channel.
type EventPublisher interface {
	Publish(Event) error
}

// An interface for marking clients that can subscribe the output of events from a channel.
type EventSubscriber interface {
	Subscribe(Event) (chan (EventOutput), CloseSubscription)
}