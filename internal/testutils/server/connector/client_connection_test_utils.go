package testconnector

import (
	"francoisgergaud/3dGame/common/event"

	"github.com/stretchr/testify/mock"
)

//MockClientConnection mocks a client-connection
type MockClientConnection struct {
	mock.Mock
}

//ReceiveEventsFromClient mocks the method of the same name
func (mock *MockClientConnection) ReceiveEventsFromClient(events []event.Event) {
	mock.Called(events)
}

//SendEventsToClient mocks the method of the same name
func (mock *MockClientConnection) SendEventsToClient(events []event.Event) {
	mock.Called(events)
}

//Close mocks the method of the same name
func (mock *MockClientConnection) Close() {
	mock.Called()
}
