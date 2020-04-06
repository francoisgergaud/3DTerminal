package testeventpublisher

import (
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/event/publisher"

	"github.com/stretchr/testify/mock"
)

type MockEventListener struct {
	mock.Mock
}

//ReceiveEvent mocks the method of the same name
func (mock *MockEventListener) ReceiveEvent(event event.Event) {
	mock.Called(event)
}

//MockEventPublisher mocks an event-publisher
type MockEventPublisher struct {
	mock.Mock
}

//PublishEvent mocks the method of the same name
func (mock *MockEventPublisher) PublishEvent(event event.Event) {
	mock.Called(event)
}

//RegisterListener mocks the method of the same name
func (mock *MockEventPublisher) RegisterListener(listener publisher.EventListener) {
	mock.Called(listener)
}
