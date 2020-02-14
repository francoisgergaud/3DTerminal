package impl

import (
	"testing"

	testplayer "francoisgergaud/3dGame/internal/testutils/client/player"
	testtcell "francoisgergaud/3dGame/internal/testutils/tcell"

	"github.com/stretchr/testify/assert"

	"github.com/gdamore/tcell"
)

func TestListenExitEvent(t *testing.T) {
	mockScreen := new(testtcell.MockScreen)
	quit := make(chan struct{})
	keyboardEvent := tcell.NewEventKey(tcell.KeyEscape, ' ', 0)
	mockScreen.On("PollEvent").Return(keyboardEvent)
	consoleEventManager := NewConsoleEventManager(mockScreen, quit)
	consoleEventManager.Listen()
}

func TestListenPlayerMoveEvent(t *testing.T) {
	mockPlayer := new(testplayer.MockPlayer)
	mockScreen := new(testtcell.MockScreen)
	quit := make(chan struct{})
	upArrowEvent := tcell.NewEventKey(tcell.KeyUp, ' ', 0)
	quitEvent := tcell.NewEventKey(tcell.KeyEscape, ' ', 0)
	mockScreen.On("PollEvent").Return(upArrowEvent).Once()
	mockScreen.On("PollEvent").Return(quitEvent).Once()
	mockPlayer.On("Action", upArrowEvent)
	consoleEventManager := NewConsoleEventManager(mockScreen, quit)
	consoleEventManager.SetPlayer(mockPlayer)
	consoleEventManager.Listen()
	_, status := <-quit
	assert.Falsef(t, status, "quit channel status invalid.")
}
