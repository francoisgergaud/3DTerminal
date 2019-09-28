package event

import (
	"francoisgergaud/3dGame/internal/testutils"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gdamore/tcell"
)

func TestListenExitEvent(t *testing.T) {
	mockCharacter := new(testutils.MockCharacter)
	mockScreen := new(testutils.MockScreen)
	quit := make(chan struct{})
	keyboardEvent := tcell.NewEventKey(tcell.KeyEscape, ' ', 0)
	mockScreen.On("PollEvent").Return(keyboardEvent)
	consoleEventManager := NewConsoleEventManager(mockScreen, mockCharacter, quit)
	consoleEventManager.Listen()
}

func TestListenPlayerMoveEvent(t *testing.T) {
	mockCharacter := new(testutils.MockCharacter)
	mockScreen := new(testutils.MockScreen)
	quit := make(chan struct{})
	upArrowEvent := tcell.NewEventKey(tcell.KeyUp, ' ', 0)
	quitEvent := tcell.NewEventKey(tcell.KeyEscape, ' ', 0)
	mockScreen.On("PollEvent").Return(upArrowEvent).Once()
	mockScreen.On("PollEvent").Return(quitEvent).Once()
	mockCharacter.On("Action", upArrowEvent)
	consoleEventManager := NewConsoleEventManager(mockScreen, mockCharacter, quit)
	consoleEventManager.Listen()
	_, status := <-quit
	assert.Falsef(t, status, "quit channel status invalid.")
}
