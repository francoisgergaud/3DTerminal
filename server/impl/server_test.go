package impl

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/math/helper"
	mathhelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/runner"
	testhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	testrunner "francoisgergaud/3dGame/internal/testutils/common/runner"
	testbot "francoisgergaud/3dGame/internal/testutils/server/bot"
	testconnector "francoisgergaud/3dGame/internal/testutils/server/connector"
	"francoisgergaud/3dGame/server/bot"
	"francoisgergaud/3dGame/server/connector"
	"testing"
	"time"

	testanimatedelement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFactories struct {
	mock.Mock
}

func (mock *MockFactories) NewWorldMap() world.WorldMap {
	args := mock.Called()
	return args.Get(0).(world.WorldMap)
}

func (mock *MockFactories) NewBot(id string, worldMap world.WorldMap, mathHelper mathhelper.MathHelper, quit chan struct{}) bot.Bot {
	args := mock.Called(id, worldMap, mathHelper, quit)
	return args.Get(0).(bot.Bot)
}

func (mock *MockFactories) NewID() uuid.UUID {
	args := mock.Called()
	return args.Get(0).(uuid.UUID)
}

func (mock *MockFactories) NewPlayer(worldMap world.WorldMap, mathHelper helper.MathHelper, quit chan struct{}) animatedelement.AnimatedElement {
	args := mock.Called(worldMap, mathHelper, quit)
	return args.Get(0).(animatedelement.AnimatedElement)
}

type mockClientEventSender struct {
	mock.Mock
}

func (mock *mockClientEventSender) Run() {
	mock.Called()
}

func (mock *mockClientEventSender) AddClient(playerID string, connectionToClient connector.ClientConnection) {
	mock.Called(playerID, connectionToClient)
}
func (mock *mockClientEventSender) RemoveClient(playerID string) {
	mock.Called(playerID)
}

func (mock *mockClientEventSender) SendEventToClient(playerID string, eventToSend event.Event) {
	mock.Called(playerID, eventToSend)
}
func (mock *mockClientEventSender) ReceiveEvent(event event.Event) {
	mock.Called(event)
}

func TestNewServer(t *testing.T) {
	quit := make(chan struct{})
	worldUpdateRate := 3
	server, error := NewServer(worldUpdateRate, quit)
	assert.Nil(t, error)
	assert.Nil(t, server.worldMap)
	assert.IsType(t, &helper.MathHelperImpl{}, server.mathHelper)
	assert.Len(t, server.players, 0)
	assert.NotNil(t, server.clientEventSender)
	assert.Len(t, server.bots, 0)
	assert.Equal(t, worldUpdateRate, server.botsUpdateRate)
	assert.IsType(t, &runner.AsyncRunner{}, server.runner)
	//TODO: create interface for factory: NotNill do not check if this is the expected factory
	assert.NotNil(t, server.worldMapFactory)
	assert.NotNil(t, server.botFactory)
	assert.NotNil(t, server.identifierFactory)
}

func TestStart(t *testing.T) {
	quit := make(chan struct{})
	eventQueue := make(chan event.Event, 100)
	mathHelper := new(testhelper.MockMathHelper)
	mockFactories := new(MockFactories)
	worldMap := new(testworld.MockWorldMap)
	uuid := uuid.New()
	clientEventSender := &clientEventSenderImp{
		eventQueue: eventQueue,
	}
	mockFactories.On("NewWorldMap").Return(worldMap)
	mockFactories.On("NewID").Return(uuid)
	mockBot := new(testbot.MockBot)
	mockFactories.On("NewBot", uuid.String(), worldMap, mathHelper, quit).Return(mockBot)
	mockBot.MockEventPublisher.On("RegisterListener", clientEventSender)
	runner := new(testrunner.MockRunner)
	server := &Impl{
		identifierFactory: mockFactories.NewID,
		worldMapFactory:   mockFactories.NewWorldMap,
		botFactory:        mockFactories.NewBot,
		mathHelper:        mathHelper,
		quit:              quit,
		clientEventSender: clientEventSender,
		runner:            runner,
		bots:              make(map[string]bot.Bot),
	}
	runner.On("Start", clientEventSender).Once()
	runner.On("Start", server).Once()
	server.Start()
	assert.Equal(t, mockBot, server.bots[uuid.String()])
	mock.AssertExpectationsForObjects(t, mockFactories, runner, &mockBot.MockAnimatedElement, &mockBot.MockEventPublisher)
}

