package impl

import (
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server"
)

//NewLocalServerConnection is the local-server-connection factory
func NewLocalServerConnection(engine client.Engine, server server.Server, quit chan struct{}) {
	localServerConnection := &LocalServerConnectionImpl{
		engine: engine,
		quit:   quit,
		server: server,
	}
	localServerConnection.playerID = localServerConnection.server.RegisterPlayer(localServerConnection)
	localServerConnection.engine.ConnectToServer(localServerConnection)
}

//LocalServerConnectionImpl is an implementation of a client connection to a local-server
type LocalServerConnectionImpl struct {
	engine   client.Engine
	quit     chan struct{}
	server   server.Server
	playerID string
}

//NotifyServer send an event to the local server
func (serverConnection *LocalServerConnectionImpl) NotifyServer(events []event.Event) error {
	for _, event := range events {
		event.PlayerID = serverConnection.playerID
		serverConnection.server.ReceiveEventFromClient(event)
	}
	return nil
}

//Disconnect does not do anything for a local-connection
func (serverConnection *LocalServerConnectionImpl) Disconnect() {}

//SendEventsToClient sends a list of events to the client
func (serverConnection *LocalServerConnectionImpl) SendEventsToClient(events []event.Event) {
	serverConnection.engine.ReceiveEventsFromServer(events)
}
