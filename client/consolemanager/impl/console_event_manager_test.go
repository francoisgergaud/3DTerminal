package impl

import (
	"testing"

	testplayer "francoisgergaud/3dGame/internal/testutils/client/player"
	testtcell "francoisgergaud/3dGame/internal/testutils/tcell"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gdamore/tcell"
)

func TestListenExitEvent(t *testing.T) {
	mockScreen := new(testtcell.MockScreen)
	quit := make(chan interface{})
	keyboardEvent := tcell.NewEventKey(tcell.KeyEscape, ' ', 0)
	mockScreen.On("PollEvent").Return(keyboardEvent)
	consoleEventManager := NewConsoleEventManager(mockScreen, quit)
	consoleEventManager.Run()
	mock.AssertExpectationsForObjects(t, mockScreen)
}

func TestListenPlayerMoveEvent(t *testing.T) {
	mockPlayer := new(testplayer.MockPlayer)
	mockScreen := new(testtcell.MockScreen)
	quit := make(chan interface{})
	upArrowEvent := tcell.NewEventKey(tcell.KeyUp, ' ', 0)
	quitEvent := tcell.NewEventKey(tcell.KeyEscape, ' ', 0)
	mockScreen.On("PollEvent").Return(upArrowEvent).Once()
	mockScreen.On("PollEvent").Return(quitEvent).Once()
	mockPlayer.On("Action", upArrowEvent)
	consoleEventManager := NewConsoleEventManager(mockScreen, quit)
	consoleEventManager.SetPlayer(mockPlayer)
	consoleEventManager.Run()
	_, status := <-quit
	assert.Falsef(t, status, "quit channel status invalid.")
	mock.AssertExpectationsForObjects(t, mockScreen, mockPlayer)
}
