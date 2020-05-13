package connector

import "francoisgergaud/3dGame/common/event"

type ServerConnector interface {
	NotifyServer([]event.Event) error
	Disconnect()
}
