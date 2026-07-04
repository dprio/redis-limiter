package events

import "time"

type (
	EventCreator struct {
		eventType EventType
	}

	event struct {
		tp       EventType
		payload  any
		dateTime time.Time
	}
)

func NewEventCreator(eventType EventType) EventCreatorInterface {
	return &EventCreator{
		eventType: eventType,
	}

}

func (ec *EventCreator) Create(payload any) Event {
	return newEvent(ec.eventType, payload)
}

func (ec *EventCreator) EventType() EventType {
	return ec.eventType
}

func newEvent(tp EventType, payload any) Event {
	return &event{
		tp:       tp,
		payload:  payload,
		dateTime: time.Now(),
	}
}

func (ev *event) SetPayload(payload any) {
	ev.payload = payload
}

func (ev *event) GetPayload() any {
	return ev.payload
}

func (ev *event) GetType() EventType {
	return ev.tp
}

func (ev *event) GetDateTime() time.Time {
	return ev.dateTime
}
