package event

// Projects an abstract class that implements the [Event] Interface.
type BaseEvent struct {
	Event
	id string
}

type BaseEventOutput struct {
	EventOutput
	origin Event
	err    error
}

func (e BaseEvent) ID() string {
	return e.id
}

func (o BaseEventOutput) Origin() Event {
	return o.origin
}

func (o BaseEventOutput) Error() error {
	return o.err
}

func NewBaseEvent(
	id string,
) Event {
	return BaseEvent{
		id: id,
	}
}

func NewBaseEventOutput(
	origin Event,
	err error,
) EventOutput {
	return BaseEventOutput{
		origin: origin,
		err:    err,
	}
}
