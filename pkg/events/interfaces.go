package events

import (
	"context"
	"sync"
	"time"
)

type (
	EventType string

	Event interface {
		GetType() EventType
		SetPayload(any)
		GetPayload() any
		GetDateTime() time.Time
	}

	EventDispatcherInterface interface {
		Dispatch(ctx context.Context, event Event) error
		RegisterHandler(eventType EventType, handler EventHandlerInterface) error
	}

	EventHandlerInterface interface {
		Handle(event Event, wg *sync.WaitGroup)
	}

	EventCreatorInterface interface {
		Create(payload any) Event
		EventType() EventType
	}
)
