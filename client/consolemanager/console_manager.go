package consolemanager

import "francoisgergaud/3dGame/client/player"

//ConsoleEventManager is a listner for events coming from the console.
type ConsoleEventManager interface {
	SetPlayer(player.Player)
	Run() error
}
