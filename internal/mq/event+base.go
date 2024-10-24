package mq

// Projects an abstract class that implements the [Event] Interface.
type BaseEvent struct {
	Event
	id    string
	topic Topic
}

type BaseEventOutput struct {
	EventOutput
	origin Event
	data   any
}

// Defines an [Event] that bundles required data for subscribers.
type DataEvent struct {
	Event
	Data any
}

func (e BaseEvent) ID() string {
	return e.id
}

func (e BaseEvent) Topic() Topic {
	return e.topic
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
	topic Topic,
) Event {
	return BaseEvent{
		id:    id,
		topic: topic,
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

func NewDataEvent(
	id string,
	topic Topic,
	data any,
) DataEvent {
	return DataEvent{
		Event: NewBaseEvent(id, topic),
		Data:  data,
	}
}
