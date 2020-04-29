package localconnector

import (
	"francoisgergaud/3dGame/client/connector"
	"francoisgergaud/3dGame/common/event"
)

//LocalServerConnection is a client-side connection to a server.
type LocalServerConnection interface {
	connector.ServerConnector
	ReceiveEventsFromServer(events []event.Event)
}
