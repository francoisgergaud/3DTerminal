package impl

import (
	"francoisgergaud/3dGame/common/event"
	testClient "francoisgergaud/3dGame/internal/testutils/client"
	"testing"
)

func TestSendEventsToClient(t *testing.T) {
	engine := new(testClient.MockEngine)
	serverConnection := LocalServerConnectionImpl{
		engine: engine,
	}
	events := make([]event.Event, 0)
	engine.On("ReceiveEventsFromServer", events)
	serverConnection.SendEventsToClient(events)
}
