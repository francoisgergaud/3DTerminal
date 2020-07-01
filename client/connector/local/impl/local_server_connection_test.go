package impl

import (
	"francoisgergaud/3dGame/common/event"
	testClient "francoisgergaud/3dGame/internal/testutils/client"
	testServer "francoisgergaud/3dGame/internal/testutils/server"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewLocalServerConnection(t *testing.T) {
	engine := new(testClient.MockEngine)
	server := new(testServer.MockServer)
	quit := make(chan interface{})
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
	server.On("ReceiveEventFromClient", mock.MatchedBy(
		func(eventParameter event.Event) bool {
			//ensure cloning happened by checking the pointers
			if &eventParameter != &eventToSend {
				return true
			}
			return false
		},
	),
	)
	serverConnection.NotifyServer(events)
	mock.AssertExpectationsForObjects(t, server)
}

func TestNotifyServerWithCloneError(t *testing.T) {
	serverConnection := LocalServerConnectionImpl{}
	eventToSend := event.Event{
		Action: "fakeAction",
		ExtraData: map[string]interface{}{
			"wrongKey": "wrongValue",
		},
	}
	events := []event.Event{eventToSend}

	assert.Error(t, serverConnection.NotifyServer(events))
}

func TestSendEventsToClient(t *testing.T) {
	engine := new(testClient.MockEngine)
	serverConnection := LocalServerConnectionImpl{
		engine: engine,
	}
	eventToSend := event.Event{Action: "fakeAction"}
	events := []event.Event{eventToSend}
	engine.On("ReceiveEventsFromServer", mock.MatchedBy(
		func(eventsParameter []event.Event) bool {
			//ensure cloning happened by checking the pointers
			if &eventsParameter[0] != &eventToSend {
				return true
			}
			return false
		},
	),
	)
	serverConnection.SendEventsToClient(events)
	mock.AssertExpectationsForObjects(t, engine)
}

func TestSendEventsToClientWithCloneError(t *testing.T) {
	serverConnection := LocalServerConnectionImpl{}
	eventToSend := event.Event{
		Action: "fakeAction",
		ExtraData: map[string]interface{}{
			"wrongKey": "wrongValue",
		},
	}
	events := []event.Event{eventToSend}

	assert.Error(t, serverConnection.SendEventsToClient(events))
}
