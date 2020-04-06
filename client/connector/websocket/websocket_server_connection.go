package websocketconnector

import (
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/common/event"
	"log"

	"github.com/gorilla/websocket"
)

//RegisterWebSocketServerConnectionToServer creates a new websocket client connection and register it to the server
func RegisterWebSocketServerConnectionToServer(engine client.Engine, quit chan struct{}) {
	websocketServerConnection := &WebSocketServerConnection{
		engine:           engine,
		wsConnection:     nil,
		eventsFromPlayer: make(chan event.Event),
		quit:             quit,
	}
	engine.GetPlayer().RegisterListener(websocketServerConnection.eventsFromPlayer)
}

//WebSocketServerConnection is a server-connection accessible through websocket
type WebSocketServerConnection struct {
	engine client.Engine
	// The websocket connection.
	wsConnection     *websocket.Conn
	eventsFromPlayer chan event.Event
	quit             chan struct{}
}

func (connection *WebSocketServerConnection) startPlayerEventSender() {
	//start the player-event sender
	go func() {
		for {
			select {
			case eventFromPlayer := <-connection.eventsFromPlayer:
				events := make([]event.Event, 0)
				events = append(events, eventFromPlayer)
				connection.wsConnection.WriteJSON(events)
			case <-connection.quit:
				connection.wsConnection.Close()
			}

		}
	}()
}

func (connection *WebSocketServerConnection) start(url string) {
	wsConnection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	connection.wsConnection = wsConnection
	defer connection.wsConnection.Close()
	//start the player-event sender
	connection.startPlayerEventSender()
	//create the websocket and lister
	go func() {
		eventsFromServer := make([]event.Event, 0)
		for {
			err := connection.wsConnection.ReadJSON(eventsFromServer)
			if err != nil {
				log.Println("read:", err)
				return
			}
			connection.engine.ReceiveEventsFromServer(eventsFromServer)
		}
	}()
}
