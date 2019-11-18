package engine

import (
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/environment/character"
	"francoisgergaud/3dGame/environment/world"
	"francoisgergaud/3dGame/environment/worldelement"
	"francoisgergaud/3dGame/internal/testutils"
	testcharacter "francoisgergaud/3dGame/internal/testutils/environment/character"
	testworldelement "francoisgergaud/3dGame/internal/testutils/environment/worldelement"
	"francoisgergaud/3dGame/internal/testutils/environment/worldmap"
	"testing"
	"time"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBackgroundRenderer struct {
	mock.Mock
}

func (mock *MockBackgroundRenderer) Render(worldMap world.WorldMap, player character.Character, worldElements []worldelement.WorldElement, screen tcell.Screen) {
	mock.Called(worldMap, player, worldElements, screen)
}

func TestStartEngine(t *testing.T) {
	screen := new(testutils.MockScreen)
	worldMap := new(worldmap.MockWorldMap)
	playerUpdateChannel := make(chan time.Time)
	playerQuitChannel := make(chan struct{})
	worldElementUpdateChannel := make(chan time.Time)
	worldElementQuitChannel := make(chan struct{})
	player := &testcharacter.MockCharacter{
		QuitChannel:   playerQuitChannel,
		UpdateChannel: playerUpdateChannel,
	}
	worldElements := make([]worldelement.WorldElement, 1)
	worldElement := &testworldelement.MockWorldElement{
		QuitChannel:   worldElementQuitChannel,
		UpdateChannel: worldElementUpdateChannel,
	}
	worldElements[0] = worldElement
	bgRender := new(MockBackgroundRenderer)
	consoleEventListener := new(testutils.MockConsoleEventManager)
	//to shorten the test of the timer. A ticker is generated every 1000/250 ms
	frameRate := 250
	updateRate := 500
	screen.On("Clear")
	screen.On("SetStyle", tcell.StyleDefault)
	consoleEventListener.On("Listen")
	bgRender.On("Render", worldMap, player, worldElements, screen)
	player.On("Start")
	player.On("GetUpdateChannel").Return(playerUpdateChannel)
	player.On("GetQuitChannel").Return(playerQuitChannel)
	worldElement.On("Start")
	worldElement.On("GetUpdateChannel").Return(worldElementUpdateChannel)
	worldElement.On("GetQuitChannel").Return(worldElementQuitChannel)
	quit := make(chan struct{})
	engine := Impl{
		screen:              screen,
		player:              player,
		worldMap:            worldMap,
		worldElements:       worldElements,
		bgRender:            bgRender,
		consoleEventManager: consoleEventListener,
		quit:                quit,
		frameRate:           frameRate,
		updateRate:          updateRate,
	}
	go func() {
		<-time.After(time.Millisecond * time.Duration((2*1000)/frameRate))
		close(quit)
	}()
	engine.StartEngine()
}

func TestNewEngine(t *testing.T) {
	screen := new(testutils.MockScreen)
	engineConfig := new(Configuration)
	grid := [][]int{}
	engineConfig.WorldMap = grid
	playerConfiguration := new(PlayerConfiguration)
	playerPosition := &common.Point2D{X: 2.6, Y: -2.9}
	playerAngle := 0.5
	playerConfiguration.InitialPosition = playerPosition
	playerConfiguration.InitialAngle = playerAngle
	engineConfig.PlayerConfiguration = playerConfiguration
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
	worldElementInitialPosition := &common.Point2D{X: 2.3, Y: -0.9}
	worldElementInitialAngle := 0.2
	worldElementSize := 0.5
	worldElementVelocity := 1.5
	worldElementStyle := tcell.StyleDefault.Foreground(tcell.Color100)
	engineConfig.WorldElementConfigurations = []WorldElementConfiguration{
		WorldElementConfiguration{
			InitialPosition: worldElementInitialPosition,
			InitialAngle:    worldElementInitialAngle,
			Size:            worldElementSize,
			Velocity:        worldElementVelocity,
			Style:           worldElementStyle,
		},
	}
	frameRate := 40
	worldUpdateRate := 50
	engineConfig.FrameRate = frameRate
	engineConfig.WorlUpdateRate = worldUpdateRate
	engine, err := NewEngine(screen, engineConfig)
	assert.Nil(t, err)
	assert.Equal(t, screen, engine.screen)
	assert.Equal(t, frameRate, engine.frameRate)
	assert.Equal(t, worldUpdateRate, engine.updateRate)
	assert.Equal(t, playerPosition, engine.player.GetPosition())
	assert.Equal(t, playerAngle, engine.player.GetAngle())
}
