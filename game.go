package main

import (
	"flag"
	"fmt"
	"francoisgergaud/3dGame/client"
	localServerConnector "francoisgergaud/3dGame/client/connector/local/impl"
	clientwebsocketconnector "francoisgergaud/3dGame/client/connector/websocket"
	"francoisgergaud/3dGame/client/consolemanager"
	consoleManagerImpl "francoisgergaud/3dGame/client/consolemanager/impl"
	clientImpl "francoisgergaud/3dGame/client/impl"
	"francoisgergaud/3dGame/server"
	serverImpl "francoisgergaud/3dGame/server/impl"
	webserver "francoisgergaud/3dGame/server/net"
	"time"

	"github.com/gdamore/tcell"
)

//InitGame initializes a game.
func InitGame(screen tcell.Screen) error {
	var mode = flag.String("mode", "local", "possible mode: 'local', 'remote', 'remoteClient', 'remoteServer'")
	var remoteAddress = flag.String("address", "127.0.0.1:9836", "remote-server host-port")
	var serverPort = flag.String("port", "9836", "remote-server host-port")
	flag.Parse()
	quit := make(chan struct{})
	worldUpdateRate := 20 //world-update frequency, for both client and server
	consoleEventManager := consoleManagerImpl.NewConsoleEventManager(screen, quit)
	go consoleEventManager.Listen()
	var err error
	var engine client.Engine
	var server server.Server
	if *mode == "local" {
		server, err = serverImpl.NewServer(worldUpdateRate, quit)
		if err != nil {
			return fmt.Errorf("error while instantiating the server: %w", err)
		}
		engine, err = createClient(screen, quit, worldUpdateRate, consoleEventManager)
		if err != nil {
			return fmt.Errorf("Error while initializing engine: %w", err)
		}
		serverConnection := localServerConnector.NewLocalServerConnection(engine, server, quit)
		serverConnection.ConnectToServer()
	} else if *mode == "remote" {
		server, err = serverImpl.NewServer(worldUpdateRate, quit)
		if err != nil {
			return fmt.Errorf("error while instantiating the server: %w", err)
		}
		engine, err = createClient(screen, quit, worldUpdateRate, consoleEventManager)
		if err != nil {
			return fmt.Errorf("Error while initializing engine: %w", err)
		}
		serverURL := "localhost:" + *serverPort
		webServer := webserver.NewWebServer(server, serverURL)
		go webServer.Start()
		time.Sleep(time.Millisecond)
		clientwebsocketconnector.RegisterWebSocketServerConnectionToServer(engine, "ws://"+serverURL+"/join", quit)
	} else if *mode == "remoteClient" {
		engine, err = createClient(screen, quit, worldUpdateRate, consoleEventManager)
		if err != nil {
			return fmt.Errorf("Error while initializing engine: %w", err)
		}
		clientwebsocketconnector.RegisterWebSocketServerConnectionToServer(engine, "ws://"+*remoteAddress+"/join", quit)
	} else if *mode == "remoteServer" {
		server, err = serverImpl.NewServer(worldUpdateRate, quit)
		if err != nil {
			return fmt.Errorf("error while instantiating the server: %w", err)
		}
		serverURL := "localhost:" + *serverPort
		webServer := webserver.NewWebServer(server, serverURL)
		go webServer.Start()
	}

	//starts the game
	//wait until quit
	<-quit
	//wait for engine graceful shutdown
	if engine != nil {
		<-engine.GetShutdown()
	}
	return nil
}

func createClient(screen tcell.Screen, quit chan struct{}, worldUpdateRate int, consoleEventManager consolemanager.ConsoleEventManager) (client.Engine, error) {
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
	return clientImpl.NewEngine(screen, consoleEventManager, engineConfiguration)
}
