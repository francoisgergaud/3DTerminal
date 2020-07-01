package impl

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/projectile"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"
	mathhelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/runner"
	testeventpublisher "francoisgergaud/3dGame/internal/testutils/common/event/publisher"
	testhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	testrunner "francoisgergaud/3dGame/internal/testutils/common/runner"
	testbot "francoisgergaud/3dGame/internal/testutils/server/bot"
	testconnector "francoisgergaud/3dGame/internal/testutils/server/connector"
	"francoisgergaud/3dGame/server/bot"
	"francoisgergaud/3dGame/server/connector"
	"testing"
	"time"

	testanimatedelement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	testprojectile "francoisgergaud/3dGame/internal/testutils/common/environment/projectile"
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

func (mock *MockFactories) NewBot(id string, worldMap world.WorldMap, mathHelper mathhelper.MathHelper, quit <-chan interface{}) bot.Bot {
	args := mock.Called(id, worldMap, mathHelper, quit)
	return args.Get(0).(bot.Bot)
}

func (mock *MockFactories) NewID() uuid.UUID {
	args := mock.Called()
	return args.Get(0).(uuid.UUID)
}

func (mock *MockFactories) NewPlayer(id string, worldMap world.WorldMap, mathHelper helper.MathHelper, quit <-chan interface{}) animatedelement.AnimatedElement {
	args := mock.Called(id, worldMap, mathHelper, quit)
	return args.Get(0).(animatedelement.AnimatedElement)
}

type mockClientEventSender struct {
	mock.Mock
}

func (mock *mockClientEventSender) Run() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *mockClientEventSender) addClient(playerID string, connectionToClient connector.ClientConnection) {
	mock.Called(playerID, connectionToClient)
}
func (mock *mockClientEventSender) removeClient(playerID string) {
	mock.Called(playerID)
}

func (mock *mockClientEventSender) sendEventToClient(playerID string, eventToSend event.Event) {
	mock.Called(playerID, eventToSend)
}
func (mock *mockClientEventSender) sendEventToAllClients(event event.Event) {
	mock.Called(event)
}

func (mock *mockClientEventSender) close() {
	mock.Called()
}

func (mock *mockClientEventSender) shutdown() {
	mock.Called()
}

type MockSpawner struct {
	mock.Mock
	testeventpublisher.MockEventPublisher
}

func (mock *MockSpawner) Spawn(animatedelementID string, moveDirection state.Direction) {
	mock.Called(animatedelementID, moveDirection)
}

func TestNewServer(t *testing.T) {
	quit := make(chan interface{})
	worldUpdateRate := 3
	server, error := NewServer(worldUpdateRate, quit)
	assert.Nil(t, error)
	assert.Nil(t, server.worldMap)
	assert.IsType(t, &helper.MathHelperImpl{}, server.mathHelper)
	assert.Len(t, server.players, 0)
	assert.NotNil(t, server.clientEventSender)
	assert.Equal(t, worldUpdateRate, server.botsUpdateRate)
	assert.IsType(t, &runner.AsyncRunner{}, server.runner)
	//TODO: create interface for factory: NotNill do not check if this is the expected factory
	assert.NotNil(t, server.worldMapFactory)
	assert.NotNil(t, server.botFactory)
	assert.NotNil(t, server.identifierFactory)
}

func TestStart(t *testing.T) {
	quit := make(chan interface{})
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
	mockFactories.On("NewBot", uuid.String(), worldMap, mathHelper, mock.MatchedBy(func(channel <-chan interface{}) bool { return channel == quit })).Return(mockBot)
	runner := new(testrunner.MockRunner)
	server := &Impl{
		identifierFactory: mockFactories.NewID,
		worldMapFactory:   mockFactories.NewWorldMap,
		botFactory:        mockFactories.NewBot,
		mathHelper:        mathHelper,
		quit:              quit,
		clientEventSender: clientEventSender,
		runner:            runner,
		players:           make(map[string]animatedelement.AnimatedElement),
	}
	runner.On("Start", clientEventSender).Once()
	runner.On("Start", server).Once()
	mockBot.MockEventPublisher.On("RegisterListener", server)
	server.Start()
	assert.Equal(t, mockBot, server.players[uuid.String()])
	mock.AssertExpectationsForObjects(t, mockFactories, runner, &mockBot.MockAnimatedElement, &mockBot.MockEventPublisher)
}

