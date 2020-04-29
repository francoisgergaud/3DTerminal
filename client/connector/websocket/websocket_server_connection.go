package websocketconnector

import (
	"fmt"
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/common/event"
	"log"

	"github.com/gorilla/websocket"
)

//NewWebSocketServerConnection creates a new websocket client connection and register it to the server
func NewWebSocketServerConnection(engine client.Engine, url string, quit chan struct{}) error {
	wsConnection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Could not dial server websocket on :"+url, err)
		return err
	}
	websocketServerConnection := &WebSocketServerConnection{
		engine:       engine,
		wsConnection: wsConnection,
		quit:         quit,
	}
	websocketServerConnection.engine.SetConnectionToServer(websocketServerConnection)
	//listen to the events from the server
	go websocketServerConnection.listenServer()
	//start the player-event sender
	//go websocketServerConnection.listenPlayer()
	return nil
}

//WebSocketServerConnection is a server-connection accessible through websocket
type WebSocketServerConnection struct {
	engine client.Engine
	// The websocket connection.
	wsConnection *websocket.Conn
	quit         chan struct{}
	playerID     string
}

// func (connection *WebSocketServerConnection) listenPlayer() {
// 	//eventsFromPlayer := make([]event.Event, 1)
// 	for {
// 		select {
// 		case /*eventsFromPlayer[0] =*/ eventFromPlayer := <-connection.playerEventQueue:
// 			if err := connection.wsConnection.WriteJSON([]event.Event{eventFromPlayer}); err != nil {
// 				fmt.Println(err)
// 			}
// 		case <-connection.quit:
// 			connection.wsConnection.Close()
// 		}
// 	}
//}

func (connection *WebSocketServerConnection) NotifyServer(events []event.Event) {
	if err := connection.wsConnection.WriteJSON(events); err != nil {
		fmt.Println(err)
	}
}

func (connection *WebSocketServerConnection) Disconnect() {
	connection.wsConnection.Close()
}

func (connection *WebSocketServerConnection) listenServer() {
	for {
		eventsFromServer := make([]event.Event, 0)
		err := connection.wsConnection.ReadJSON(&eventsFromServer)
		if err != nil {
			log.Println("read:", err)
		} else {
			connection.engine.ReceiveEventsFromServer(eventsFromServer)
		}
	}
}
