package websocketconnector

import (
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server"

	"github.com/gorilla/websocket"
)

//RegisterWebSocketClientConnectionToServer creates a new websocket client connection and register it to the server
func RegisterWebSocketClientConnectionToServer(wsConnection *websocket.Conn, server server.Server) {
	websocketClientConnection := &WebSocketClientConnection{
		server:       server,
		wsConnection: wsConnection,
		send:         make(chan event.Event),
	}
	server.RegisterPlayer(websocketClientConnection)
}

//WebSocketClientConnection is a client-connection accessible through websocket
type WebSocketClientConnection struct {
	server server.Server
	// The websocket connection.
	wsConnection *websocket.Conn
	send         chan event.Event
}

//SendEventsToClient sends events to a client
func (clientConnection *WebSocketClientConnection) SendEventsToClient(timeFrame uint32, events []event.Event) {
	for _, event := range events {
		clientConnection.send <- event
	}
}

func (clientConnection *WebSocketClientConnection) readFromWebSocket() {
	defer func() {
		clientConnection.wsConnection.Close()
	}()
	for {
		eventsFromClient := make([]event.Event, 0)
		clientConnection.wsConnection.ReadJSON(eventsFromClient)
		for _, event := range eventsFromClient {
			clientConnection.server.ReceiveEventFromClient(event)
		}
	}
}

func (clientConnection *WebSocketClientConnection) writeToWebSocket() {
	defer func() {
		clientConnection.wsConnection.Close()
	}()
	for {
		select {
		case event := <-clientConnection.send:
			clientConnection.wsConnection.WriteJSON(event)
		}
	}
}
