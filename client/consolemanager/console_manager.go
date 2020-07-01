package consolemanager

import (
	"francoisgergaud/3dGame/client"
)

//ConsoleEventManager is a listner for events coming from the console.
type ConsoleEventManager interface {
	SetPlayer(client.Engine)
	Run() error
}
