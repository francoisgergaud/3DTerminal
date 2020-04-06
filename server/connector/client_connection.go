package connector

import (
	"francoisgergaud/3dGame/common/event"
)

//ClientConnection is a server-side connection to a client.
type ClientConnection interface {
	SendEventsToClient(timeFrame uint32, events []event.Event)
}
