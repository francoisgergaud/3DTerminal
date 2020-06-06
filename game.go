package main

import (
	"flag"
	"fmt"
	"francoisgergaud/3dGame/client"
	localServerConnector "francoisgergaud/3dGame/client/connector/local/impl"
	clienWwebsocketconnector "francoisgergaud/3dGame/client/connector/websocket"
	clientwebsocketconnector "francoisgergaud/3dGame/client/connector/websocket"
	"francoisgergaud/3dGame/client/consolemanager"
	consoleManagerImpl "francoisgergaud/3dGame/client/consolemanager/impl"
	clientImpl "francoisgergaud/3dGame/client/impl"
	"francoisgergaud/3dGame/server"
	websocketconnector "francoisgergaud/3dGame/server/connector/websocket"
	serverImpl "francoisgergaud/3dGame/server/impl"
	webserver "francoisgergaud/3dGame/server/net"
	"os"
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"

	"github.com/gdamore/tcell"
)

func createScreen() tcell.Screen {
	tcell.SetEncodingFallback(tcell.EncodingFallbackUTF8)
	screen, e := tcell.NewScreen()
	if e != nil {
		panic(e)
	}
	if e = screen.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		panic(e)
	}
	return screen
}

//InitGame initializes a game.
func InitGame() error {
	var mode = flag.String("mode", "local", "possible mode: 'local', 'remote', 'remoteClient', 'remoteServer'")
	var remoteAddress = flag.String("address", "127.0.0.1:9836", "remote-server host-port")
	var serverPort = flag.String("port", "9836", "remote-server host-port")
	flag.Parse()
	quit := make(chan struct{})
	worldUpdateRate := 20 //world-update frequency, for both client and server
	var err error
	var engine client.Engine
	var server server.Server
	if *mode == "local" {
		screen := createScreen()
		consoleEventManager := consoleManagerImpl.NewConsoleEventManager(screen, quit)
		server, err = serverImpl.NewServer(worldUpdateRate, quit)
		if err != nil {
			return fmt.Errorf("error while instantiating the server: %w", err)
		}
		server.Start()
		engine, err = createClient(screen, quit, worldUpdateRate, consoleEventManager)
		if err != nil {
			return fmt.Errorf("Error while initializing engine: %w", err)
		}
		localServerConnector.NewLocalServerConnection(engine, server, quit)
		//serverConnection.ConnectToServer()
	} else if *mode == "remote" {
		screen := createScreen()
		consoleEventManager := consoleManagerImpl.NewConsoleEventManager(screen, quit)
		server, err = serverImpl.NewServer(worldUpdateRate, quit)
		if err != nil {
			return fmt.Errorf("error while instantiating the server: %w", err)
		}
		server.Start()
		engine, err = createClient(screen, quit, worldUpdateRate, consoleEventManager)
		if err != nil {
			return fmt.Errorf("Error while initializing engine: %w", err)
		}
		serverURL := "localhost:" + *serverPort

		websocketUpgrader := websocketconnector.NewWebsocketUpgraderWwrapper(&gorillaWebsocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		})

		webServer := webserver.NewWebServer(server, serverURL, websocketUpgrader)
		go webServer.Run()
		time.Sleep(time.Millisecond)
		dialer := clienWwebsocketconnector.NewWebsocketDialerWrapper()
		webserverConnection, err := clientwebsocketconnector.NewWebSocketServerConnection(engine, "ws://"+serverURL+"/join", dialer, quit)
		if err != nil {
			return fmt.Errorf("Error while initializing connection to server: %w", err)
		}
		go webserverConnection.Run()
	} else if *mode == "remoteClient" {
		screen := createScreen()
		consoleEventManager := consoleManagerImpl.NewConsoleEventManager(screen, quit)
		engine, err = createClient(screen, quit, worldUpdateRate, consoleEventManager)
		if err != nil {
			return fmt.Errorf("Error while initializing engine: %w", err)
		}
		dialer := clienWwebsocketconnector.NewWebsocketDialerWrapper()
		webserverConnection, err := clientwebsocketconnector.NewWebSocketServerConnection(engine, "ws://"+*remoteAddress+"/join", dialer, quit)
		if err != nil {
			return fmt.Errorf("Error while initializing connection to server: %w", err)
		}
		go webserverConnection.Run()
	} else if *mode == "remoteServer" {
		server, err = serverImpl.NewServer(worldUpdateRate, quit)
		if err != nil {
			return fmt.Errorf("error while instantiating the server: %w", err)
		}
		server.Start()
		serverURL := "localhost:" + *serverPort
		websocketUpgrader := websocketconnector.NewWebsocketUpgraderWwrapper(&gorillaWebsocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		})
		webServer := webserver.NewWebServer(server, serverURL, websocketUpgrader)
		go webServer.Run()
	}

	//starts the game
	//wait until quit
	<-quit
	//wait for engine graceful shutdown
	if engine != nil {
		<-engine.Shutdown()
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
