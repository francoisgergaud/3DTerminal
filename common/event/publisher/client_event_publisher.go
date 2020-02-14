package publisher

import "francoisgergaud/3dGame/common/event"

//EventPublisher manage the sending of events
type EventPublisher interface {
	PublishEvent(event.Event)
	RegisterListener(chan<- event.Event)
}