func TestRegisterPlayer(t *testing.T) {
	quit := make(chan struct{})
	worldMap := new(testworld.MockWorldMap)
	serverPlayers := make(map[string]animatedelement.AnimatedElement)
	bots := make(map[string]bot.Bot)
	uuid := uuid.New()
	mockFactories := new(MockFactories)
	mockFactories.On("NewID").Return(uuid)
	clientEventSender := new(mockClientEventSender)
	mathHelper := new(testhelper.MockMathHelper)
	animatedElement := new(testanimatedelement.MockAnimatedElement)
	mockFactories.On("NewPlayer", worldMap, mathHelper, quit).Return(animatedElement)
	animatedElementState := &state.AnimatedElementState{}
	animatedElement.On("GetState").Return(animatedElementState)
	animatedElement.On("Start")

	otherPlayerID := "otherPlayer1"
	otherPlayer := new(testanimatedelement.MockAnimatedElement)
	otherPlayerState := &state.AnimatedElementState{
		Velocity: 0.0025,
	}
	otherPlayer.On("GetState").Return(otherPlayerState)
	serverPlayers[otherPlayerID] = otherPlayer

	bots = make(map[string]bot.Bot)
	botID := "bot1"
	bot := &testbot.MockBot{}
	botState := &state.AnimatedElementState{
		StepAngle: 0.0025,
	}
	bot.MockAnimatedElement.On("GetState").Return(botState)
	bots[botID] = bot

	server := Impl{
		worldMap:          worldMap,
		players:           serverPlayers,
		bots:              bots,
		quit:              quit,
		identifierFactory: mockFactories.NewID,
		clientEventSender: clientEventSender,
		playerFactory:     mockFactories.NewPlayer,
		mathHelper:        mathHelper,
	}
	clientConnection := new(testconnector.MockClientConnection)
	clientEventSender.On("AddClient", uuid.String(), clientConnection)
	worldMap.On("Clone").Return(worldMap)
	var eventForOtherPlayerCapture, eventForPlayerCapture event.Event
	clientEventSender.On(
		"ReceiveEvent",
		mock.MatchedBy(
			func(event event.Event) bool {
				eventForOtherPlayerCapture = event
				return true
			},
		),
	)
	clientEventSender.On(
		"SendEventToClient",
		uuid.String(),
		mock.MatchedBy(
			func(event event.Event) bool {
				eventForPlayerCapture = event
				return true
			},
		),
	)
	server.RegisterPlayer(clientConnection)

	assert.NotEmpty(t, eventForOtherPlayerCapture.PlayerID)
	assert.Nil(t, eventForOtherPlayerCapture.ExtraData["worldMap"])
	assert.Same(t, animatedElementState, eventForOtherPlayerCapture.State)
	assert.Nil(t, eventForOtherPlayerCapture.ExtraData["otherPlayers"])
	assert.Equal(t, "join", eventForOtherPlayerCapture.Action)

	assert.Equal(t, "init", eventForPlayerCapture.Action)
	assert.Equal(t, uuid.String(), eventForPlayerCapture.PlayerID)
	assert.Same(t, animatedElementState, eventForPlayerCapture.State)
	assert.Same(t, worldMap, eventForPlayerCapture.ExtraData["worldMap"])
	assert.Equal(t, *otherPlayerState, eventForPlayerCapture.ExtraData["otherPlayers"].(map[string]state.AnimatedElementState)[otherPlayerID])
	assert.Equal(t, *botState, eventForPlayerCapture.ExtraData["otherPlayers"].(map[string]state.AnimatedElementState)[botID])

	assert.Equal(t, animatedElement, serverPlayers[uuid.String()])

	mock.AssertExpectationsForObjects(t, mockFactories, clientEventSender, animatedElement, worldMap)
}

func TestUnregisterClient(t *testing.T) {
	clientEventSender := new(mockClientEventSender)
	palyers := make(map[string]animatedelement.AnimatedElement)
	playerID := "playerTest"
	palyers[playerID] = new(testanimatedelement.MockAnimatedElement)
	server := Impl{
		clientEventSender: clientEventSender,
		players:           palyers,
	}
	clientEventSender.On("RemoveClient", playerID)
	var eventCapture event.Event
	clientEventSender.On(
		"ReceiveEvent",
		mock.MatchedBy(
			func(event event.Event) bool {
				eventCapture = event
				return true
			},
		),
	)
	server.UnregisterClient(playerID)

	assert.NotContains(t, playerID, server.players)
	assert.Equal(t, "quit", eventCapture.Action)
	assert.Equal(t, playerID, eventCapture.PlayerID)
	mock.AssertExpectationsForObjects(t, clientEventSender)
}

func TestReceiveEventFromClient(t *testing.T) {
	clientEventSender := new(mockClientEventSender)
	palyers := make(map[string]animatedelement.AnimatedElement)
	playerID := "playerTest"
	player := new(testanimatedelement.MockAnimatedElement)
	palyers[playerID] = player
	server := Impl{
		clientEventSender: clientEventSender,
		players:           palyers,
	}
	var eventCapture event.Event
	clientEventSender.On(
		"ReceiveEvent",
		mock.MatchedBy(
			func(event event.Event) bool {
				eventCapture = event
				return true
			},
		),
	)
	eventState := &state.AnimatedElementState{
		MoveDirection: state.Backward,
	}
	eventReceived := event.Event{
		PlayerID: playerID,
		State:    eventState,
	}
	player.On("SetState", eventState)
	server.ReceiveEventFromClient(eventReceived)
	assert.Equal(t, eventReceived, eventCapture)
	mock.AssertExpectationsForObjects(t, player, clientEventSender)
}

