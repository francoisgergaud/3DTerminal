package testclient

import (
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/event"

	"github.com/stretchr/testify/mock"
)

//MockEngine is a mocked engine
type MockEngine struct {
	mock.Mock
}

//StartEngine mocks the method of the same name
func (mock *MockEngine) StartEngine() {
	mock.Called()
}

//ReceiveEventsFromServer mocks the method of the same name
func (mock *MockEngine) ReceiveEventsFromServer(events []event.Event) {
	mock.Called(events)
}

//GetPlayer mocks the method of the same name
func (mock *MockEngine) GetPlayer() player.Player {
	args := mock.Called()
	return args.Get(0).(player.Player)
}
