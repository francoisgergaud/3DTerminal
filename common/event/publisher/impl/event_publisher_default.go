package impl

import (
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/event/publisher"
)

//NewEventPublisherImpl is factory creating the default implementation of a EventPublisher
func NewEventPublisherImpl() *EventPublisherImpl {
	return &EventPublisherImpl{
		listeners: make([]publisher.EventListener, 0),
	}
}

//EventPublisherImpl is the implementation of a EventPublisher
type EventPublisherImpl struct {
	listeners []publisher.EventListener
}

//PublishEvent fanout an event to the listeners.
func (eventPublisher *EventPublisherImpl) PublishEvent(event event.Event) {
	for _, listener := range eventPublisher.listeners {
		listener.ReceiveEvent(event)
	}
}

//RegisterListener registers a new listener
func (eventPublisher *EventPublisherImpl) RegisterListener(eventListener publisher.EventListener) {
	eventPublisher.listeners = append(eventPublisher.listeners, eventListener)
}
