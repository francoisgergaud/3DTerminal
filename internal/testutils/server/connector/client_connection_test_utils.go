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
func (mock *MockClientConnection) ReceiveEventsFromClient(timeFrame uint32, events []event.Event) {
	mock.Called(timeFrame, events)
}

//SendEventsToClient mocks the method of the same name
func (mock *MockClientConnection) SendEventsToClient(timeFrame uint32, events []event.Event) {
	mock.Called(timeFrame, events)
}
