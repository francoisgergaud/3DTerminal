package testclient

import (
	"francoisgergaud/3dGame/client/connector"
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/event"

	"github.com/stretchr/testify/mock"
)

//MockEngine is a mocked engine
type MockEngine struct {
	mock.Mock
}

//ReceiveEventsFromServer mocks the method of the same name
func (mock *MockEngine) ReceiveEventsFromServer(events []event.Event) {
	mock.Called(events)
}

//Player mocks the method of the same name
func (mock *MockEngine) Player() player.Player {
	args := mock.Called()
	return args.Get(0).(player.Player)
}

//Shutdown mocks the method of the same name
func (mock *MockEngine) Shutdown() <-chan interface{} {
	args := mock.Called()
	return args.Get(0).(chan interface{})
}

//ConnectToServer mocks the method of the same name
func (mock *MockEngine) ConnectToServer(connectionToServer connector.ServerConnector) {
	mock.Called(connectionToServer)
}

//OtherPlayers mocks the method of the name
func (mock *MockEngine) OtherPlayers() map[string]animatedelement.AnimatedElement {
	args := mock.Called()
	return args.Get(0).(map[string]animatedelement.AnimatedElement)
}
