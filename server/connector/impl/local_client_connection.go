package impl

import (
	"francoisgergaud/3dGame/client/connector"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server"
)

//LocalClientConnection represents local-connection
type LocalClientConnection struct {
	Server           server.Server
	ServerConnection connector.ServerConnection
}

//SendEventsToClient sends a list of events to the client
func (serverConnection *LocalClientConnection) SendEventsToClient(timeFrame uint32, events []event.Event) {
	serverConnection.ServerConnection.ReceiveEventsFromServer(timeFrame, events)
}

//ReceiveEventsFromClient receive events from a client
func (serverConnection *LocalClientConnection) ReceiveEventsFromClient(timeFrame uint32, events []event.Event) {
	serverConnection.Server.ReceiveEventsFromClient(events)
}
