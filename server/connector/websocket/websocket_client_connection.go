package websocketconnector

import (
	"fmt"
	websocket "francoisgergaud/3dGame/common/connector"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/runner"
	"francoisgergaud/3dGame/server"
)

func bufferProvider() []event.Event {
	return make([]event.Event, 0)
}

//NewWebSocketClientConnection is a factory for WebSocketClientConnection
func NewWebSocketClientConnection(eventToSendToCLient chan event.Event, clientWebsocketSender ClientWebSocketSender, wsConnection websocket.WebsocketConnection) *WebSocketClientConnection {
	return &WebSocketClientConnection{
		eventToSendToCLient:   eventToSendToCLient,
		clientWebsocketSender: clientWebsocketSender,
		wsConnection:          wsConnection,
	}
}

//WebSocketClientConnection is a client-connection accessible through websocket
type WebSocketClientConnection struct {
	eventToSendToCLient   chan event.Event
	clientWebsocketSender ClientWebSocketSender
	wsConnection          websocket.WebsocketConnection
}

//SendEventsToClient sends events to a client
func (clientConnection *WebSocketClientConnection) SendEventsToClient(events []event.Event) error {
	for _, event := range events {
		clientConnection.eventToSendToCLient <- event
	}
	return nil
}

//Close closes the connection. It stops the client-websocket-sender.
func (clientConnection *WebSocketClientConnection) Close() {
	clientConnection.clientWebsocketSender.Stop()
	//closing the websocket connection will stop both client-listener and client-sender
	clientConnection.wsConnection.Close()
}

//NewClientWebSocketListener is a factory for ClientWebSocketListener
func NewClientWebSocketListener(playerID string, wsConnection websocket.WebsocketConnection, server server.Server) *ClientWebSocketListener {
	return &ClientWebSocketListener{
		playerID:       playerID,
		wsConnection:   wsConnection,
		server:         server,
		bufferProvider: bufferProvider,
	}
}

//ClientWebSocketListener is a runnable which listen from incoming websocket messages from a client
type ClientWebSocketListener struct {
	wsConnection   websocket.WebsocketConnection
	server         server.Server
	playerID       string
	bufferProvider func() []event.Event
}

//Run is a blocking loop to listen on incoming websocket events from a client
func (clientWebSocketListener *ClientWebSocketListener) Run() error {
	for {
		//TODO: optimize the reader: I had to create a new array on each read, otherwise, object set during
		//the first array initialization are re-used and override
		eventsFromClient := clientWebSocketListener.bufferProvider()
		if err := clientWebSocketListener.wsConnection.ReadJSON(&eventsFromClient); err != nil {
			clientWebSocketListener.server.UnregisterClient(clientWebSocketListener.playerID)
			return fmt.Errorf("%w", err)
		}
		for _, event := range eventsFromClient {
			event.PlayerID = clientWebSocketListener.playerID
			clientWebSocketListener.server.ReceiveEventFromClient(event)
		}
	}
}

//NewClientWebSocketSender is a factory for ClientWebSocketSender
func NewClientWebSocketSender(wsConnection websocket.WebsocketConnection, eventToSendToClient chan event.Event) *ClientWebSocketSenderImpl {
	return &ClientWebSocketSenderImpl{
		wsConnection:        wsConnection,
		eventToSendToClient: eventToSendToClient,
		quit:                make(chan interface{}),
	}
}

//ClientWebSocketSender is a runnable which send events from the server to the client using the websocket-connection opened by the client
type ClientWebSocketSender interface {
	runner.Runnable
	Stop()
}

//ClientWebSocketSenderImpl is a default implementation of the ClientWebSocketSender
type ClientWebSocketSenderImpl struct {
	wsConnection        websocket.WebsocketConnection
	eventToSendToClient chan event.Event
	quit                chan interface{}
}

//Run is a blocking loop waiting for events from server and to be sent to the client
func (clientWebSocketSender *ClientWebSocketSenderImpl) Run() error {
	for {
		select {
		case <-clientWebSocketSender.quit:
			return nil
		case eventToClient := <-clientWebSocketSender.eventToSendToClient:
			if err := clientWebSocketSender.wsConnection.WriteJSON([]event.Event{eventToClient}); err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	}
}

//Stop stops the client-event-sender
func (clientWebSocketSender *ClientWebSocketSenderImpl) Stop() {
	close(clientWebSocketSender.quit)
}
