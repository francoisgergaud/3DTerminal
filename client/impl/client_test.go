package impl

import (
	"francoisgergaud/3dGame/client/configuration"
	"francoisgergaud/3dGame/client/render/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/projectile"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/runner"
	testclient "francoisgergaud/3dGame/internal/testutils/client"
	testconnector "francoisgergaud/3dGame/internal/testutils/client/connector"
	testConsoleManager "francoisgergaud/3dGame/internal/testutils/client/consolemanager"
	testanimatedelement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	testprojectile "francoisgergaud/3dGame/internal/testutils/common/environment/projectile"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	testrunner "francoisgergaud/3dGame/internal/testutils/common/runner"
	testtcell "francoisgergaud/3dGame/internal/testutils/tcell"
	"testing"
	"time"

	"github.com/gdamore/tcell"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBackgroundRenderer struct {
	mock.Mock
}

func (mock *MockBackgroundRenderer) Render(playerID string, worldMap world.WorldMap, player animatedelement.AnimatedElement, worldElements map[string]animatedelement.AnimatedElement, projectiles map[string]projectile.Projectile, screen tcell.Screen) {
	mock.Called(playerID, worldMap, player, worldElements, projectiles, screen)
}

type MockFactories struct {
	mock.Mock
}

func (mock *MockFactories) NewID() uuid.UUID {
	args := mock.Called()
	return args.Get(0).(uuid.UUID)
}

func TestNewEngine(t *testing.T) {
	screen := new(testtcell.MockScreen)
	engineConfig := &configuration.Configuration{
		GradientRSBackgroundRange:  []float32{0.5},
		GradientRSBackgroundColors: []int{0, 1},
		GradientRSMultiplicator:    2.0,
		GradientRSLimit:            3.0,
		GradientRSFirst:            0.5,
		FrameRate:                  40,
		WorlUpdateRate:             50,
	}
	consoleManager := new(testConsoleManager.MockConsoleEventManager)
	quit := make(chan interface{})
	engine, err := NewEngine(screen, consoleManager, engineConfig, quit)
	assert.Nil(t, err)
	assert.Equal(t, screen, engine.screen)
	assert.IsType(t, &helper.MathHelperImpl{}, engine.mathHelper)
	assert.IsType(t, &impl.RendererImpl{}, engine.renderer)
	assert.NotNil(t, engine.preInitializationEventFromServerQueue)
	assert.Equal(t, engineConfig.FrameRate, engine.frameRate)
	assert.True(t, quit == engine.quit)
	assert.Equal(t, consoleManager, engine.consoleEventManager)
	assert.NotNil(t, engine.shutdown)
	assert.False(t, engine.initialized)
	assert.NotNil(t, engine.animatedElementFactory)
	//test the player-listener
	assert.True(t, quit == engine.playerListener.quit)
	assert.NotNil(t, engine.playerListener.playerEventQueue)
	assert.Nil(t, engine.playerListener.connectionToServer)
	assert.Equal(t, engineConfig.WorlUpdateRate, engine.worldElementUpdater.updateRate)
	assert.True(t, quit == engine.worldElementUpdater.quit)
	assert.Equal(t, engine, engine.worldElementUpdater.engine)
	assert.IsType(t, &runner.AsyncRunner{}, engine.Runner)
	mock.AssertExpectationsForObjects(t, screen)
}

func TestEngineRun(t *testing.T) {
	screen := new(testtcell.MockScreen)
	worldMap := new(testworld.MockWorldMap)
	quitChannel := make(chan interface{})
	player := new(testanimatedelement.MockAnimatedElement)
	playerID := "fakePlayerID"
	worldElements := make(map[string]animatedelement.AnimatedElement)
	projectiles := make(map[string]projectile.Projectile)
	bgRender := new(MockBackgroundRenderer)
	//to shorten the test of the timer. A ticker is generated every 1000/250 ms
	frameRate := 1000
	screen.On("Clear")
	screen.On("SetStyle", tcell.StyleDefault)
	screen.On("Fini")
	bgRender.On("Render", playerID, worldMap, player, worldElements, projectiles, screen)
	shutdown := make(chan interface{})
	connectionToServer := new(testconnector.MockServerConnection)
	connectionToServer.On("Disconnect")
	engine := Impl{
		screen:             screen,
		player:             player,
		playerID:           playerID,
		worldMap:           worldMap,
		otherPlayers:       worldElements,
		projectiles:        projectiles,
		renderer:           bgRender,
		quit:               quitChannel,
		frameRate:          frameRate,
		shutdown:           shutdown,
		connectionToServer: connectionToServer,
	}
	//Run is blocking
	go engine.Run()
	<-time.After(time.Millisecond * 2)
	close(quitChannel)
	<-shutdown
	mock.AssertExpectationsForObjects(t, bgRender, screen, player, connectionToServer)
}

func TestWorldUpdaterRun(t *testing.T) {
	quitChannel := make(chan interface{})
	player := new(testanimatedelement.MockAnimatedElement)
	worldElements := make(map[string]animatedelement.AnimatedElement)
	worldElement := &testanimatedelement.MockAnimatedElement{}
	worldElements["worldElementID"] = worldElement

	projectiles := make(map[string]projectile.Projectile)
	projectile := &testprojectile.MockProjectile{}
	projectiles["projectileID"] = projectile

	player.On("Move")
	worldElement.On("Move")
	projectile.MockAnimatedElement.On("Move")

	engine := new(testclient.MockEngine)
	engine.On("Player").Return(player)
	engine.On("OtherPlayers").Return(worldElements)
	engine.On("Projectiles").Return(projectiles)

	worldElementUpdater := worldElementUpdaterImpl{
		updateRate: 1000,
		engine:     engine,
		quit:       quitChannel,
	}
	go worldElementUpdater.Run()
	<-time.After(time.Millisecond * 5)
	close(quitChannel)
	mock.AssertExpectationsForObjects(t, player, worldElement, projectile, engine)
}

func TestReceiveEventFromServerJoin(t *testing.T) {
	engine := &Impl{
		otherPlayers:           make(map[string]animatedelement.AnimatedElement),
		otherPlayerLastUpdates: make(map[string]uint32),
		initialized:            true,
	}
	events := make([]event.Event, 0)
	newPlayerState := state.AnimatedElementState{}
	events = append(events,
		event.Event{
			Action:   "join",
			PlayerID: "player1",
			State:    &newPlayerState,
		},
	)
	engine.ReceiveEventsFromServer(events)
	playerRegistered, ok := engine.otherPlayers["player1"]
	assert.True(t, ok)
	assert.Equal(t, &newPlayerState, playerRegistered.State())
}

func TestReceiveEventFromServerMoveWithOldTimeframe(t *testing.T) {
	otherPlayerID := "otherPlayerID"
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	mockAnimatedElement := testanimatedelement.MockAnimatedElement{}
	otherPlayers[otherPlayerID] = &mockAnimatedElement
	engine := &Impl{
		otherPlayers: otherPlayers,
		initialized:  true,
	}
	events := make([]event.Event, 0)
	otherPlayerState := state.AnimatedElementState{}
	events = append(events,
		event.Event{
			Action:   "move",
			PlayerID: otherPlayerID,
			State:    &otherPlayerState,
		},
	)
	engine.ReceiveEventsFromServer(events)
	mock.AssertExpectationsForObjects(t, &mockAnimatedElement)
}

func TestReceiveEventFromServerMove(t *testing.T) {
	otherPlayerID := "otherPlayerID"
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	otherPlayerLastUpdates := make(map[string]uint32)
	otherPlayerLastUpdates["otherPlayerID"] = 1
	engine := &Impl{
		otherPlayers:           otherPlayers,
		initialized:            true,
		otherPlayerLastUpdates: otherPlayerLastUpdates,
	}
	mockAnimatedElement := testanimatedelement.MockAnimatedElement{}
	otherPlayers[otherPlayerID] = &mockAnimatedElement
	events := make([]event.Event, 0)
	otherPlayerState := state.AnimatedElementState{}
	events = append(events,
		event.Event{
			Action:    "move",
			PlayerID:  otherPlayerID,
			State:     &otherPlayerState,
			TimeFrame: 2,
		},
	)
	mockAnimatedElement.On("SetState", &otherPlayerState)
	engine.ReceiveEventsFromServer(events)
	mock.AssertExpectationsForObjects(t, &mockAnimatedElement)
}

//TODO: decompose the client: this method is too complex to test
func TestReceiveEventsFromServerInit(t *testing.T) {
	playerID := "playerID"
	playerState := state.AnimatedElementState{}
	worldMap := new(testworld.MockWorldMap)
	otherPlayerStates := make(map[string]state.AnimatedElementState)
	otherPlayerID := "otherPlayerIDTest1"
	otherPlayerState := state.AnimatedElementState{}
	otherPlayerStates[otherPlayerID] = otherPlayerState
	projectileStates := make(map[string]state.AnimatedElementState)
	projectileID := "projectileTest1"
	projectilePosition := &math.Point2D{
		X: 0.025,
		Y: 0.025,
	}
	projectileAngle := 0.075
	projectileState := state.AnimatedElementState{
		Position: projectilePosition,
		Angle:    projectileAngle,
	}
	projectileStates[projectileID] = projectileState
	projectile := new(testprojectile.MockProjectile)
	mathHelper := new(testhelper.MockMathHelper)
	quit := make(chan interface{})
	playerListener := &playerListenerImpl{}
	consoleEventManager := new(testConsoleManager.MockConsoleEventManager)
	animatedElementFactory := testanimatedelement.MockAnimatedElementFactory{}
	player := new(testanimatedelement.MockAnimatedElement)
	otherPlayerAnimatedElement := new(testanimatedelement.MockAnimatedElement)
	animatedElementFactory.On("NewAnimatedElementWithState", playerID, &playerState, worldMap, mathHelper).Return(player)
	animatedElementFactory.On("NewAnimatedElementWithState", otherPlayerID, &otherPlayerState, worldMap, mathHelper).Return(otherPlayerAnimatedElement)
	projectileFactory := new(testprojectile.MockProjectileFactory)
	projectileFactory.On("CreateProjectile", projectileID, projectilePosition, projectileAngle, worldMap, mock.MatchedBy(
		func(otherPlayers map[string]animatedelement.AnimatedElement) bool {
			for id := range otherPlayers {
				if id == otherPlayerID {
					return true
				}
			}
			return false
		},
	), mathHelper).Return(projectile)
	runner := new(testrunner.MockRunner)
	worldElementUpdater := &worldElementUpdaterImpl{}
	engine := &Impl{
		initialized:                           false,
		mathHelper:                            mathHelper,
		animatedElementFactory:                animatedElementFactory.NewAnimatedElementWithState,
		projectileFactory:                     projectileFactory.CreateProjectile,
		consoleEventManager:                   consoleEventManager,
		quit:                                  quit,
		playerListener:                        playerListener,
		worldElementUpdater:                   worldElementUpdater,
		Runner:                                runner,
		preInitializationEventFromServerQueue: make(chan event.Event, 100),
	}
	consoleEventManager.On("SetPlayer", engine)
	runner.On("Start", engine)
	runner.On("Start", worldElementUpdater)
	runner.On("Start", playerListener)
	runner.On("Start", consoleEventManager)
	preInitializationEvent := event.Event{}
	initEvent := event.Event{
		PlayerID: playerID,
		Action:   "init",
		State:    &playerState,
		ExtraData: map[string]interface{}{
			"worldMap": worldMap,
			"otherPlayers": map[string]*state.AnimatedElementState{
				otherPlayerID: &otherPlayerState,
			},
			"projectiles": map[string]*state.AnimatedElementState{
				projectileID: &projectileState,
			},
		},
	}

	engine.ReceiveEventsFromServer([]event.Event{preInitializationEvent, initEvent})

	assert.Equal(t, playerID, engine.playerID)
	assert.Equal(t, worldMap, engine.worldMap)
	assert.Equal(t, player, engine.player)
	assert.Equal(t, otherPlayerAnimatedElement, engine.otherPlayers[otherPlayerID])
	assert.Equal(t, projectile, engine.projectiles[projectileID])
	assert.True(t, engine.initialized)
	mock.AssertExpectationsForObjects(t, player, worldMap, consoleEventManager, &animatedElementFactory, runner, projectileFactory)
}

func TestReceiveEventsFromServerQuit(t *testing.T) {
	otherPlayerID := "otherPlayerID"
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	otherPlayerLastUpdates := make(map[string]uint32)
	otherPlayerLastUpdates[otherPlayerID] = 1
	engine := &Impl{
		otherPlayers:           otherPlayers,
		initialized:            true,
		otherPlayerLastUpdates: otherPlayerLastUpdates,
	}
	mockAnimatedElement := testanimatedelement.MockAnimatedElement{}
	otherPlayers[otherPlayerID] = &mockAnimatedElement
	events := make([]event.Event, 0)
	events = append(events,
		event.Event{
			Action:   "quit",
			PlayerID: otherPlayerID,
		},
	)
	engine.ReceiveEventsFromServer(events)
	assert.NotContains(t, engine.otherPlayers, otherPlayerID)
	assert.NotContains(t, engine.otherPlayerLastUpdates, otherPlayerID)
}

func TestReceiveEventsFromServerKillOtherPlayer(t *testing.T) {
	otherPlayerID := "otherPlayerID"
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	otherPlayerLastUpdates := make(map[string]uint32)
	otherPlayerLastUpdates[otherPlayerID] = 1
	engine := &Impl{
		otherPlayers:           otherPlayers,
		initialized:            true,
		otherPlayerLastUpdates: otherPlayerLastUpdates,
	}
	mockAnimatedElement := testanimatedelement.MockAnimatedElement{}
	otherPlayers[otherPlayerID] = &mockAnimatedElement
	events := make([]event.Event, 0)
	events = append(events,
		event.Event{
			Action:   "kill",
			PlayerID: otherPlayerID,
		},
	)
	engine.ReceiveEventsFromServer(events)
	assert.NotContains(t, engine.otherPlayers, otherPlayerID)
	assert.NotContains(t, engine.otherPlayerLastUpdates, otherPlayerID)
}

func TestReceiveEventsFromServerFire(t *testing.T) {
	projectileID := "projectileID"
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	projectileFactoryBuilder := new(testprojectile.MockProjectileFactory)
	mathHelper := new(testhelper.MockMathHelper)
	worldMap := new(testworld.MockWorldMap)
	engine := &Impl{
		playerID:          "playerID",
		otherPlayers:      otherPlayers,
		projectileFactory: projectileFactoryBuilder.CreateProjectile,
		mathHelper:        mathHelper,
		worldMap:          worldMap,
		initialized:       true,
		projectiles:       make(map[string]projectile.Projectile),
	}
	position := &math.Point2D{}
	angle := 0.765
	events := make([]event.Event, 0)
	events = append(events,
		event.Event{
			PlayerID: "otherPlayerID",
			Action:   "fire",
			State: &state.AnimatedElementState{
				Position: position,
				Angle:    angle,
			},
			ExtraData: map[string]interface{}{
				"projectileID": projectileID,
			},
		},
	)
	projectileToReturn := &projectile.ProjectileImpl{}
	projectileFactoryBuilder.On("CreateProjectile", projectileID, position, angle, worldMap, otherPlayers, mathHelper).Return(projectileToReturn)

	engine.ReceiveEventsFromServer(events)

	mock.AssertExpectationsForObjects(t, projectileFactoryBuilder)
	assert.Same(t, projectileToReturn, engine.projectiles[projectileID])
}

func TestReceiveEventsFromServerProjectileImpact(t *testing.T) {
	projectiles := make(map[string]projectile.Projectile)
	projectileID := "projectileID"
	projectiles[projectileID] = &projectile.ProjectileImpl{}
	engine := &Impl{
		playerID:    "playerID",
		projectiles: projectiles,
		initialized: true,
	}
	events := make([]event.Event, 0)
	events = append(events,
		event.Event{
			PlayerID: projectileID,
			Action:   "projectileImpact",
		},
	)

	engine.ReceiveEventsFromServer(events)

	assert.NotContains(t, projectiles, projectileID)
}

func TestReceiveEventsFromServerKillPlayer(t *testing.T) {
	playerID := "playerID"
	engine := &Impl{
		playerID:            playerID,
		initialized:         true,
		waitSpawnFromServer: false,
	}
	events := make([]event.Event, 0)
	events = append(events,
		event.Event{
			PlayerID: playerID,
			Action:   "kill",
		},
	)

	engine.ReceiveEventsFromServer(events)

	assert.True(t, engine.waitSpawnFromServer)
}

func TestReceiveEventsFromServerSpawnPlayer(t *testing.T) {
	playerID := "playerID"
	player := new(testanimatedelement.MockAnimatedElement)
	engine := &Impl{
		playerID:            playerID,
		initialized:         true,
		waitSpawnFromServer: true,
		player:              player,
	}
	events := make([]event.Event, 0)
	stateForSpawn := &state.AnimatedElementState{}
	events = append(events,
		event.Event{
			PlayerID: playerID,
			Action:   "spawn",
			State:    stateForSpawn,
		},
	)
	player.On("SetState", stateForSpawn)

	engine.ReceiveEventsFromServer(events)

	assert.False(t, engine.waitSpawnFromServer)
	mock.AssertExpectationsForObjects(t, player)
}

func TestPlayerListenerRun(t *testing.T) {
	quit := make(chan interface{})
	playerEventQueue := make(chan event.Event)
	eventFromPlayer := event.Event{}
	serverConnection := new(testconnector.MockServerConnection)
	serverConnection.On("NotifyServer", []event.Event{eventFromPlayer}).Return(nil)
	playerListener := playerListenerImpl{
		playerEventQueue:   playerEventQueue,
		quit:               quit,
		connectionToServer: serverConnection,
	}
	go playerListener.Run()
	playerEventQueue <- eventFromPlayer
	<-time.After(time.Millisecond * 2)
	close(quit)
	mock.AssertExpectationsForObjects(t, serverConnection)
}

func TestOtherPlayers(t *testing.T) {
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	engine := &Impl{
		otherPlayers: otherPlayers,
	}
	assert.Equal(t, otherPlayers, engine.OtherPlayers())
}

func TestShutdown(t *testing.T) {
	shutdown := make(chan interface{})
	engine := &Impl{
		shutdown: shutdown,
	}
	close(shutdown)
	engine.Shutdown()
}

func TestConnectToServer(t *testing.T) {
	serverConnection := new(testconnector.MockServerConnection)
	engine := &Impl{
		playerListener: &playerListenerImpl{},
	}
	engine.ConnectToServer(serverConnection)
	assert.Same(t, serverConnection, engine.connectionToServer)
	assert.Same(t, serverConnection, engine.playerListener.connectionToServer)
}

func playerMoveTest(t *testing.T, initialRotationDirection, expectedRotationDirection, initialMoveDirection, expectedMoveDirection state.Direction, eventKey *tcell.EventKey) {
	playerState := state.AnimatedElementState{
		Position:        &math.Point2D{X: 1, Y: 3},
		Angle:           0.0,
		StepAngle:       0.01,
		Size:            0.5,
		Velocity:        1.0,
		Style:           tcell.StyleDefault.Background(tcell.Color111),
		MoveDirection:   initialMoveDirection,
		RotateDirection: initialRotationDirection,
	}
	player := new(testanimatedelement.MockAnimatedElement)
	player.On("State").Return(&playerState)
	playerEventQueue := make(chan event.Event, 1)
	engine := &Impl{
		player: player,
		playerListener: &playerListenerImpl{
			playerEventQueue: playerEventQueue,
		},
		waitSpawnFromServer: false,
	}
	engine.Action(eventKey)
	assert.Equal(t, expectedRotationDirection, playerState.RotateDirection)
	assert.Equal(t, expectedMoveDirection, playerState.MoveDirection)
	eventSent := <-playerEventQueue
	assert.Equal(t, "move", eventSent.Action)
	assert.Equal(t, player.State(), eventSent.State)
}

func TestMoveAction(t *testing.T) {
	playerMoveTest(t, state.None, state.None, state.None, state.Forward, tcell.NewEventKey(tcell.KeyUp, 0, 0))
	playerMoveTest(t, state.None, state.None, state.Backward, state.None, tcell.NewEventKey(tcell.KeyUp, 0, 0))
	playerMoveTest(t, state.None, state.None, state.None, state.Backward, tcell.NewEventKey(tcell.KeyDown, 0, 0))
	playerMoveTest(t, state.None, state.None, state.Forward, state.None, tcell.NewEventKey(tcell.KeyDown, 0, 0))
	playerMoveTest(t, state.None, state.Left, state.None, state.None, tcell.NewEventKey(tcell.KeyLeft, 0, 0))
	playerMoveTest(t, state.Right, state.None, state.None, state.None, tcell.NewEventKey(tcell.KeyLeft, 0, 0))
	playerMoveTest(t, state.None, state.Right, state.None, state.None, tcell.NewEventKey(tcell.KeyRight, 0, 0))
	playerMoveTest(t, state.Left, state.None, state.None, state.None, tcell.NewEventKey(tcell.KeyRight, 0, 0))
}

func TestFireAction(t *testing.T) {
	playerState := state.AnimatedElementState{
		Position:        &math.Point2D{X: 1, Y: 3},
		Angle:           1.5,
		StepAngle:       0.01,
		Size:            0.5,
		Velocity:        1.0,
		Style:           tcell.StyleDefault.Background(tcell.Color111),
		MoveDirection:   state.None,
		RotateDirection: state.None,
	}
	player := new(testanimatedelement.MockAnimatedElement)
	player.On("State").Return(&playerState)
	playerEventQueue := make(chan event.Event, 1)
	projectileFactoryBuilder := new(testprojectile.MockProjectileFactory)
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	mathHelper := new(testhelper.MockMathHelper)
	worldMap := new(testworld.MockWorldMap)
	mockFactories := new(MockFactories)
	playerID := "playerID"
	randomID := uuid.New()
	engine := &Impl{
		playerID: playerID,
		player:   player,
		playerListener: &playerListenerImpl{
			playerEventQueue: playerEventQueue,
		},
		projectileFactory:   projectileFactoryBuilder.CreateProjectile,
		waitSpawnFromServer: false,
		otherPlayers:        otherPlayers,
		mathHelper:          mathHelper,
		worldMap:            worldMap,
		identifierFactory:   mockFactories.NewID,
		projectiles:         make(map[string]projectile.Projectile),
	}
	mockFactories.On("NewID").Return(randomID)
	epextedPosition := &math.Point2D{X: 0.9999999999999999, Y: 2.25}
	projectileToReturn := new(testprojectile.MockProjectile)
	projectileState := &state.AnimatedElementState{}
	expectedProjectileID := playerID + "." + randomID.String()
	projectileFactoryBuilder.On("CreateProjectile", expectedProjectileID, epextedPosition, playerState.Angle, worldMap, otherPlayers, mathHelper).Return(projectileToReturn)
	projectileToReturn.MockAnimatedElement.On("State").Return(projectileState)

	engine.Action(tcell.NewEventKey(tcell.KeyEnter, 0, 0))

	assert.Same(t, projectileToReturn, engine.projectiles[expectedProjectileID])
	eventSentToServer := <-playerEventQueue
	assert.Equal(t, "fire", eventSentToServer.Action)
	assert.Equal(t, expectedProjectileID, eventSentToServer.ExtraData["projectileID"])
	assert.Same(t, projectileState, eventSentToServer.State)
	mock.AssertExpectationsForObjects(t, projectileFactoryBuilder, mockFactories)
}

func TestActionWhithWiatingSpwanFromServer(t *testing.T) {
	playerState := state.AnimatedElementState{}
	player := new(testanimatedelement.MockAnimatedElement)
	player.On("State").Return(&playerState)
	engine := &Impl{
		player:              player,
		waitSpawnFromServer: true,
	}
	engine.Action(tcell.NewEventKey(tcell.KeyUp, 0, 0))
}
