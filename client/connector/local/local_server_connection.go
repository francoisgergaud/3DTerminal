package localconnector

import (
	"francoisgergaud/3dGame/common/event"
)

//LocalServerConnection is a client-side connection to a server.
type LocalServerConnection interface {
	ReceiveEventsFromServer(timeFrame uint32, events []event.Event)
}
