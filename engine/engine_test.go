package engine

import (
	"francoisgergaud/3dGame/environment"
	"francoisgergaud/3dGame/internal/testutils"
	"testing"
	"time"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/mock"
)

type MockBackgroundRenderer struct {
	mock.Mock
}

func (mock *MockBackgroundRenderer) Render(worldMap environment.WorldMap, player environment.Character, screen tcell.Screen) {
	mock.Called(worldMap, player, screen)
}

func TestStartEngine(t *testing.T) {
	screen := new(testutils.MockScreen)
	worldMap := new(testutils.MockWorldMap)
	player := environment.NewPlayer(nil, 0.0, 0.0, worldMap)
	bgRender := new(MockBackgroundRenderer)
	consoleEventListener := new(testutils.MockConsoleEventManager)
	//to shorten the test of the timer. A ticker is generated every 1000/250 ms
	frameRate := 250
	screen.On("Clear")
	consoleEventListener.On("Listen")
	bgRender.On("Render", worldMap, player, screen)
	quit := make(chan struct{})
	engine := EngineImpl{
		screen:              screen,
		player:              player,
		worldMap:            worldMap,
		bgRender:            bgRender,
		consoleEventManager: consoleEventListener,
		quit:                quit,
		frameRate:           frameRate,
	}
	go func() {
		<-time.After(time.Millisecond * time.Duration((2*1000)/frameRate))
		close(quit)
	}()
	engine.StartEngine()
}
