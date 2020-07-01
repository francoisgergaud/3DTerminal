package impl

import (
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server"
)

//NewLocalServerConnection is the local-server-connection factory
func NewLocalServerConnection(engine client.Engine, server server.Server, quit <-chan interface{}) {
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
	quit     <-chan interface{}
	server   server.Server
	playerID string
}

//NotifyServer send an event to the local server
func (serverConnection *LocalServerConnectionImpl) NotifyServer(events []event.Event) error {
	for _, event := range events {
		event.PlayerID = serverConnection.playerID
		eventClone, err := event.Clone()
		if err != nil {
			return err
		}
		serverConnection.server.ReceiveEventFromClient(*eventClone)
	}
	return nil
}

//Disconnect does not do anything for a local-connection
func (serverConnection *LocalServerConnectionImpl) Disconnect() {}

//SendEventsToClient sends a list of events to the client
func (serverConnection *LocalServerConnectionImpl) SendEventsToClient(events []event.Event) error {
	eventsClone := make([]event.Event, len(events))
	for i, eventToClone := range events {
		eventClone, err := eventToClone.Clone()
		if err != nil {
			return err
		}
		eventsClone[i] = *eventClone
	}
	serverConnection.engine.ReceiveEventsFromServer(eventsClone)
	return nil
}

//Close closes the connection (doesn't do anything for a local-connection)
func (serverConnection *LocalServerConnectionImpl) Close() {
}
