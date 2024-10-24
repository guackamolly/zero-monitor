package event

// Projects an abstract class that implements the [Event] Interface.
type BaseEvent struct {
	Event
	id string
}

type BaseEventOutput struct {
	EventOutput
	origin Event
	data   any
}

func (e BaseEvent) ID() string {
	return e.id
}

func (o BaseEventOutput) Origin() Event {
	return o.origin
}

func (o BaseEventOutput) Data() any {
	if o.Error() == nil {
		return o.data
	}

	return nil
}

func (o BaseEventOutput) Error() error {
	if e, ok := o.data.(error); ok {
		return e
	}

	return nil
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
	data any,
) EventOutput {
	return BaseEventOutput{
		origin: origin,
		data:   data,
	}
}
