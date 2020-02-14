package main

import (
	"testing"

	testClient "francoisgergaud/3dGame/internal/testutils/client"
	testConsoleManager "francoisgergaud/3dGame/internal/testutils/client/consolemanager"
	testtcell "francoisgergaud/3dGame/internal/testutils/tcell"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	engine := new(testClient.MockEngine)
	consoleEventManager := new(testConsoleManager.MockConsoleEventManager)
	engine.On("StartEngine")
	consoleEventManager.On("Listen")
	game := Game{
		engine:              engine,
		consoleEventManager: consoleEventManager,
	}
	game.Start()
}

func TestInitGame(t *testing.T) {
	screen := new(testtcell.MockScreen)
	game, err := InitGame(screen)
	assert.Nil(t, err)
	assert.NotNil(t, game)
}
