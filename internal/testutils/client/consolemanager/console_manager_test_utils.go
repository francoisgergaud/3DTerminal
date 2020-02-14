package testconsolemanager

import (
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/event"

	"github.com/stretchr/testify/mock"
)

//MockConsoleEventManager mocks the calls to the ConsoleEventManager interface.
type MockConsoleEventManager struct {
	mock.Mock
}

//Listen mocks the call to the the method of the same name.
func (mock *MockConsoleEventManager) Listen() {
	mock.Called()
}

//SetPlayer mocks the call to the the method of the same name.
func (mock *MockConsoleEventManager) SetPlayer(player player.Player) {
	mock.Called(player)
}

//MockEventPublisher mocks an event-publisher
type MockEventPublisher struct {
	mock.Mock
}

//PublishEvent mocks the method of the same name
func (mock *MockEventPublisher) PublishEvent(event event.Event) {
	mock.Called(event)
}
