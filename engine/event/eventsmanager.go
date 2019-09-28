package event

import (
	"francoisgergaud/3dGame/environment"

	"github.com/gdamore/tcell"
)

//ConsoleEventManager is a listner for events coming from the console.
type ConsoleEventManager interface {
	Listen()
}

//NewConsoleEventManager builds a new ConsoleEventManagerImpl.
func NewConsoleEventManager(screen tcell.Screen, player environment.Character, quit chan struct{}) ConsoleEventManager {
	return &ConsoleEventManagerImpl{
		screen: screen,
		player: player,
		quit:   quit,
	}
}

//ConsoleEventManagerImpl is the implementation of the ConsoleEventManager interface.
type ConsoleEventManagerImpl struct {
	screen tcell.Screen
	player environment.Character
	quit   chan struct{}
}

//Listen listens to the events emit by the terminal.
func (consoleEventManager *ConsoleEventManagerImpl) Listen() {
	for {
		ev := consoleEventManager.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyEnter:
				close(consoleEventManager.quit)
				return
			default:
				consoleEventManager.player.Action(ev)
			}
		case *tcell.EventResize:
			consoleEventManager.screen.Sync()
		}
	}
}
