package events

import (
	"context"
	"errors"
	"slices"
	"sync"
)

var (
	ErrHandlerAlreadyRegistered = errors.New("handler already registered")
)

type eventDispatcher struct {
	handlers map[EventType][]EventHandlerInterface
}

func NewEventDispatcher() EventDispatcherInterface {
	return &eventDispatcher{
		handlers: make(map[EventType][]EventHandlerInterface),
	}
}

func (ed *eventDispatcher) Dispatch(ctx context.Context, event Event) error {
	if handlers, ok := ed.handlers[event.GetType()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (ed *eventDispatcher) RegisterHandler(eventType EventType, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventType]; ok {
		if slices.Contains(ed.handlers[eventType], handler) {
			return ErrHandlerAlreadyRegistered
		}
	}
	ed.handlers[eventType] = append(ed.handlers[eventType], handler)
	return nil
}
