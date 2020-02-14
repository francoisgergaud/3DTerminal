package testconnector

import (
	"francoisgergaud/3dGame/common/event"

	"github.com/stretchr/testify/mock"
)

//MockServerConnection mocks a server-connection from client
type MockServerConnection struct {
	mock.Mock
}

//SendEventsToServer mocks the method of the same name
func (mock *MockServerConnection) SendEventsToServer(timeFrame uint32, events []event.Event) {
	mock.Called(timeFrame, events)
}

//ReceiveEventsFromServer mocks the method of the same name
func (mock *MockServerConnection) ReceiveEventsFromServer(timeFrame uint32, events []event.Event) {
	mock.Called(timeFrame, events)
}
