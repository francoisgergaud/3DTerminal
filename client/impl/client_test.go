package impl

import (
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/math"
	testconnector "francoisgergaud/3dGame/internal/testutils/client/connector"
	testConsoleManager "francoisgergaud/3dGame/internal/testutils/client/consolemanager"
	testPlayer "francoisgergaud/3dGame/internal/testutils/client/player"
	testanimatedelement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
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

func (mock *MockBackgroundRenderer) Render(worldMap world.WorldMap, player player.Player, worldElements map[string]animatedelement.AnimatedElement, screen tcell.Screen) {
	mock.Called(worldMap, player, worldElements, screen)
}

func TestStartEngine(t *testing.T) {
	screen := new(testtcell.MockScreen)
	worldMap := new(testworld.MockWorldMap)
	playerUpdateChannel := make(chan time.Time)
	quitChannel := make(chan struct{})
	worldElementUpdateChannel := make(chan time.Time)
	player := &testPlayer.MockPlayer{
		UpdateChannel: playerUpdateChannel,
		QuitChannel:   quitChannel,
	}
	worldElements := make(map[string]animatedelement.AnimatedElement)
	worldElement := &testanimatedelement.MockAnimatedElement{
		QuitChannel:   quitChannel,
		UpdateChannel: worldElementUpdateChannel,
	}
	worldElements["worldElementID"] = worldElement
	bgRender := new(MockBackgroundRenderer)
	consoleEventListener := new(testConsoleManager.MockConsoleEventManager)
	//to shorten the test of the timer. A ticker is generated every 1000/250 ms
	frameRate := 250
	updateRate := 500
	screen.On("Clear")
	screen.On("SetStyle", tcell.StyleDefault)
	consoleEventListener.On("Listen")
	bgRender.On("Render", worldMap, player, worldElements, screen)
	player.On("Start")
	player.On("GetUpdateChannel")
	worldElement.On("Start")
	worldElement.On("GetUpdateChannel")
	engine := Impl{
		screen:       screen,
		player:       player,
		worldMap:     worldMap,
		otherPlayers: worldElements,
		renderer:     bgRender,
		quit:         quitChannel,
		frameRate:    frameRate,
		updateRate:   updateRate,
	}
	go func() {
		<-time.After(time.Millisecond * time.Duration((2*1000)/frameRate))
		close(quitChannel)
	}()
	engine.StartEngine()
}

func TestNewEngine(t *testing.T) {
	screen := new(testtcell.MockScreen)
	engineConfig := createEngineConfigForTest()
	serverConnection := new(testconnector.MockServerConnection)
	engine, err := NewEngine(screen, engineConfig, serverConnection)
	assert.Nil(t, err)
	assert.Equal(t, screen, engine.screen)
	assert.Equal(t, engineConfig.FrameRate, engine.frameRate)
	assert.Equal(t, engineConfig.WorlUpdateRate, engine.updateRate)
	assert.Equal(t, &engineConfig.PlayerConfiguration, engine.GetPlayer().GetState())
	assert.Equal(t, engineConfig.OtherPlayerConfigurations["otherPlayerID"], *engine.otherPlayers["otherPlayerID"].GetState())
}

func TestReceiveJoinEventFromServer(t *testing.T) {
	engine, _ := NewEngine(new(testtcell.MockScreen), createEngineConfigForTest(), nil)
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
	assert.Equal(t, &newPlayerState, playerRegistered.GetState())
}

func TestReceiveMoveEventFromServer(t *testing.T) {
	engine, _ := NewEngine(new(testtcell.MockScreen), createEngineConfigForTest(), nil)
	events := make([]event.Event, 0)
	newPlayerState := state.AnimatedElementState{}
	events = append(events,
		event.Event{
			Action:   "move",
			PlayerID: "otherPlayerID",
			State:    &newPlayerState,
		},
	)
	engine.ReceiveEventsFromServer(events)
	playerRegistered, ok := engine.otherPlayers["otherPlayerID"]
	assert.True(t, ok)
	assert.Equal(t, &newPlayerState, playerRegistered.GetState())
}

func TestRegisterPlayer(t *testing.T) {
	engine, _ := NewEngine(new(testtcell.MockScreen), createEngineConfigForTest(), nil)
	events := make([]event.Event, 0)
	newPlayerState := state.AnimatedElementState{}
	events = append(events,
		event.Event{
			Action:   "move",
			PlayerID: "otherPlayerID",
			State:    &newPlayerState,
		},
	)
	engine.ReceiveEventsFromServer(events)
	playerRegistered, ok := engine.otherPlayers["otherPlayerID"]
	assert.True(t, ok)
	assert.Equal(t, &newPlayerState, playerRegistered.GetState())
}

func createEngineConfigForTest() *client.Configuration {
	engineConfig := new(client.Configuration)
	engineConfig.WorldMap = new(testworld.MockWorldMap)
	playerState := state.AnimatedElementState{
		Position: &math.Point2D{X: 2.6, Y: -2.9},
		Angle:    0.5,
	}
	engineConfig.PlayerConfiguration = playerState
	backgroundRange := []float32{0.5}
	engineConfig.GradientRSBackgroundRange = backgroundRange
	backgroundColors := []int{0, 1}
	engineConfig.GradientRSBackgroundColors = backgroundColors
	gradientRaySamplerMultiplicator := 2.0
	engineConfig.GradientRSMultiplicator = gradientRaySamplerMultiplicator
	gradientRaySamplerMaxLimit := 3.0
	engineConfig.GradientRSLimit = gradientRaySamplerMaxLimit
	gradientRaySamplerFirst := 0.5
	engineConfig.GradientRSFirst = gradientRaySamplerFirst
	engineConfig.OtherPlayerConfigurations = make(map[string]state.AnimatedElementState)
	otherPlayerState := state.AnimatedElementState{}
	engineConfig.OtherPlayerConfigurations["otherPlayerID"] = otherPlayerState
	frameRate := 40
	worldUpdateRate := 50
	engineConfig.FrameRate = frameRate
	engineConfig.WorlUpdateRate = worldUpdateRate
	return engineConfig
}
