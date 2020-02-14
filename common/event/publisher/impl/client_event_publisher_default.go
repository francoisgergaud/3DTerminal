package impl

import (
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/event/publisher"
)

//NewEventPublisherImpl is factory creating the default implementation of a EventPublisher
func NewEventPublisherImpl() publisher.EventPublisher {
	return &EventPublisherImpl{
		listeners: make([]chan<- event.Event, 0),
	}
}

//EventPublisherImpl is the implementation of a EventPublisher
type EventPublisherImpl struct {
	listeners []chan<- event.Event
}

//PublishEvent fanout an event to the listeners.
func (eventPublisher *EventPublisherImpl) PublishEvent(event event.Event) {
	for _, listener := range eventPublisher.listeners {
		listener <- event
	}
}

//RegisterListener registers a new listener
func (eventPublisher *EventPublisherImpl) RegisterListener(listener chan<- event.Event) {
	eventPublisher.listeners = append(eventPublisher.listeners, listener)
}