func TestSRun(t *testing.T) {
	quit := make(chan struct{})
	palyers := make(map[string]animatedelement.AnimatedElement)
	playerID := "playerTest"
	player := new(testanimatedelement.MockAnimatedElement)
	palyers[playerID] = player
	bots := make(map[string]bot.Bot)
	botID := "botTest"
	bot := new(testbot.MockBot)
	bots[botID] = bot
	server := Impl{
		botsUpdateRate: 1000,
		players:        palyers,
		bots:           bots,
		quit:           quit,
	}
	bot.MockAnimatedElement.On("Move")
	player.On("Move")
	go server.Run()
	<-time.After(time.Millisecond * 5)
	close(quit)
	mock.AssertExpectationsForObjects(t, player, &bot.MockAnimatedElement)
}

func TestClientEventSenderRun(t *testing.T) {
	clientConnection := new(testconnector.MockClientConnection)
	eventQueue := make(chan event.Event, 2)
	quit := make(chan struct{})
	clientConnections := make(map[string]connector.ClientConnection)
	playerID := "playerID"
	clientConnections[playerID] = clientConnection
	clientEventSenderInitalTimeFrame := uint32(2)
	clientEventSender := &clientEventSenderImp{
		clientConnections: clientConnections,
		timeFrame:         clientEventSenderInitalTimeFrame,
		clientUpdateRate:  1000,
		eventQueue:        eventQueue,
		quit:              quit,
	}
	var eventsToCapture []event.Event
	clientConnection.On(
		"SendEventsToClient",
		mock.MatchedBy(
			func(event []event.Event) bool {
				eventsToCapture = event
				return true
			},
		),
	)
	eventToSend := event.Event{
		Action:    "actionTest",
		TimeFrame: 0,
	}
	eventQueue <- eventToSend
	go clientEventSender.Run()
	<-time.After(time.Millisecond * 5)
	close(quit)
	assert.Equal(t, 1, len(eventsToCapture))
	//verify the timeframe is updated before sending the event
	eventToSend.TimeFrame = clientEventSenderInitalTimeFrame
	assert.Equal(t, eventToSend, eventsToCapture[0])
	assert.Greater(t, clientEventSender.timeFrame, clientEventSenderInitalTimeFrame)
	mock.AssertExpectationsForObjects(t, clientConnection)
}

func TestClientEventSenderAddClient(t *testing.T) {
	clientConnection := new(testconnector.MockClientConnection)
	clientConnections := make(map[string]connector.ClientConnection)
	playerID := "playerID"
	clientEventSender := &clientEventSenderImp{
		clientConnections: clientConnections,
	}
	clientEventSender.AddClient(playerID, clientConnection)
	assert.Same(t, clientConnection, clientEventSender.clientConnections[playerID])
}

func TestClientEventSenderRemoveClient(t *testing.T) {
	clientConnection := new(testconnector.MockClientConnection)
	clientConnections := make(map[string]connector.ClientConnection)
	playerID := "playerID"
	clientConnections[playerID] = clientConnection
	clientEventSender := &clientEventSenderImp{
		clientConnections: clientConnections,
	}
	clientEventSender.RemoveClient(playerID)
	assert.Nil(t, clientConnections[playerID])
}

func TestClientEventSenderSendEventToClient(t *testing.T) {
	clientConnection := new(testconnector.MockClientConnection)
	clientConnections := make(map[string]connector.ClientConnection)
	playerID := "playerID"
	clientConnections[playerID] = clientConnection
	clientEventSenderInitalTimeFrame := uint32(2)
	clientEventSender := &clientEventSenderImp{
		clientConnections: clientConnections,
		timeFrame:         clientEventSenderInitalTimeFrame,
	}
	var eventsToCapture []event.Event
	clientConnection.On(
		"SendEventsToClient",
		mock.MatchedBy(
			func(event []event.Event) bool {
				eventsToCapture = event
				return true
			},
		),
	)
	eventToSend := event.Event{
		Action:    "actionTest",
		TimeFrame: 0,
	}
	clientEventSender.SendEventToClient(playerID, eventToSend)
	assert.Equal(t, 1, len(eventsToCapture))
	//verify the timeframe is updated before sending the event
	eventToSend.TimeFrame = clientEventSenderInitalTimeFrame
	assert.Equal(t, eventToSend, eventsToCapture[0])
	mock.AssertExpectationsForObjects(t, clientConnection)
}

func TestClientEventSenderReceiveEvent(t *testing.T) {
	eventQueue := make(chan event.Event, 2)
	clientEventSender := &clientEventSenderImp{
		eventQueue: eventQueue,
	}
	eventToSend := event.Event{
		Action: "actionTest",
	}
	clientEventSender.ReceiveEvent(eventToSend)
	eventReceived := <-eventQueue
	assert.Equal(t, eventToSend, eventReceived)
}
