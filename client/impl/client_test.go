package impl

import (
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
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

func (mock *MockBackgroundRenderer) Render(playerID string, worldMap world.WorldMap, player player.Player, worldElements map[string]animatedelement.AnimatedElement, screen tcell.Screen) {
	mock.Called(playerID, worldMap, player, worldElements, screen)
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
	playerID := "fakePlayerID"
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
	bgRender.On("Render", playerID, worldMap, player, worldElements, screen)
	player.On("Start")
	player.On("GetUpdateChannel")
	worldElement.On("Start")
	worldElement.On("GetUpdateChannel")
	engine := Impl{
		screen:       screen,
		player:       player,
		playerID:     playerID,
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
	consoleManager := new(testConsoleManager.MockConsoleEventManager)
	consoleManager.On("Listen")
	engine, err := NewEngine(screen, consoleManager, engineConfig)
	assert.Nil(t, err)
	assert.Equal(t, screen, engine.screen)
	assert.Equal(t, engineConfig.FrameRate, engine.frameRate)
	assert.Equal(t, engineConfig.WorlUpdateRate, engine.updateRate)
}

func TestReceiveJoinEventFromServer(t *testing.T) {
	//engine, _ := NewEngine(new(testtcell.MockScreen), createEngineConfigForTest())
	engine := &Impl{
		otherPlayers:           make(map[string]animatedelement.AnimatedElement),
		otherPlayerLastUpdates: make(map[string]uint32),
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
	assert.Equal(t, &newPlayerState, playerRegistered.GetState())
}

func TestReceiveMoveEventFromServer(t *testing.T) {
	otherPlayerID := "otherPlayerID"
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	engine := &Impl{
		otherPlayers: otherPlayers,
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
	mockAnimatedElement.On("SetState", &otherPlayerState)
	engine.ReceiveEventsFromServer(events)
}

func createEngineConfigForTest() *client.Configuration {
	engineConfig := new(client.Configuration)
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
	frameRate := 40
	worldUpdateRate := 50
	engineConfig.FrameRate = frameRate
	engineConfig.WorlUpdateRate = worldUpdateRate
	return engineConfig
}
