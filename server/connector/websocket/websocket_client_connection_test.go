package websocketconnector

import (
	"errors"
	"francoisgergaud/3dGame/common/event"
	testwebsocket "francoisgergaud/3dGame/internal/testutils/common/connector"
	testrunner "francoisgergaud/3dGame/internal/testutils/common/runner"
	testserver "francoisgergaud/3dGame/internal/testutils/server"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClientWebSocketSender struct {
	mock.Mock
	testrunner.MockRunnable
}

func (mock *MockClientWebSocketSender) Stop() {
	mock.Called()
}

func TestNewWebSocketClientConnection(t *testing.T) {
	eventToSendToCLient := make(chan event.Event)
	clientEventSender := new(MockClientWebSocketSender)
	wsConnection := new(testwebsocket.MockWebsockeConnection)
	clientConnection := NewWebSocketClientConnection(eventToSendToCLient, clientEventSender, wsConnection)
	assert.Equal(t, eventToSendToCLient, clientConnection.eventToSendToCLient)
	assert.Same(t, clientEventSender, clientConnection.clientWebsocketSender)
	assert.Same(t, wsConnection, clientConnection.wsConnection)
}

func TestSendEventsToClient(t *testing.T) {
	eventToSendToCLient := make(chan event.Event, 2)
	clientConnection := WebSocketClientConnection{
		eventToSendToCLient: eventToSendToCLient,
	}
	event1 := event.Event{PlayerID: "playerID1"}
	event2 := event.Event{PlayerID: "playerID2"}
	clientConnection.SendEventsToClient([]event.Event{event1, event2})
	assert.Equal(t, event1, <-eventToSendToCLient)
	assert.Equal(t, event2, <-eventToSendToCLient)
}

func TestClose(t *testing.T) {
	clientWebsocketSender := new(MockClientWebSocketSender)
	websocketConnection := new(testwebsocket.MockWebsockeConnection)
	clientConnection := WebSocketClientConnection{
		clientWebsocketSender: clientWebsocketSender,
		wsConnection:          websocketConnection,
	}
	clientWebsocketSender.On("Stop")
	websocketConnection.On("Close").Return(nil)
	clientConnection.Close()
	mock.AssertExpectationsForObjects(t, clientWebsocketSender, websocketConnection)
}

func TestNewClientWebSocketListener(t *testing.T) {
	playerID := "playerID"
	wsConnection := new(testwebsocket.MockWebsockeConnection)
	server := new(testserver.MockServer)
	websocketClientListener := NewClientWebSocketListener(playerID, wsConnection, server)
	assert.Equal(t, playerID, websocketClientListener.playerID)
	assert.Equal(t, server, websocketClientListener.server)
	assert.Equal(t, wsConnection, websocketClientListener.wsConnection)
}

func TestClientWebSocketListenerRun(t *testing.T) {
	playerID := "playerID"
	wsConnection := new(testwebsocket.MockWebsockeConnection)
	server := new(testserver.MockServer)
	websocketClientListener := &ClientWebSocketListener{
		playerID:     playerID,
		wsConnection: wsConnection,
		server:       server,
	}
	eventFromClient := event.Event{
		Action:   "testAction",
		PlayerID: "testPlayerID",
	}
	wsConnection.On("ReadJSON", mock.MatchedBy(
		func(eventsRead *[]event.Event) bool {
			(*eventsRead) = append(*eventsRead, eventFromClient)
			return true
		},
	),
	).Return(nil).Once()
	wsConnection.On("ReadJSON", mock.AnythingOfType("*[]event.Event")).Return(errors.New("second read error"))
	eventFromClientWithReplacedPlayerID := eventFromClient
	eventFromClientWithReplacedPlayerID.PlayerID = playerID
	server.On("ReceiveEventFromClient", eventFromClientWithReplacedPlayerID)
	server.On("UnregisterClient", playerID)
	websocketClientListener.Run()
	mock.AssertExpectationsForObjects(t, wsConnection, server)
}

func TestNewClientWebSocketSender(t *testing.T) {
	wsConnection := new(testwebsocket.MockWebsockeConnection)
	eventToSendToClient := make(chan event.Event, 0)
	websocketClientSender := NewClientWebSocketSender(wsConnection, eventToSendToClient)
	assert.Equal(t, eventToSendToClient, websocketClientSender.eventToSendToClient)
	assert.Equal(t, wsConnection, websocketClientSender.wsConnection)
}

func TestClientWebSocketSenderStop(t *testing.T) {
	clientWebSocketSender := &ClientWebSocketSenderImpl{
		quit: make(chan interface{}),
	}
	go clientWebSocketSender.Run()
	clientWebSocketSender.Stop()
}

func TestClientWebSocketSenderRunWithError(t *testing.T) {
	wsConnection := new(testwebsocket.MockWebsockeConnection)
	eventsToSendToClient := make(chan event.Event, 0)
	clientWebSocketSender := &ClientWebSocketSenderImpl{
		wsConnection:        wsConnection,
		eventToSendToClient: eventsToSendToClient,
	}
	eventToSend := event.Event{
		PlayerID: "testPlayerID",
	}
	wsConnection.On("WriteJSON", []event.Event{eventToSend}).Return(errors.New("test-error"))
	go clientWebSocketSender.Run()
	eventsToSendToClient <- eventToSend
	mock.AssertExpectationsForObjects(t, wsConnection)
}
