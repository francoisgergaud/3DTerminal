package websocketconnector

import (
	"fmt"
	"francoisgergaud/3dGame/client"
	websocket "francoisgergaud/3dGame/common/connector"
	"francoisgergaud/3dGame/common/event"
)

//NewWebSocketServerConnection creates a new websocket client connection and register it to the server
func NewWebSocketServerConnection(engine client.Engine, url string, dialer WebsocketDialer, quit chan<- interface{}) (*WebSocketServerConnection, error) {
	wsConnection, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("Could not dial server websocket on :"+url+", %w", err)
	}
	websocketServerConnection := &WebSocketServerConnection{
		engine:       engine,
		wsConnection: wsConnection,
		quit:         quit,
	}
	websocketServerConnection.engine.ConnectToServer(websocketServerConnection)
	//listen to the events from the server
	return websocketServerConnection, nil
}

//WebSocketServerConnection is a server-connection accessible through websocket
type WebSocketServerConnection struct {
	engine client.Engine
	// The websocket connection.
	wsConnection websocket.WebsocketConnection
	playerID     string
	quit         chan<- interface{}
}

//NotifyServer sends an event to s server
func (connection *WebSocketServerConnection) NotifyServer(events []event.Event) error {
	if err := connection.wsConnection.WriteJSON(events); err != nil {
		return fmt.Errorf("quit client-websocket sender because of write-error: %w", err)
	}
	return nil
}

//Disconnect disconnect the client websocket to the server
func (connection *WebSocketServerConnection) Disconnect() {
	connection.wsConnection.Close()
}

//Run is a blocking loop listening events from server
func (connection *WebSocketServerConnection) Run() error {
	for {
		eventsFromServer := make([]event.Event, 0)
		err := connection.wsConnection.ReadJSON(&eventsFromServer)
		if err != nil {
			close(connection.quit)
			return fmt.Errorf("quit client-websocket listener because of read-error: %w", err)
		}
		connection.engine.ReceiveEventsFromServer(eventsFromServer)
	}
}
