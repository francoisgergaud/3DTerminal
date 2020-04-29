package testclient

import (
	"francoisgergaud/3dGame/client/connector"
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
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

//Initialize mocks the method of the same name
func (mock *MockEngine) Initialize(playerID string, playerState state.AnimatedElementState, worldMap world.WorldMap, otherPlayersState map[string]state.AnimatedElementState, serverTimeFrame uint32) {
	mock.Called(playerID, playerState, worldMap, otherPlayersState, serverTimeFrame)
}

//GetShutdown mocks the method of the same name
func (mock *MockEngine) GetShutdown() <-chan interface{} {
	args := mock.Called()
	return args.Get(0).(chan interface{})
}

//SetConnectionToServer mocks the method of the same name
func (mock *MockEngine) SetConnectionToServer(connectionToServer connector.ServerConnector) {
	mock.Called(connectionToServer)
}
