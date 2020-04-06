package main

import (
	"testing"
	"time"

	testClient "francoisgergaud/3dGame/internal/testutils/client"
	testConsoleManager "francoisgergaud/3dGame/internal/testutils/client/consolemanager"
)

func TestStart(t *testing.T) {
	engine := new(testClient.MockEngine)
	consoleEventManager := new(testConsoleManager.MockConsoleEventManager)
	engine.On("StartEngine")
	consoleEventManager.On("Listen")
	quit := make(chan struct{})
	engineShutdown := make(chan interface{})
	engine.On("GetShutdown").Return(engineShutdown)
	game := Game{
		engine:              engine,
		consoleEventManager: consoleEventManager,
		quit:                quit,
	}
	timer := time.NewTimer(time.Microsecond)
	go func() {
		<-timer.C
		close(quit)
		close(engineShutdown)
	}()
	game.Start()

}

// func TestInitGame(t *testing.T) {
// 	screen := new(testtcell.MockScreen)
// 	screen.On("Clear")
// 	game, err := InitGame(screen)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, game)
// }
