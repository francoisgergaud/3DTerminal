package testserver

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server/connector"

	"github.com/stretchr/testify/mock"
)

//MockServer is the mock for a server
type MockServer struct {
	mock.Mock
}

//RegisterPlayer mocks the method of the same name
func (mock *MockServer) RegisterPlayer(clientConnection connector.ClientConnection) (playerID string, worldMap world.WorldMap, animatedElementState state.AnimatedElementState, otherPlayers map[string]state.AnimatedElementState) {
	args := mock.Called(clientConnection)
	playerID = args.String(0)
	worldMap = args.Get(1).(world.WorldMap)
	animatedElementState = args.Get(2).(state.AnimatedElementState)
	otherPlayers = args.Get(3).(map[string]state.AnimatedElementState)
	return
}

//UnregisterClient mocks the method of the same name
func (mock *MockServer) UnregisterClient(playerID string) {
	mock.Called(playerID)
}

//ReceiveEventsFromClient mocks the method of the same name
func (mock *MockServer) ReceiveEventsFromClient(events []event.Event) {
	mock.Called(events)
}
