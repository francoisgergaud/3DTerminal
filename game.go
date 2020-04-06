package main

import (
	"fmt"
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/client/configuration"
	localServerConnector "francoisgergaud/3dGame/client/connector/local/impl"
	clienWwebsocketconnector "francoisgergaud/3dGame/client/connector/websocket"
	clientwebsocketconnector "francoisgergaud/3dGame/client/connector/websocket"
	"francoisgergaud/3dGame/client/consolemanager"
	consoleManagerImpl "francoisgergaud/3dGame/client/consolemanager/impl"
	clientImpl "francoisgergaud/3dGame/client/impl"
	"francoisgergaud/3dGame/common/runner"
	"francoisgergaud/3dGame/server"
	websocketconnector "francoisgergaud/3dGame/server/connector/websocket"
	serverImpl "francoisgergaud/3dGame/server/impl"
	webserver "francoisgergaud/3dGame/server/net"
	"os"
	"os/signal"
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"

	"github.com/gdamore/tcell"
)

//NewGame is a Game factory
func NewGame() *Game {
	return &Game{
		runner:                    new(runner.AsyncRunner),
		createScreen:              createScreen,
		createConsoleEventManager: consoleManagerImpl.NewConsoleEventManager,
		createServer:              createServer,
		createClient:              createClient,
		localServerConnection:     localServerConnector.NewLocalServerConnection,
		createWebServer:           createWebServer,
		connectToWebserver:        connectToWebserver,
		createSignalListener:      createSignalListener,
		quit:                      make(chan interface{}),
	}
}

//Game represent a game instance which can be started
type Game struct {
	runner                    runner.Runner
	createScreen              func() tcell.Screen
	createConsoleEventManager func(screen tcell.Screen, quit chan<- interface{}) consolemanager.ConsoleEventManager
	createServer              func(quit chan interface{}, worldUpdateRate int) server.Server
	createClient              func(quit chan interface{}, worldUpdateRate int, consoleEventManager consolemanager.ConsoleEventManager, screen tcell.Screen) client.Engine
	localServerConnection     func(engine client.Engine, server server.Server, quit <-chan interface{})
	createWebServer           func(address, port string, server server.Server) *webserver.WebServer
	connectToWebserver        func(quit chan<- interface{}, client client.Engine, remoteAddress string) *clienWwebsocketconnector.WebSocketServerConnection
	createSignalListener      func(quit chan<- interface{})
	quit                      chan interface{}
}

//InitLocalGame initializes a local server and a client connecting locally to it
func (game *Game) InitLocalGame() error {
	screen := game.createScreen()
	consoleEventManager := game.createConsoleEventManager(screen, game.quit)
	worldUpdateRate := 20 //world-update frequency, for both client and server
	var engine client.Engine
	var server server.Server
	server = game.createServer(game.quit, worldUpdateRate)
	server.Start()
	engine = game.createClient(game.quit, worldUpdateRate, consoleEventManager, screen)
	game.localServerConnection(engine, server, game.quit)
	//wait for components graceful shutdown
	engine.Shutdown()
	server.Shutdown()
	return nil
}

//InitRemoteGame initializes a remote server and a client connecting remotly to it.
func (game *Game) InitRemoteGame(serverPort string) error {
	screen := game.createScreen()
	consoleEventManager := game.createConsoleEventManager(screen, game.quit)
	worldUpdateRate := 20 //world-update frequency, for both client and server
	var engine client.Engine
	var server server.Server
	server = game.createServer(game.quit, worldUpdateRate)
	server.Start()
	engine = game.createClient(game.quit, worldUpdateRate, consoleEventManager, screen)
	webServer := game.createWebServer("localhost:", serverPort, server)
	game.runner.Start(webServer)
	time.Sleep(time.Millisecond)
	webserverConnection := game.connectToWebserver(game.quit, engine, "localhost:"+serverPort)
	game.runner.Start(webserverConnection)
	//wait for engine graceful shutdown
	engine.Shutdown()
	server.Shutdown()
	return nil
}

//InitRemoteClient initializes a client connecting to a remote server
func (game *Game) InitRemoteClient(remoteAddress string) error {
	screen := game.createScreen()
	consoleEventManager := game.createConsoleEventManager(screen, game.quit)
	worldUpdateRate := 20 //world-update frequency, for both client and server
	var engine client.Engine
	engine = game.createClient(game.quit, worldUpdateRate, consoleEventManager, screen)
	webserverConnection := game.connectToWebserver(game.quit, engine, remoteAddress)
	game.runner.Start(webserverConnection)
	//wait for engine graceful shutdown
	engine.Shutdown()
	return nil
}

//InitRemoteServer initializes a server accesssible remotly.
func (game *Game) InitRemoteServer(serverPort string) error {
	//Remote server does not have a console-manager associated. The server will be close using the following close-handler
	game.createSignalListener(game.quit)
	worldUpdateRate := 20 //world-update frequency, for both client and server
	var server server.Server
	server = game.createServer(game.quit, worldUpdateRate)
	server.Start()
	webServer := game.createWebServer("localhost:", serverPort, server)
	game.runner.Start(webServer)
	//starts the game and wait until quit
	server.Shutdown()
	return nil
}

func createScreen() tcell.Screen {
	tcell.SetEncodingFallback(tcell.EncodingFallbackUTF8)
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(fmt.Errorf("error while instantiating the screen: %w", err))
	}
	if err = screen.Init(); err != nil {
		panic(fmt.Errorf("error while initializing the screen: %w", err))
	}
	return screen
}

func createClient(quit chan interface{}, worldUpdateRate int, consoleEventManager consolemanager.ConsoleEventManager, screen tcell.Screen) client.Engine {
	engineConfiguration := configuration.NewConfiguration(worldUpdateRate)
	client, err := clientImpl.NewEngine(screen, consoleEventManager, engineConfiguration, quit)
	if err != nil {
		panic(fmt.Errorf("error while instantiating the client: %w", err))
	}
	return client
}

func createServer(quit chan interface{}, worldUpdateRate int) server.Server {
	server, err := serverImpl.NewServer(worldUpdateRate, quit)
	if err != nil {
		panic(fmt.Errorf("error while instantiating the server: %w", err))
	}
	return server
}

func createWebServer(address, port string, server server.Server) *webserver.WebServer {
	websocketUpgrader := websocketconnector.NewWebsocketUpgraderWwrapper(&gorillaWebsocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})

	return webserver.NewWebServer(server, address+port, websocketUpgrader)
}

func connectToWebserver(quit chan<- interface{}, client client.Engine, remoteAddress string) *clienWwebsocketconnector.WebSocketServerConnection {
	dialer := clienWwebsocketconnector.NewWebsocketDialerWrapper()
	webserverConnection, err := clientwebsocketconnector.NewWebSocketServerConnection(client, "ws://"+remoteAddress+"/join", dialer, quit)
	if err != nil {
		panic(fmt.Errorf("Error while initializing connection to server: %w", err))
	}
	return webserverConnection
}

func createSignalListener(quit chan<- interface{}) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		close(quit)
	}()
}