func TestRegisterPlayer(t *testing.T) {
	quit := make(chan interface{})
	worldMap := new(testworld.MockWorldMap)
	serverPlayers := make(map[string]animatedelement.AnimatedElement)
	serverProjectiles := make(map[string]projectile.Projectile)
	uuid := uuid.New()
	mockFactories := new(MockFactories)
	mockFactories.On("NewID").Return(uuid)
	clientEventSender := new(mockClientEventSender)
	mathHelper := new(testhelper.MockMathHelper)
	animatedElement := new(testanimatedelement.MockAnimatedElement)
	mockFactories.On("NewPlayer", uuid.String(), worldMap, mathHelper, mock.MatchedBy(func(channel <-chan interface{}) bool { return channel == quit })).Return(animatedElement)
	animatedElementState := &state.AnimatedElementState{}
	animatedElement.On("State").Return(animatedElementState)
	otherPlayerID := "otherPlayer1"
	otherPlayer := new(testanimatedelement.MockAnimatedElement)
	otherPlayerState := &state.AnimatedElementState{
		Velocity: 0.0025,
	}
	otherPlayer.On("State").Return(otherPlayerState)
	serverPlayers[otherPlayerID] = otherPlayer
	projectileID := "projectile1"
	projectile := new(testprojectile.MockProjectile)
	projectileState := &state.AnimatedElementState{
		Velocity: 0.0075,
		Angle:    0.75,
	}
	projectile.MockAnimatedElement.On("State").Return(projectileState)
	serverProjectiles[projectileID] = projectile
	server := Impl{
		worldMap:          worldMap,
		players:           serverPlayers,
		projectiles:       serverProjectiles,
		quit:              quit,
		identifierFactory: mockFactories.NewID,
		clientEventSender: clientEventSender,
		playerFactory:     mockFactories.NewPlayer,
		mathHelper:        mathHelper,
	}
	clientConnection := new(testconnector.MockClientConnection)
	clientEventSender.On("addClient", uuid.String(), clientConnection)
	var eventForOtherPlayerCapture, eventForPlayerCapture event.Event
	clientEventSender.On(
		"sendEventToAllClients",
		mock.MatchedBy(
			func(event event.Event) bool {
				eventForOtherPlayerCapture = event
				return true
			},
		),
	)
	clientEventSender.On(
		"sendEventToClient",
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
	assert.Equal(t, otherPlayerState, eventForPlayerCapture.ExtraData["otherPlayers"].(map[string]*state.AnimatedElementState)[otherPlayerID])
	assert.Equal(t, projectileState, eventForPlayerCapture.ExtraData["projectiles"].(map[string]*state.AnimatedElementState)[projectileID])
	assert.Equal(t, animatedElement, serverPlayers[uuid.String()])
	mock.AssertExpectationsForObjects(t, mockFactories, clientEventSender, animatedElement, projectile, worldMap)
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
	clientEventSender.On("removeClient", playerID)
	var eventCapture event.Event
	clientEventSender.On(
		"sendEventToAllClients",
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

func TestReceiveMoveEventFromClient(t *testing.T) {
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
		"sendEventToAllClients",
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
		Action:   "move",
		State:    eventState,
	}
	player.On("SetState", eventState)
	server.ReceiveEventFromClient(eventReceived)
	assert.Equal(t, eventReceived, eventCapture)
	mock.AssertExpectationsForObjects(t, player, clientEventSender)
}

func TestReceiveFireEventFromClient(t *testing.T) {
	clientEventSender := new(mockClientEventSender)
	palyers := make(map[string]animatedelement.AnimatedElement)
	playerID := "playerTest"
	projectileFactoryBuilder := new(testprojectile.MockProjectileFactory)
	mathHelper := new(testhelper.MockMathHelper)
	worldMap := new(testworld.MockWorldMap)
	server := Impl{
		clientEventSender: clientEventSender,
		players:           palyers,
		mathHelper:        mathHelper,
		worldMap:          worldMap,
		projectileFactory: projectileFactoryBuilder.CreateProjectile,
		projectiles:       make(map[string]projectile.Projectile),
	}
	var eventCapture event.Event
	clientEventSender.On(
		"sendEventToAllClients",
		mock.MatchedBy(
			func(event event.Event) bool {
				eventCapture = event
				return true
			},
		),
	)
	projectileID := "projectileIDTest"
	projectilePosition := &math.Point2D{
		X: 2.0,
		Y: 4.0,
	}
	projectileAngle := 1.5
	eventState := &state.AnimatedElementState{
		Position: projectilePosition,
		Angle:    projectileAngle,
	}
	eventReceived := event.Event{
		PlayerID: playerID,
		Action:   "fire",
		State:    eventState,
		ExtraData: map[string]interface{}{
			"projectileID": projectileID,
		},
	}
	projectileToReturn := new(testprojectile.MockProjectile)
	projectileToReturn.MockEventPublisher.On("RegisterListener", &server)
	projectileFactoryBuilder.On("CreateProjectile", projectileID, projectilePosition, projectileAngle, worldMap, palyers, mathHelper).Return(projectileToReturn)

	server.ReceiveEventFromClient(eventReceived)

	assert.Equal(t, eventReceived, eventCapture)
	assert.Equal(t, server.projectiles[projectileID], projectileToReturn)
	mock.AssertExpectationsForObjects(t, projectileToReturn, projectileFactoryBuilder, clientEventSender)
}

func TestRun(t *testing.T) {
	quit := make(chan interface{})
	players := make(map[string]animatedelement.AnimatedElement)
	playerID := "playerTest"
	player := new(testanimatedelement.MockAnimatedElement)
	players[playerID] = player
	botID := "botTest"
	bot := new(testbot.MockBot)
	players[botID] = bot
	projectiles := make(map[string]projectile.Projectile)
	projectileID := "projectileID"
	projectile := new(testprojectile.MockProjectile)
	projectiles[projectileID] = projectile
	server := Impl{
		botsUpdateRate: 1000,
		players:        players,
		quit:           quit,
		projectiles:    projectiles,
	}
	bot.MockAnimatedElement.On("Move")
	player.On("Move")
	projectile.MockAnimatedElement.On("Move")
	go server.Run()
	<-time.After(time.Millisecond * 5)
	close(quit)
	mock.AssertExpectationsForObjects(t, player, &bot.MockAnimatedElement, projectile)
}

func TestReceiveEventMove(t *testing.T) {
	playerID := "playerTest"
	player := new(testanimatedelement.MockAnimatedElement)
	players := make(map[string]animatedelement.AnimatedElement)
	players[playerID] = player
	clientEventSender := new(mockClientEventSender)
	server := Impl{
		players:           players,
		clientEventSender: clientEventSender,
	}
	eventAnimatedElementState := &state.AnimatedElementState{}
	moveEvent := event.Event{
		Action:   "move",
		PlayerID: playerID,
		State:    eventAnimatedElementState,
	}
	clientEventSender.On("sendEventToAllClients", moveEvent)
	player.On("SetState", eventAnimatedElementState)

	server.ReceiveEvent(moveEvent)

	mock.AssertExpectationsForObjects(t, player, clientEventSender)
}

func TestReceiveEventSpawn(t *testing.T) {
	playerID := "playerTest"
	player := new(testanimatedelement.MockAnimatedElement)
	players := make(map[string]animatedelement.AnimatedElement)
	players[playerID] = player
	clientEventSender := new(mockClientEventSender)
	server := Impl{
		players:           players,
		clientEventSender: clientEventSender,
	}
	eventAnimatedElementState := &state.AnimatedElementState{}
	spawnEvent := event.Event{
		Action:   "spawn",
		PlayerID: playerID,
		State:    eventAnimatedElementState,
	}
	clientEventSender.On("sendEventToAllClients", spawnEvent)
	server.ReceiveEvent(spawnEvent)

	mock.AssertExpectationsForObjects(t, player, clientEventSender)
}

func TestReceiveEventProjectileWallImpact(t *testing.T) {
	projectileID := "projectileIDTest"
	projectiles := make(map[string]projectile.Projectile)
	projectile := new(testprojectile.MockProjectile)
	projectiles[projectileID] = projectile
	clientEventSender := new(mockClientEventSender)
	server := Impl{
		projectiles:       projectiles,
		clientEventSender: clientEventSender,
	}
	projectileWallImpactEvent := event.Event{
		Action:   "projectileWallImpact",
		PlayerID: projectileID,
	}
	clientEventSender.On("sendEventToAllClients", mock.MatchedBy(
		func(eventToSend event.Event) bool {
			//The originl event is transformed before being sent to the clients
			if eventToSend.Action == "projectileImpact" && eventToSend.PlayerID == projectileID {
				return true
			}
			return false
		},
	))
	server.ReceiveEvent(projectileWallImpactEvent)

	assert.NotContains(t, projectiles, projectile)
	mock.AssertExpectationsForObjects(t, clientEventSender)
}

func TestReceiveEventProjectilePlayerImpact(t *testing.T) {
	projectileID := "projectileIDTest"
	projectiles := make(map[string]projectile.Projectile)
	projectile := new(testprojectile.MockProjectile)
	projectiles[projectileID] = projectile
	playerID := "playerIDTest"
	clientEventSender := new(mockClientEventSender)
	spawner := new(MockSpawner)
	server := Impl{
		projectiles:       projectiles,
		clientEventSender: clientEventSender,
		spawner:           spawner,
	}
	projectilePlayerImpactEvent := event.Event{
		Action:   "projectilePlayerImpact",
		PlayerID: projectileID,
		ExtraData: map[string]interface{}{
			"playerID": playerID,
		},
	}
	clientEventSender.On("sendEventToAllClients", mock.MatchedBy(
		func(eventToSend event.Event) bool {
			//The originl event is transformed before being sent to the clients
			if eventToSend.Action == "projectileImpact" && eventToSend.PlayerID == projectileID {
				return true
			}
			return false
		},
	))
	clientEventSender.On("sendEventToAllClients", mock.MatchedBy(
		func(eventToSend event.Event) bool {
			//The originl event is transformed before being sent to the clients
			if eventToSend.Action == "kill" && eventToSend.PlayerID == playerID {
				return true
			}
			return false
		},
	))
	spawner.On("Spawn", playerID, state.None)

	server.ReceiveEvent(projectilePlayerImpactEvent)

	assert.NotContains(t, projectiles, projectile)
	mock.AssertExpectationsForObjects(t, clientEventSender, spawner)
}

func TestReceiveEventProjectileBotImpact(t *testing.T) {
	projectileID := "projectileIDTest"
	projectiles := make(map[string]projectile.Projectile)
	projectile := new(testprojectile.MockProjectile)
	projectiles[projectileID] = projectile
	playerID := "playerIDTest"
	clientEventSender := new(mockClientEventSender)
	spawner := new(MockSpawner)
	botIDs := []string{playerID}
	server := Impl{
		projectiles:       projectiles,
		clientEventSender: clientEventSender,
		spawner:           spawner,
		botIDs:            botIDs,
	}
	projectilePlayerImpactEvent := event.Event{
		Action:   "projectilePlayerImpact",
		PlayerID: projectileID,
		ExtraData: map[string]interface{}{
			"playerID": playerID,
		},
	}
	clientEventSender.On("sendEventToAllClients", mock.MatchedBy(
		func(eventToSend event.Event) bool {
			//The originl event is transformed before being sent to the clients
			if eventToSend.Action == "projectileImpact" && eventToSend.PlayerID == projectileID {
				return true
			}
			return false
		},
	))
	clientEventSender.On("sendEventToAllClients", mock.MatchedBy(
		func(eventToSend event.Event) bool {
			//The originl event is transformed before being sent to the clients
			if eventToSend.Action == "kill" && eventToSend.PlayerID == playerID {
				return true
			}
			return false
		},
	))
	spawner.On("Spawn", playerID, state.Forward)

	server.ReceiveEvent(projectilePlayerImpactEvent)

	assert.NotContains(t, projectiles, projectile)
	mock.AssertExpectationsForObjects(t, clientEventSender, spawner)
}

func TestClientEventSenderRun(t *testing.T) {
	clientConnection := new(testconnector.MockClientConnection)
	eventQueue := make(chan event.Event, 2)
	quit := make(chan interface{})
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
		shutdownCompleted: make(chan interface{}),
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
	).Return(nil)
	clientConnection.On("Close")
	eventToSend := event.Event{
		Action:    "actionTest",
		TimeFrame: 0,
	}
	eventQueue <- eventToSend
	go clientEventSender.Run()
	<-time.After(time.Millisecond * 5)
	close(quit)
	<-clientEventSender.shutdownCompleted
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
	clientEventSender.addClient(playerID, clientConnection)
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
	clientConnection.On("Close")
	clientEventSender.removeClient(playerID)
	assert.Nil(t, clientConnections[playerID])
	mock.AssertExpectationsForObjects(t, clientConnection)
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
	).Return(nil)
	eventToSend := event.Event{
		Action:    "actionTest",
		TimeFrame: 0,
	}
	clientEventSender.sendEventToClient(playerID, eventToSend)
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
	clientEventSender.sendEventToAllClients(eventToSend)
	eventReceived := <-eventQueue
	assert.Equal(t, eventToSend, eventReceived)
}

func TestClientEventSenderClose(t *testing.T) {
	clientConnection := new(testconnector.MockClientConnection)
	clientConnections := make(map[string]connector.ClientConnection)
	playerID := "playerID"
	clientConnections[playerID] = clientConnection
	clientEventSender := &clientEventSenderImp{
		clientConnections: clientConnections,
		shutdownCompleted: make(chan interface{}),
	}
	clientConnection.On("Close")
	clientEventSender.close()
	mock.AssertExpectationsForObjects(t, clientConnection)
}
