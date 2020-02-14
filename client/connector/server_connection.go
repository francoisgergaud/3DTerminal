package connector

import "francoisgergaud/3dGame/common/event"

//ServerConnection is a client-side connection to a server.
type ServerConnection interface {
	//send a list of events to the server
	SendEventsToServer(timeFrame uint32, events []event.Event)
	//receive a list of events from the sever
	ReceiveEventsFromServer(timeFrame uint32, events []event.Event)
}
