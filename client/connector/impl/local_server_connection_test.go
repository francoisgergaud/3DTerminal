package impl

import (
	"francoisgergaud/3dGame/common/event"
	testClient "francoisgergaud/3dGame/internal/testutils/client"
	testconnector "francoisgergaud/3dGame/internal/testutils/server/connector"
	"testing"
)

func TestSendEventsToServer(t *testing.T) {
	clientConnection := new(testconnector.MockClientConnection)
	serverConnection := LocalServerConnection{
		ClientConnection: clientConnection,
	}
	timeFramfe := uint32(0)
	events := make([]event.Event, 0)
	clientConnection.On("ReceiveEventsFromClient", timeFramfe, events)
	serverConnection.SendEventsToServer(timeFramfe, events)
}

func TestReceiveEventsFromServer(t *testing.T) {
	engine := new(testClient.MockEngine)
	serverConnection := LocalServerConnection{
		Engine: engine,
	}
	events := make([]event.Event, 0)
	engine.On("ReceiveEventsFromServer", events)
	serverConnection.ReceiveEventsFromServer(0, events)
}
