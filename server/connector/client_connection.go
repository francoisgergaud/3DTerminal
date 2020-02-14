package connector

import "francoisgergaud/3dGame/common/event"

//ClientConnection is a server-side connection to a client.
type ClientConnection interface {
	//receive a list of events froma client
	ReceiveEventsFromClient(timeFrame uint32, events []event.Event)
	//send a list of events to the client
	SendEventsToClient(timeFrame uint32, events []event.Event)
}
