package main

import (
	"fmt"
	"francoisgergaud/3dGame/client"
	localServerConnector "francoisgergaud/3dGame/client/connector/local/impl"
	"francoisgergaud/3dGame/client/consolemanager"
	consoleManagerImpl "francoisgergaud/3dGame/client/consolemanager/impl"
	clientImpl "francoisgergaud/3dGame/client/impl"
	serverImpl "francoisgergaud/3dGame/server/impl"

	"github.com/gdamore/tcell"
)

//Game defines the game entity.
type Game struct {
	engine              client.Engine
	consoleEventManager consolemanager.ConsoleEventManager
	quit                chan struct{}
}

//Start the game.
func (game *Game) Start() {
	go game.consoleEventManager.Listen()
	//wait until quit
	<-game.quit
	<-game.engine.GetShutdown()
}

//InitGame initializes a game.
func InitGame(screen tcell.Screen) (*Game, error) {
	quit := make(chan struct{})
	worldUpdateRate := 20 //world-update frequency, for both client and server
	server, err := serverImpl.NewServer(worldUpdateRate, quit)
	if err != nil {
		return nil, fmt.Errorf("error while instantiating the server: %w", err)
	}
	engineConfiguration := &client.Configuration{
		FrameRate:                  20,
		WorlUpdateRate:             worldUpdateRate,
		ScreenHeight:               40,
		ScreenWidth:                120,
		PlayerFieldOfViewAngle:     0.4,
		Visibility:                 20.0,
		GradientRSFirst:            1.0,
		GradientRSMultiplicator:    2.0,
		GradientRSLimit:            10.0,
		GradientRSWallStartColor:   255,
		GradientRSWallEndColor:     240,
		GradientRSBackgroundRange:  []float32{0.5, 0.55, 0.65},
		GradientRSBackgroundColors: []int{63, 58, 64, 70},
		QuitChannel:                quit,
	}
	consoleEventManager := consoleManagerImpl.NewConsoleEventManager(screen, quit)
	engine, err := clientImpl.NewEngine(screen, consoleEventManager, engineConfiguration)
	if err != nil {
		return nil, fmt.Errorf("Error while initializing engine: %w", err)
	}
	//client-server initialization
	serverConnection := localServerConnector.NewLocalServerConnection(engine, server, quit)
	serverConnection.ConnectToServer()

	//starts the game
	return &Game{
		engine:              engine,
		consoleEventManager: consoleEventManager,
		quit:                quit,
	}, nil

}
