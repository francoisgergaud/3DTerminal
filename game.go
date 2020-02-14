package main

import (
	"fmt"
	"francoisgergaud/3dGame/client"
	serverConnectorImpl "francoisgergaud/3dGame/client/connector/impl"
	"francoisgergaud/3dGame/client/consolemanager"
	consoleManagerImpl "francoisgergaud/3dGame/client/consolemanager/impl"
	clientImpl "francoisgergaud/3dGame/client/impl"
	clientconnectorImpl "francoisgergaud/3dGame/server/connector/impl"
	serverImpl "francoisgergaud/3dGame/server/impl"

	"github.com/gdamore/tcell"
)

//Game defines the game entity.
type Game struct {
	engine              client.Engine
	consoleEventManager consolemanager.ConsoleEventManager
}

//Start the game.
func (game *Game) Start() {
	go game.consoleEventManager.Listen()
	game.engine.StartEngine()
}

//InitGame initializes a game.
func InitGame(screen tcell.Screen) (*Game, error) {
	clientConnection := &clientconnectorImpl.LocalClientConnection{}
	serverConnection := &serverConnectorImpl.LocalServerConnection{}
	quit := make(chan struct{})
	server, err := serverImpl.NewServer(quit)
	if err != nil {
		return nil, fmt.Errorf("error while instantiating the server: %w", err)
	}
	playerID, worldMap, playerState, otherPlayerConfigurations := server.RegisterPlayer(clientConnection)
	engineConfiguration := &client.Configuration{
		FrameRate:                  20,
		WorlUpdateRate:             40,
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
		WorldMap:                   worldMap,
		PlayerID:                   playerID,
		PlayerConfiguration:        playerState,
		OtherPlayerConfigurations:  otherPlayerConfigurations,
		QuitChannel:                quit,
	}
	engine, err := clientImpl.NewEngine(screen, engineConfiguration, serverConnection)
	if err != nil {
		return nil, fmt.Errorf("Error while initializing engine: %w", err)
	}
	consoleEventManager := consoleManagerImpl.NewConsoleEventManager(screen, quit)
	consoleEventManager.SetPlayer(engine.GetPlayer())
	clientConnection.Server = server
	clientConnection.ServerConnection = serverConnection
	serverConnection.Engine = engine
	serverConnection.ClientConnection = clientConnection
	return &Game{
		engine:              engine,
		consoleEventManager: consoleEventManager,
	}, nil

}
