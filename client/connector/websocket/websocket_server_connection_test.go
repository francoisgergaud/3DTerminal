package websocketconnector

import (
	"errors"
	websocket "francoisgergaud/3dGame/common/connector"
	"francoisgergaud/3dGame/common/event"
	testClient "francoisgergaud/3dGame/internal/testutils/client"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	testwebsocket "francoisgergaud/3dGame/internal/testutils/common/connector"
)

type MockFactories struct {
	mock.Mock
}

func (mockFactory *MockFactories) bufferProvider() []event.Event {
	args := mockFactory.Called()
	return args.Get(0).([]event.Event)
}

type MockWebsocketDialer struct {
	mock.Mock
}

func (dialer *MockWebsocketDialer) Dial(urlStr string, requestHeader http.Header) (websocket.WebsocketConnection, *http.Response, error) {
	args := dialer.Mock.Called(urlStr, requestHeader)
	var errCast error
	var err = args.Get(2)
	if err != nil {
		errCast = err.(error)
	}
	var httpResponseCast *http.Response
	var httpResponse = args.Get(2)
	if httpResponse != nil {
		httpResponseCast = err.(*http.Response)
	}
	return args.Get(0).(websocket.WebsocketConnection), httpResponseCast, errCast
}

func TestNewWebSocketServerConnection(t *testing.T) {
	engine := new(testClient.MockEngine)
	var webSocketServerConnectionCapture *WebSocketServerConnection
	quit := make(chan interface{})
	//capture
	engine.On("ConnectToServer", mock.MatchedBy(
		func(wsConnector *WebSocketServerConnection) bool {
			webSocketServerConnectionCapture = wsConnector
			return true
		},
	))
	url := "testURL"
	mockWebsocketDialer := new(MockWebsocketDialer)
	mockWebsocketConnection := new(testwebsocket.MockWebsockeConnection)
	mockWebsocketDialer.On("Dial", url, mock.AnythingOfType("http.Header")).Return(mockWebsocketConnection, nil, nil)
	mockWebsocketConnection.On("ReadJSON", mock.AnythingOfType("*[]event.Event")).Return(nil)
	NewWebSocketServerConnection(engine, url, mockWebsocketDialer, quit)
	assert.Same(t, webSocketServerConnectionCapture.engine, engine)
	assert.Same(t, webSocketServerConnectionCapture.wsConnection, mockWebsocketConnection)
	assert.True(t, webSocketServerConnectionCapture.quit == quit)
	mock.AssertExpectationsForObjects(t, mockWebsocketDialer, engine)
}

func TestRun(t *testing.T) {
	engine := new(testClient.MockEngine)
	mockWebsocketConnection := new(testwebsocket.MockWebsockeConnection)
	playerID := "playerTest"
	mockFactories := new(MockFactories)
	webSocketServerConnection := &WebSocketServerConnection{
		engine:         engine,
		wsConnection:   mockWebsocketConnection,
		playerID:       playerID,
		quit:           make(chan interface{}),
		bufferProvider: mockFactories.bufferProvider,
	}
	eventsFromServer := make([]event.Event, 0)
	eventsFromServer2 := make([]event.Event, 0)
	mockFactories.On("bufferProvider").Return(eventsFromServer).Once()
	mockFactories.On("bufferProvider").Return(eventsFromServer2).Once()
	engine.On("ReceiveEventsFromServer", eventsFromServer)
	errorFromReader := errors.New("test read-error")
	mockWebsocketConnection.On("ReadJSON", &eventsFromServer).Return(nil).Once()
	mockWebsocketConnection.On("ReadJSON", &eventsFromServer2).Return(errorFromReader).Once()

	webSocketServerConnection.Run()

	mock.AssertExpectationsForObjects(t, mockWebsocketConnection, engine, mockFactories)
}

func TestNotifyServer(t *testing.T) {
	mockWebsocketConnection := new(testwebsocket.MockWebsockeConnection)
	eventsToSend := []event.Event{}
	webSocketServerConnection := WebSocketServerConnection{
		wsConnection: mockWebsocketConnection,
	}
	mockWebsocketConnection.On("WriteJSON", eventsToSend).Return(nil).Once()
	err := webSocketServerConnection.NotifyServer(eventsToSend)
	assert.NoError(t, err)
	mock.AssertExpectationsForObjects(t, mockWebsocketConnection)
}

func TestNotifyServerWithError(t *testing.T) {
	mockWebsocketConnection := new(testwebsocket.MockWebsockeConnection)
	eventsToSend := []event.Event{}
	webSocketServerConnection := WebSocketServerConnection{
		wsConnection: mockWebsocketConnection,
	}
	errorFromRWriter := errors.New("test write-error")
	mockWebsocketConnection.On("WriteJSON", eventsToSend).Return(errorFromRWriter).Once()
	err := webSocketServerConnection.NotifyServer(eventsToSend)
	assert.Error(t, err)
	mock.AssertExpectationsForObjects(t, mockWebsocketConnection)
}

func TestDisconnect(t *testing.T) {
	mockWebsocketConnection := new(testwebsocket.MockWebsockeConnection)
	webSocketServerConnection := WebSocketServerConnection{
		wsConnection: mockWebsocketConnection,
	}
	mockWebsocketConnection.On("Close").Return(nil)
	webSocketServerConnection.Disconnect()
	mock.AssertExpectationsForObjects(t, mockWebsocketConnection)
}
