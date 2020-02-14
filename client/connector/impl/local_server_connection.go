package impl

import (
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server/connector"
)

//LocalServerConnection is an implementation of a client connection to a local-server
type LocalServerConnection struct {
	Engine           client.Engine
	ClientConnection connector.ClientConnection
}

//ReceiveEventsFromServer receive events from a local server
func (clientConnection *LocalServerConnection) ReceiveEventsFromServer(timeFrame uint32, events []event.Event) {
	clientConnection.Engine.ReceiveEventsFromServer(events)
}

//SendEventsToServer sends event to a local server
func (clientConnection *LocalServerConnection) SendEventsToServer(timeFrame uint32, events []event.Event) {
	clientConnection.ClientConnection.ReceiveEventsFromClient(timeFrame, events)
}
