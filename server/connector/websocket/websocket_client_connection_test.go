package websocketconnector

import (
	"errors"
	"francoisgergaud/3dGame/common/event"
	testwebsocket "francoisgergaud/3dGame/internal/testutils/common/connector"
	testserver "francoisgergaud/3dGame/internal/testutils/server"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewWebSocketClientConnection(t *testing.T) {
	eventToSendToCLient := make(chan event.Event)
	clientConnection := NewWebSocketClientConnection(eventToSendToCLient)
	assert.Equal(t, eventToSendToCLient, clientConnection.eventToSendToCLient)
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

func TestClientWebSocketSenderRun(t *testing.T) {
	wsConnection := new(testwebsocket.MockWebsockeConnection)
	eventsToSendToClient := make(chan event.Event, 0)
	clientWebSocketSender := &ClientWebSocketSender{
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
