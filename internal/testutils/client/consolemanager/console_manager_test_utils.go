package testconsolemanager

import (
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/common/event"

	"github.com/stretchr/testify/mock"
)

//MockConsoleEventManager mocks the calls to the ConsoleEventManager interface.
type MockConsoleEventManager struct {
	mock.Mock
}

//Run mocks the call to the the method of the same name.
func (mock *MockConsoleEventManager) Run() error {
	args := mock.Called()
	return args.Error(0)
}

//SetPlayer mocks the call to the the method of the same name.
func (mock *MockConsoleEventManager) SetPlayer(engine client.Engine) {
	mock.Called(engine)
}

//MockEventPublisher mocks an event-publisher
type MockEventPublisher struct {
	mock.Mock
}

//PublishEvent mocks the method of the same name
func (mock *MockEventPublisher) PublishEvent(event event.Event) {
	mock.Called(event)
}
