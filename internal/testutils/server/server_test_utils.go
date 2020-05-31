package testserver

import (
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server/connector"

	"github.com/stretchr/testify/mock"
)

//MockServer is the mock for a server
type MockServer struct {
	mock.Mock
}

//RegisterPlayer mocks the method of the same name
func (mock *MockServer) RegisterPlayer(clientConnection connector.ClientConnection) string {
	args := mock.Called(clientConnection)
	return args.String(0)
}

//UnregisterClient mocks the method of the same name
func (mock *MockServer) UnregisterClient(playerID string) {
	mock.Called(playerID)
}

//ReceiveEventFromClient mocks the method of the same name
func (mock *MockServer) ReceiveEventFromClient(event event.Event) {
	mock.Called(event)
}

//Start mocks the method of the same name
func (mock *MockServer) Start() {
	mock.Called()
}
