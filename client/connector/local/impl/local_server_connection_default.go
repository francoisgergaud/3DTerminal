package impl

import (
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server"
)

//NewLocalServerConnection is the local-server-connection factory
func NewLocalServerConnection(engine client.Engine, server server.Server, quit chan struct{}) *LocalServerConnectionImpl {
	localServerConnection := &LocalServerConnectionImpl{
		engine: engine,
		quit:   quit,
		server: server,
	}
	localServerConnection.playerID = localServerConnection.server.RegisterPlayer(localServerConnection)
	localServerConnection.engine.SetConnectionToServer(localServerConnection)
	return localServerConnection
}

//LocalServerConnectionImpl is an implementation of a client connection to a local-server
type LocalServerConnectionImpl struct {
	engine   client.Engine
	quit     chan struct{}
	server   server.Server
	playerID string
}

//ConnectToServer connects to the server and listen to the events from:
// - the player to send them to the server
// - the server to update the local environment
// func (serverConnection *LocalServerConnectionImpl) ConnectToServer() {
// 	serverConnection.playerID = serverConnection.server.RegisterPlayer(serverConnection)
// 	serverConnection.engine.SetConnectionToServer(serverConnection)
// 	go func() {
// 		for {
// 			select {
// 			case eventFromPlayer := <-serverConnection.playerEventQueue:
// 				eventFromPlayer.PlayerID = serverConnection.playerID
// 				serverConnection.server.ReceiveEventFromClient(eventFromPlayer)
// 			case <-serverConnection.quit:
// 				return
// 			}
// 		}
// 	}()
// }

func (serverConnection *LocalServerConnectionImpl) NotifyServer(events []event.Event) {
	for _, event := range events {
		event.PlayerID = serverConnection.playerID
		serverConnection.server.ReceiveEventFromClient(event)
	}
}

func (serverConnection *LocalServerConnectionImpl) Disconnect() {}

//SendEventsToClient sends a list of events to the client
func (serverConnection *LocalServerConnectionImpl) SendEventsToClient(events []event.Event) {
	serverConnection.engine.ReceiveEventsFromServer(events)
}
