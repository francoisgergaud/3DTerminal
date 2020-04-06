package impl

import (
	"francoisgergaud/3dGame/client/configuration"
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/client/render/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/runner"
	testclient "francoisgergaud/3dGame/internal/testutils/client"
	testconnector "francoisgergaud/3dGame/internal/testutils/client/connector"
	testConsoleManager "francoisgergaud/3dGame/internal/testutils/client/consolemanager"
	testPlayer "francoisgergaud/3dGame/internal/testutils/client/player"
	testanimatedelement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	testrunner "francoisgergaud/3dGame/internal/testutils/common/runner"
	testtcell "francoisgergaud/3dGame/internal/testutils/tcell"
	"testing"
	"time"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBackgroundRenderer struct {
	mock.Mock
}

func (mock *MockBackgroundRenderer) Render(playerID string, worldMap world.WorldMap, player player.Player, worldElements map[string]animatedelement.AnimatedElement, screen tcell.Screen) {
	mock.Called(playerID, worldMap, player, worldElements, screen)
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
	assert.NotNil(t, engine.playerFactory)
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
	player := &testPlayer.MockPlayer{}
	playerID := "fakePlayerID"
	worldElements := make(map[string]animatedelement.AnimatedElement)
	bgRender := new(MockBackgroundRenderer)
	//to shorten the test of the timer. A ticker is generated every 1000/250 ms
	frameRate := 1000
	screen.On("Clear")
	screen.On("SetStyle", tcell.StyleDefault)
	screen.On("Fini")
	bgRender.On("Render", playerID, worldMap, player, worldElements, screen)
	shutdown := make(chan interface{})
	connectionToServer := new(testconnector.MockServerConnection)
	connectionToServer.On("Disconnect")
	engine := Impl{
		screen:             screen,
		player:             player,
		playerID:           playerID,
		worldMap:           worldMap,
		otherPlayers:       worldElements,
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
	player := &testPlayer.MockPlayer{}
	worldElements := make(map[string]animatedelement.AnimatedElement)
	worldElement := &testanimatedelement.MockAnimatedElement{}
	worldElements["worldElementID"] = worldElement

	player.MockAnimatedElement.On("Move")
	worldElement.On("Move")

	engine := new(testclient.MockEngine)
	engine.On("Player").Return(player)
	engine.On("OtherPlayers").Return(worldElements)

	worldElementUpdater := worldElementUpdaterImpl{
		updateRate: 1000,
		engine:     engine,
		quit:       quitChannel,
	}
	go worldElementUpdater.Run()
	<-time.After(time.Millisecond * 5)
	close(quitChannel)
	mock.AssertExpectationsForObjects(t, player, worldElement, engine)
}

func TestReceiveJoinEventFromServer(t *testing.T) {
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

func TestReceiveMovePastEventFromServer(t *testing.T) {
	otherPlayerID := "otherPlayerID"
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	engine := &Impl{
		otherPlayers: otherPlayers,
		initialized:  true,
	}
	mockAnimatedElement := testanimatedelement.MockAnimatedElement{}
	otherPlayers[otherPlayerID] = &mockAnimatedElement
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

func TestReceiveMoveEventFromServer(t *testing.T) {
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
func TestReceiveInitializationEvents(t *testing.T) {
	playerID := "playerID"
	playerState := state.AnimatedElementState{}
	worldMap := new(testworld.MockWorldMap)
	worldMapCloned := new(testworld.MockWorldMap)
	worldMap.On("Clone").Return(worldMapCloned)
	otherPlayerStates := make(map[string]state.AnimatedElementState)
	otherPlayerID := "otherPlayerID"
	otherPlayerState := state.AnimatedElementState{}
	otherPlayerStates[otherPlayerID] = otherPlayerState
	playerFactory := testPlayer.MockPlayerFactory{}
	mathHelper := new(testhelper.MockMathHelper)
	quit := make(chan interface{})
	playerListener := &playerListenerImpl{}
	player := new(testPlayer.MockPlayer)
	player.MockEventPublisher.On("RegisterListener", playerListener)
	playerFactory.On("NewPlayer", &playerState, worldMapCloned, mathHelper).Return(player)
	consoleEventManager := new(testConsoleManager.MockConsoleEventManager)
	consoleEventManager.On("SetPlayer", player)
	animatedElementFactory := testanimatedelement.MockAnimatedElementFactory{}
	otherPlayerAnimatedElement := new(testanimatedelement.MockAnimatedElement)
	animatedElementFactory.On("NewAnimatedElementWithState", &otherPlayerState, worldMapCloned, mathHelper).Return(otherPlayerAnimatedElement)
	runner := new(testrunner.MockRunner)
	worldElementUpdater := &worldElementUpdaterImpl{}
	engine := &Impl{
		initialized:                           false,
		mathHelper:                            mathHelper,
		playerFactory:                         playerFactory.NewPlayer,
		animatedElementFactory:                animatedElementFactory.NewAnimatedElementWithState,
		consoleEventManager:                   consoleEventManager,
		quit:                                  quit,
		playerListener:                        playerListener,
		worldElementUpdater:                   worldElementUpdater,
		Runner:                                runner,
		preInitializationEventFromServerQueue: make(chan event.Event, 100),
	}
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
			"otherPlayers": map[string]state.AnimatedElementState{
				otherPlayerID: otherPlayerState,
			},
		},
	}
	engine.ReceiveEventsFromServer([]event.Event{preInitializationEvent, initEvent})
	assert.Equal(t, playerID, engine.playerID)
	assert.Equal(t, worldMapCloned, engine.worldMap)
	assert.Equal(t, player, engine.player)
	assert.Equal(t, otherPlayerAnimatedElement, engine.otherPlayers[otherPlayerID])
	assert.True(t, engine.initialized)
	mock.AssertExpectationsForObjects(t, player, worldMap, &playerFactory, consoleEventManager, &animatedElementFactory, runner)
}

func TestReceiveQuitEventFromServer(t *testing.T) {
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

func TestPlayerListenerReceiveEvent(t *testing.T) {
	playerEventQueue := make(chan event.Event)
	playerListener := playerListenerImpl{
		playerEventQueue: playerEventQueue,
	}
	event := event.Event{
		PlayerID: "fakePlayerID",
	}
	go playerListener.ReceiveEvent(event)
	eventReceived := <-playerEventQueue
	assert.Equal(t, eventReceived, event)
}
