package publisher

import "francoisgergaud/3dGame/common/event"

//EventListener defines an event-listener
type EventListener interface {
	ReceiveEvent(event.Event)
}

//EventPublisher manage the sending of events
type EventPublisher interface {
	PublishEvent(event.Event)
	RegisterListener(EventListener)
}
