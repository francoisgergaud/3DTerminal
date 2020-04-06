package main

import (
	"francoisgergaud/3dGame/client"
	clienWwebsocketconnector "francoisgergaud/3dGame/client/connector/websocket"
	"francoisgergaud/3dGame/client/consolemanager"
	"francoisgergaud/3dGame/common/runner"
	testclient "francoisgergaud/3dGame/internal/testutils/client"
	testconsolemanager "francoisgergaud/3dGame/internal/testutils/client/consolemanager"
	testrunner "francoisgergaud/3dGame/internal/testutils/common/runner"
	testserver "francoisgergaud/3dGame/internal/testutils/server"
	testtcell "francoisgergaud/3dGame/internal/testutils/tcell"
	"francoisgergaud/3dGame/server"
	webserver "francoisgergaud/3dGame/server/net"
	"testing"
	"time"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGameFactories struct {
	mock.Mock
}

func (mock *mockGameFactories) createServer(quit chan interface{}, worldUpdateRate int) server.Server {
	args := mock.Called(quit, worldUpdateRate)
	return args.Get(0).(server.Server)
}

func (mock *mockGameFactories) createScreen() tcell.Screen {
	args := mock.Called()
	return args.Get(0).(tcell.Screen)
}

func (mock *mockGameFactories) createConsoleEventManager(screen tcell.Screen, quit chan<- interface{}) consolemanager.ConsoleEventManager {
	args := mock.Called(screen, quit)
	return args.Get(0).(consolemanager.ConsoleEventManager)
}

func (mock *mockGameFactories) createClient(quit chan interface{}, worldUpdateRate int, consoleEventManager consolemanager.ConsoleEventManager, screen tcell.Screen) client.Engine {
	args := mock.Called(quit, worldUpdateRate, consoleEventManager, screen)
	return args.Get(0).(client.Engine)
}

func (mock *mockGameFactories) localServerConnection(engine client.Engine, server server.Server, quit <-chan interface{}) {
	mock.Called(engine, server, quit)
}

func (mock *mockGameFactories) createWebServer(address, port string, server server.Server) *webserver.WebServer {
	args := mock.Called(address, port, server)
	return args.Get(0).(*webserver.WebServer)
}

func (mock *mockGameFactories) connectToWebserver(quit chan<- interface{}, client client.Engine, remoteAddress string) *clienWwebsocketconnector.WebSocketServerConnection {
	args := mock.Called(quit, client, remoteAddress)
	return args.Get(0).(*clienWwebsocketconnector.WebSocketServerConnection)
}

func (mock *mockGameFactories) createSignalListener(quit chan<- interface{}) {
	mock.Called(quit)
}

func TestNewGame(t *testing.T) {
	game := NewGame()
	assert.IsType(t, &runner.AsyncRunner{}, game.runner)
	assert.NotNil(t, game.connectToWebserver)
	assert.NotNil(t, game.createClient)
	assert.NotNil(t, game.createServer)
	assert.NotNil(t, game.createWebServer)
	assert.NotNil(t, game.localServerConnection)
}

func TestInitLocal(t *testing.T) {
	mockGameFactories := new(mockGameFactories)
	server := new(testserver.MockServer)
	client := new(testclient.MockEngine)
	quit := make(chan interface{})
	screen := new(testtcell.MockScreen)
	consoleEventManager := new(testconsolemanager.MockConsoleEventManager)
	mockGameFactories.On("createScreen").Return(screen)
	mockGameFactories.On("createConsoleEventManager", screen, mock.MatchedBy(func(channel chan<- interface{}) bool { return channel == quit })).Return(consoleEventManager)
	mockGameFactories.On("createClient", quit, 20, consoleEventManager, screen).Return(client)
	mockGameFactories.On("createServer", mock.MatchedBy(func(channel <-chan interface{}) bool { return channel == quit }), 20).Return(server)
	mockGameFactories.On("localServerConnection", client, server, mock.MatchedBy(func(channel <-chan interface{}) bool { return channel == quit }))
	server.On("Start")
	client.On("Shutdown")
	server.On("Shutdown")
	game := &Game{
		createScreen:              mockGameFactories.createScreen,
		createConsoleEventManager: mockGameFactories.createConsoleEventManager,
		createClient:              mockGameFactories.createClient,
		createServer:              mockGameFactories.createServer,
		localServerConnection:     mockGameFactories.localServerConnection,
		quit:                      quit,
	}
	go func() {
		<-time.After(time.Millisecond)
		close(game.quit)
	}()
	game.InitLocalGame()
	mock.AssertExpectationsForObjects(t, mockGameFactories, client, server)
}

func TestInitRemote(t *testing.T) {
	port := "portNumber"
	mockGameFactories := new(mockGameFactories)
	server := new(testserver.MockServer)
	client := new(testclient.MockEngine)
	runner := new(testrunner.MockRunner)
	quit := make(chan interface{})
	screen := new(testtcell.MockScreen)
	consoleEventManager := new(testconsolemanager.MockConsoleEventManager)
	webServer := &webserver.WebServer{}
	websocketServerConnection := &clienWwebsocketconnector.WebSocketServerConnection{}
	mockGameFactories.On("createScreen").Return(screen)
	mockGameFactories.On("createConsoleEventManager", screen, mock.MatchedBy(func(channel chan<- interface{}) bool { return channel == quit })).Return(consoleEventManager)
	mockGameFactories.On("createClient", quit, 20, consoleEventManager, screen).Return(client).Return(client)
	mockGameFactories.On("createServer", mock.MatchedBy(func(channel <-chan interface{}) bool { return channel == quit }), 20).Return(server)
	mockGameFactories.On("createWebServer", "localhost:", port, server).Return(webServer)
	mockGameFactories.On("connectToWebserver", mock.MatchedBy(func(channel chan<- interface{}) bool { return channel == quit }), client, "localhost:"+port).Return(websocketServerConnection)
	runner.On("Start", webServer)
	runner.On("Start", websocketServerConnection)
	server.On("Start")
	client.On("Shutdown")
	server.On("Shutdown")
	game := &Game{
		runner:                    runner,
		createScreen:              mockGameFactories.createScreen,
		createConsoleEventManager: mockGameFactories.createConsoleEventManager,
		createClient:              mockGameFactories.createClient,
		createServer:              mockGameFactories.createServer,
		connectToWebserver:        mockGameFactories.connectToWebserver,
		createWebServer:           mockGameFactories.createWebServer,
		quit:                      quit,
	}
	go func() {
		<-time.After(time.Millisecond)
		close(game.quit)
	}()
	game.InitRemoteGame(port)
	mock.AssertExpectationsForObjects(t, mockGameFactories, client, server, runner)
}

func TestInitRemoteClient(t *testing.T) {
	remoteAddress := "address"
	mockGameFactories := new(mockGameFactories)
	client := new(testclient.MockEngine)
	runner := new(testrunner.MockRunner)
	quit := make(chan interface{})
	screen := new(testtcell.MockScreen)
	consoleEventManager := new(testconsolemanager.MockConsoleEventManager)
	websocketServerConnection := &clienWwebsocketconnector.WebSocketServerConnection{}
	mockGameFactories.On("createScreen").Return(screen)
	mockGameFactories.On("createConsoleEventManager", screen, mock.MatchedBy(func(channel chan<- interface{}) bool { return channel == quit })).Return(consoleEventManager)
	mockGameFactories.On("createClient", quit, 20, consoleEventManager, screen).Return(client).Return(client)
	mockGameFactories.On("connectToWebserver", mock.MatchedBy(func(channel chan<- interface{}) bool { return channel == quit }), client, remoteAddress).Return(websocketServerConnection)
	runner.On("Start", websocketServerConnection)
	client.On("Shutdown")
	game := &Game{
		runner:                    runner,
		createScreen:              mockGameFactories.createScreen,
		createConsoleEventManager: mockGameFactories.createConsoleEventManager,
		createClient:              mockGameFactories.createClient,
		connectToWebserver:        mockGameFactories.connectToWebserver,
		quit:                      quit,
	}
	go func() {
		<-time.After(time.Millisecond)
		close(game.quit)
	}()
	game.InitRemoteClient(remoteAddress)
	mock.AssertExpectationsForObjects(t, mockGameFactories, client, runner)
}

func TestInitRemoteServer(t *testing.T) {
	port := "portNumber"
	mockGameFactories := new(mockGameFactories)
	server := new(testserver.MockServer)
	runner := new(testrunner.MockRunner)
	quit := make(chan interface{})
	webServer := &webserver.WebServer{}
	mockGameFactories.On("createServer", mock.MatchedBy(func(channel <-chan interface{}) bool { return channel == quit }), 20).Return(server)
	mockGameFactories.On("createWebServer", "localhost:", port, server).Return(webServer)
	mockGameFactories.On("createSignalListener", mock.MatchedBy(func(channel chan<- interface{}) bool { return channel == quit }))
	runner.On("Start", webServer)
	server.On("Start")
	server.On("Shutdown")
	game := &Game{
		runner:               runner,
		createServer:         mockGameFactories.createServer,
		createWebServer:      mockGameFactories.createWebServer,
		createSignalListener: mockGameFactories.createSignalListener,
		quit:                 quit,
	}
	go func() {
		<-time.After(time.Millisecond)
		close(game.quit)
	}()
	game.InitRemoteServer(port)
	mock.AssertExpectationsForObjects(t, mockGameFactories, server, runner)
}
