package impl

import (
	"francoisgergaud/3dGame/client/consolemanager"
	"francoisgergaud/3dGame/client/player"

	"github.com/gdamore/tcell"
)

//NewConsoleEventManager builds a new ConsoleEventManagerImpl.
func NewConsoleEventManager(screen tcell.Screen, quit chan struct{}) consolemanager.ConsoleEventManager {
	return &ConsoleEventManagerImpl{
		screen:      screen,
		quitChannel: quit,
	}
}

//ConsoleEventManagerImpl is the implementation of the ConsoleEventManager interface.
type ConsoleEventManagerImpl struct {
	screen      tcell.Screen
	player      player.Player
	quitChannel chan struct{}
}

//SetPlayer set the player the console-event will be sent to
func (consoleEventManager *ConsoleEventManagerImpl) SetPlayer(player player.Player) {
	consoleEventManager.player = player
}

//Listen listens to the events emit by the terminal.
func (consoleEventManager *ConsoleEventManagerImpl) Listen() {
	for {
		ev := consoleEventManager.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyEnter:
				close(consoleEventManager.quitChannel)
				return
			default:
				if consoleEventManager.player != nil {
					consoleEventManager.player.Action(ev)
				}
			}
		case *tcell.EventResize:
			consoleEventManager.screen.Sync()
		}
	}
}
