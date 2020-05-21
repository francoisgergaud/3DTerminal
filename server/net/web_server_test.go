package webserver

import (
	"errors"
	testwebsocket "francoisgergaud/3dGame/internal/testutils/common/connector"
	testrunner "francoisgergaud/3dGame/internal/testutils/common/runner"
	testserver "francoisgergaud/3dGame/internal/testutils/server"
	"francoisgergaud/3dGame/server"
	websocketconnector "francoisgergaud/3dGame/server/connector/websocket"
	"net/http"
	"testing"

	websocket "francoisgergaud/3dGame/common/connector"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/runner"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWebsocketUpgrader struct {
	mock.Mock
}

func (websocketUpgrader *mockWebsocketUpgrader) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (websocket.WebsocketConnection, error) {
	args := websocketUpgrader.Called(w, r, responseHeader)
	return args.Get(0).(websocket.WebsocketConnection), args.Error(1)
}

type mockHttpServer struct {
	mock.Mock
}

func (httpServer *mockHttpServer) Handle(pattern string, handler http.Handler) {
	httpServer.Called(pattern, handler)
}
func (httpServer *mockHttpServer) ListenAndServe(addr string, handler http.Handler) error {
	args := httpServer.Called(addr, handler)
	return args.Error(0)
}

type mockResponseWriter struct {
	mock.Mock
}

func (responseWriter *mockResponseWriter) Header() http.Header {
	args := responseWriter.Called()
	return args.Get(0).(http.Header)
}

func (responseWriter *mockResponseWriter) Write(content []byte) (int, error) {
	args := responseWriter.Called(content)
	return args.Int(0), args.Error(1)
}

func (responseWriter *mockResponseWriter) WriteHeader(statusCode int) {
	responseWriter.Called(statusCode)
}

type mockPlayerJoinHandlerFactories struct {
	mock.Mock
}

func (mock *mockPlayerJoinHandlerFactories) websocketClientConnectionFactory(eventToSendToCLient chan event.Event) *websocketconnector.WebSocketClientConnection {
	args := mock.Called(eventToSendToCLient)
	return args.Get(0).(*websocketconnector.WebSocketClientConnection)
}

func (mock *mockPlayerJoinHandlerFactories) websocketClientListenerFactory(playerID string, wsConnection websocket.WebsocketConnection, server server.Server) *websocketconnector.ClientWebSocketListener {
	args := mock.Called(playerID, wsConnection, server)
	return args.Get(0).(*websocketconnector.ClientWebSocketListener)
}

func (mock *mockPlayerJoinHandlerFactories) websocketClientSenderFactory(wsConnection websocket.WebsocketConnection, eventToSendToCLient chan event.Event) *websocketconnector.ClientWebSocketSender {
	args := mock.Called(wsConnection, eventToSendToCLient)
	return args.Get(0).(*websocketconnector.ClientWebSocketSender)
}

func TestNewWebServer(t *testing.T) {
	server := new(testserver.MockServer)
	address := "testurl"
	websocketUpgrader := new(mockWebsocketUpgrader)
	webServer := NewWebServer(server, address, websocketUpgrader)
	assert.Equal(t, address, webServer.serverAddress)
	assert.IsType(t, &HttpServerWrapper{}, webServer.httpServer)
	assert.Equal(t, server, webServer.playerJoinHandler.server)
	assert.Equal(t, websocketUpgrader, webServer.playerJoinHandler.upgrader)
	assert.IsType(t, &runner.AsyncRunner{}, webServer.playerJoinHandler.runner)
	assert.IsType(t, websocketconnector.NewWebSocketClientConnection, webServer.playerJoinHandler.websocketClientConnectionFactory)
	assert.IsType(t, websocketconnector.NewClientWebSocketListener, webServer.playerJoinHandler.websocketClientListenerFactory)
	assert.IsType(t, websocketconnector.NewClientWebSocketSender, webServer.playerJoinHandler.websocketClientSenderFactory)
}

func TestWebServerRun(t *testing.T) {
	httpServer := new(mockHttpServer)
	playerJoinHandler := new(PlayerJoinHandler)
	serverAddress := "serverUrl"
	err := errors.New("test-error")
	webServer := &WebServer{
		serverAddress:     serverAddress,
		httpServer:        httpServer,
		playerJoinHandler: playerJoinHandler,
	}
	httpServer.On("Handle", "/join", playerJoinHandler)
	httpServer.On("ListenAndServe", serverAddress, nil).Return(err)
	webServer.Run()
	mock.AssertExpectationsForObjects(t, httpServer)
}

func TestPlayerJoinHandlerServeHTTP(t *testing.T) {
	websocketUpgrader := new(mockWebsocketUpgrader)
	playerJoinHandlerFactories := new(mockPlayerJoinHandlerFactories)
	server := new(testserver.MockServer)
	playerID := "playerID"
	runner := new(testrunner.MockRunner)
	playerJoinHandler := PlayerJoinHandler{
		runner:                           runner,
		upgrader:                         websocketUpgrader,
		server:                           server,
		websocketClientConnectionFactory: playerJoinHandlerFactories.websocketClientConnectionFactory,
		websocketClientSenderFactory:     playerJoinHandlerFactories.websocketClientSenderFactory,
		websocketClientListenerFactory:   playerJoinHandlerFactories.websocketClientListenerFactory,
	}
	reponseWriter := new(mockResponseWriter)
	reader := &http.Request{}
	websocketConnection := new(testwebsocket.MockWebsockeConnection)
	websocketUpgrader.On("Upgrade", reponseWriter, reader, http.Header(nil)).Return(websocketConnection, nil)
	clientConnection := new(websocketconnector.WebSocketClientConnection)
	var eventsChannelForClientFactory chan event.Event
	playerJoinHandlerFactories.On("websocketClientConnectionFactory", mock.MatchedBy(
		func(eventsChannel chan event.Event) bool {
			eventsChannelForClientFactory = eventsChannel
			return true
		},
	)).Return(clientConnection)
	clientWebsocketSender := new(websocketconnector.ClientWebSocketSender)
	var eventsChannelForClientSenderFactory chan event.Event
	playerJoinHandlerFactories.On("websocketClientSenderFactory", websocketConnection, mock.MatchedBy(
		func(eventsChannel chan event.Event) bool {
			eventsChannelForClientSenderFactory = eventsChannel
			return true
		},
	)).Return(clientWebsocketSender)
	server.On("RegisterPlayer", clientConnection).Return(playerID)
	clientWebsocketListener := new(websocketconnector.ClientWebSocketListener)
	playerJoinHandlerFactories.On("websocketClientListenerFactory", playerID, websocketConnection, server).Return(clientWebsocketListener)
	runner.On("Start", clientWebsocketSender)
	runner.On("Start", clientWebsocketListener)
	playerJoinHandler.ServeHTTP(reponseWriter, reader)
	assert.Equal(t, eventsChannelForClientFactory, eventsChannelForClientSenderFactory)
	mock.AssertExpectationsForObjects(t, websocketUpgrader, playerJoinHandlerFactories, server, runner)
}
