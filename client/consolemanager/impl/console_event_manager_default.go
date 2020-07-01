package impl

import (
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/client/consolemanager"

	"github.com/gdamore/tcell"
)

//NewConsoleEventManager builds a new ConsoleEventManagerImpl.
func NewConsoleEventManager(screen tcell.Screen, quit chan<- interface{}) consolemanager.ConsoleEventManager {
	return &ConsoleEventManagerImpl{
		screen:      screen,
		quitChannel: quit,
	}
}

//ConsoleEventManagerImpl is the implementation of the ConsoleEventManager interface.
type ConsoleEventManagerImpl struct {
	screen      tcell.Screen
	engine      client.Engine
	quitChannel chan<- interface{}
}

//SetPlayer set the player the console-event will be sent to
func (consoleEventManager *ConsoleEventManagerImpl) SetPlayer(engine client.Engine) {
	consoleEventManager.engine = engine
}

//Run is a blocking loop listening to the events emited by the terminal.
func (consoleEventManager *ConsoleEventManagerImpl) Run() error {
	for {
		ev := consoleEventManager.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				close(consoleEventManager.quitChannel)
				return nil
			default:
				if consoleEventManager.engine != nil {
					consoleEventManager.engine.Action(ev)
				}
			}
		case *tcell.EventResize:
			consoleEventManager.screen.Sync()
		}
	}
}
