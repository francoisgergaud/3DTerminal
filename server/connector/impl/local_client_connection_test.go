package impl

import (
	"francoisgergaud/3dGame/common/event"
	testconnector "francoisgergaud/3dGame/internal/testutils/client/connector"
	testserver "francoisgergaud/3dGame/internal/testutils/server"
	"testing"
)

func TestReceiveEventsFromClient(t *testing.T) {
	server := new(testserver.MockServer)
	localConnection := LocalClientConnection{
		Server: server,
	}
	timeframe := uint32(1)
	events := make([]event.Event, 0)
	server.On("ReceiveEventsFromClient", events)
	localConnection.ReceiveEventsFromClient(timeframe, events)
}

func TestSendEventsToClient(t *testing.T) {
	serverConnection := new(testconnector.MockServerConnection)
	localConnection := LocalClientConnection{
		ServerConnection: serverConnection,
	}
	timeframe := uint32(1)
	events := make([]event.Event, 0)
	serverConnection.On("ReceiveEventsFromServer", timeframe, events)
	localConnection.SendEventsToClient(timeframe, events)
}
