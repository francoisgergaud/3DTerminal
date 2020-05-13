package testconnector

import (
	"francoisgergaud/3dGame/common/event"

	"github.com/stretchr/testify/mock"
)

//MockServerConnection mocks a server-connection from client
type MockServerConnection struct {
	mock.Mock
}

//NotifyServer mocks the method of the same name
func (mock *MockServerConnection) NotifyServer(events []event.Event) error {
	args := mock.Called(events)
	return args.Error(0)
}

//Disconnect mocks the method of the same name
func (mock *MockServerConnection) Disconnect() {
	mock.Called()
}
