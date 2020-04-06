package impl

import (
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server"
)

//NewLocalServerConnection is the local-server-connection factory
func NewLocalServerConnection(engine client.Engine, server server.Server, quit chan struct{}) *LocalServerConnectionImpl {
	return &LocalServerConnectionImpl{
		engine:                                engine,
		playerEventQueue:                      make(chan event.Event),
		preInitializationEventFromServerQueue: make(chan event.Event, 100),
		quit:                                  quit,
		server:                                server,
		intialized:                            false,
	}
}

//LocalServerConnectionImpl is an implementation of a client connection to a local-server
type LocalServerConnectionImpl struct {
	engine                                client.Engine
	playerEventQueue                      chan event.Event
	preInitializationEventFromServerQueue chan event.Event
	quit                                  chan struct{}
	server                                server.Server
	intialized                            bool
	playerID                              string
}

//ReceiveEventsFromServer receive events from a local server
func (serverConnection *LocalServerConnectionImpl) ReceiveEventsFromServer(timeFrame uint32, events []event.Event) {
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

//ConnectToServer connects to the server and listen to the events from:
// - the player to send them to the server
// - the server to update the local environment
func (serverConnection *LocalServerConnectionImpl) ConnectToServer() {
	serverConnection.server.RegisterPlayer(serverConnection)
	go func() {
		for {
			select {
			case eventFromPlayer := <-serverConnection.playerEventQueue:
				eventFromPlayer.PlayerID = serverConnection.playerID
				serverConnection.server.ReceiveEventFromClient(eventFromPlayer)
			case <-serverConnection.quit:
				return
			}
		}
	}()
}

//SendEventsToClient sends a list of events to the client
func (serverConnection *LocalServerConnectionImpl) SendEventsToClient(timeFrame uint32, events []event.Event) {
	serverConnection.ReceiveEventsFromServer(timeFrame, events)
}
