package websocketconnector

import (
	"fmt"
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"log"

	"github.com/gorilla/websocket"
)

//RegisterWebSocketServerConnectionToServer creates a new websocket client connection and register it to the server
func RegisterWebSocketServerConnectionToServer(engine client.Engine, url string, quit chan struct{}) error {
	wsConnection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
		return err
	}
	websocketServerConnection := &WebSocketServerConnection{
		engine:                                engine,
		wsConnection:                          wsConnection,
		playerEventQueue:                      make(chan event.Event),
		quit:                                  quit,
		preInitializationEventFromServerQueue: make(chan event.Event, 100),
		intialized:                            false,
	}
	//listen to the events from the server
	go websocketServerConnection.listenServer()
	//start the player-event sender
	go websocketServerConnection.listenPlayer()
	return nil
}

//WebSocketServerConnection is a server-connection accessible through websocket
type WebSocketServerConnection struct {
	engine client.Engine
	// The websocket connection.
	wsConnection                          *websocket.Conn
	quit                                  chan struct{}
	playerEventQueue                      chan event.Event
	preInitializationEventFromServerQueue chan event.Event
	intialized                            bool
	playerID                              string
}

func (connection *WebSocketServerConnection) listenPlayer() {
	//eventsFromPlayer := make([]event.Event, 1)
	for {
		select {
		case /*eventsFromPlayer[0] =*/ eventFromPlayer := <-connection.playerEventQueue:
			if err := connection.wsConnection.WriteJSON([]event.Event{eventFromPlayer}); err != nil {
				fmt.Println(err)
			}
		case <-connection.quit:
			connection.wsConnection.Close()
		}
	}
}

func (connection *WebSocketServerConnection) listenServer() {
	for {
		eventsFromServer := make([]event.Event, 0)
		err := connection.wsConnection.ReadJSON(&eventsFromServer)
		if err != nil {
			log.Println("read:", err)
		} else {
			connection.ReceiveEventsFromServer(0, eventsFromServer)
		}
	}
}

//ReceiveEventsFromServer receive events from a remote server
func (serverConnection *WebSocketServerConnection) ReceiveEventsFromServer(timeFrame uint32, events []event.Event) {
	if serverConnection.intialized {
		serverConnection.engine.ReceiveEventsFromServer(events)
	} else {
		for _, eventFromServer := range events {
			if eventFromServer.Action == "init" {
				//initialize and start the client
				serverConnection.playerID = eventFromServer.PlayerID
				worldMap, _ := eventFromServer.ExtraData["worldMap"].(world.WorldMap)
				otherPlayerStates, _ := eventFromServer.ExtraData["otherPlayers"].(map[string]state.AnimatedElementState)
				otherPlayerStatesClone := make(map[string]state.AnimatedElementState)
				for otherPlayerID, otherPlayerState := range otherPlayerStates {
					otherPlayerStatesClone[otherPlayerID] = otherPlayerState.Clone()
				}
				serverConnection.engine.Initialize(eventFromServer.PlayerID, eventFromServer.State.Clone(), worldMap.Clone(), otherPlayerStatesClone, eventFromServer.TimeFrame)
				serverConnection.engine.StartEngine()
				//propagate events from player
				serverConnection.engine.GetPlayer().RegisterListener(serverConnection.playerEventQueue)
				//process all previous events
				numberOfPreInitializationEvents := len(serverConnection.preInitializationEventFromServerQueue)
				if numberOfPreInitializationEvents > 0 {
					preInitializationEvents := make([]event.Event, numberOfPreInitializationEvents)
					for i := 0; i < numberOfPreInitializationEvents; i++ {
						preInitializationEvents[i] = <-serverConnection.preInitializationEventFromServerQueue
					}
					serverConnection.engine.ReceiveEventsFromServer(preInitializationEvents)
				}
				//change the state
				serverConnection.intialized = true
			} else {
				serverConnection.preInitializationEventFromServerQueue <- eventFromServer
			}
		}
	}
}
