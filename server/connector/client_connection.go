package connector

import (
	"francoisgergaud/3dGame/common/event"
)

//ClientConnection is a server-side connection to a client.
type ClientConnection interface {
	SendEventsToClient(events []event.Event)
	Close()
}
