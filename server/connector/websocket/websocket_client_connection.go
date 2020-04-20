package websocketconnector

import (
	"fmt"
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

	go websocketClientConnection.writeToWebSocket()
	websocketClientConnection.playerID = server.RegisterPlayer(websocketClientConnection)
	go websocketClientConnection.readFromWebSocket()
}

//WebSocketClientConnection is a client-connection accessible through websocket
type WebSocketClientConnection struct {
	server server.Server
	// The websocket connection.
	wsConnection *websocket.Conn
	send         chan event.Event
	playerID     string
}

//SendEventsToClient sends events to a client
func (clientConnection *WebSocketClientConnection) SendEventsToClient(timeFrame uint32, events []event.Event) {
	for _, event := range events {
		clientConnection.send <- event
	}
}

func (clientConnection *WebSocketClientConnection) readFromWebSocket() {
	eventsFromClient := make([]event.Event, 0)
	for {
		if err := clientConnection.wsConnection.ReadJSON(&eventsFromClient); err != nil {
			fmt.Println(err)
		}
		for _, event := range eventsFromClient {
			event.PlayerID = clientConnection.playerID
			clientConnection.server.ReceiveEventFromClient(event)
		}
	}
}

func (clientConnection *WebSocketClientConnection) writeToWebSocket() {
	for {
		//eventsToClient := make([]event.Event, 1)
		select {
		case /*eventsToClient[0] =*/ eventToClient := <-clientConnection.send:
			if err := clientConnection.wsConnection.WriteJSON([]event.Event{eventToClient}); err != nil {
				fmt.Println(err)
			}
		}
	}
}
