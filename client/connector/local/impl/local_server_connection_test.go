package impl

import (
	"francoisgergaud/3dGame/common/event"
	testClient "francoisgergaud/3dGame/internal/testutils/client"
	testServer "francoisgergaud/3dGame/internal/testutils/server"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestNewLocalServerConnection(t *testing.T) {
	engine := new(testClient.MockEngine)
	server := new(testServer.MockServer)
	quit := make(chan struct{})
	engine.On("ConnectToServer", mock.AnythingOfType("*impl.LocalServerConnectionImpl"))
	server.On("RegisterPlayer", mock.AnythingOfType("*impl.LocalServerConnectionImpl")).Return("playerID")
	NewLocalServerConnection(engine, server, quit)
}

func TestNotifyServer(t *testing.T) {
	server := new(testServer.MockServer)
	serverConnection := LocalServerConnectionImpl{
		server: server,
	}
	eventToSend := event.Event{Action: "fakeAction"}
	events := []event.Event{eventToSend}
	server.On("ReceiveEventFromClient", eventToSend)
	serverConnection.NotifyServer(events)
}

func TestDisconnect(t *testing.T) {}

func TestSendEventsToClient(t *testing.T) {
	engine := new(testClient.MockEngine)
	serverConnection := LocalServerConnectionImpl{
		engine: engine,
	}
	events := make([]event.Event, 0)
	engine.On("ReceiveEventsFromServer", events)
	serverConnection.SendEventsToClient(events)
}
