package engine

import (
	"francoisgergaud/3dGame/engine/event"
	"francoisgergaud/3dGame/engine/render"
	"francoisgergaud/3dGame/environment"
	"time"

	"github.com/gdamore/tcell"
)

//Configuration contains the required parametrable parameters for the engine.
type Configuration struct {
	//The frame-rate per second.
	FrameRate int
	//the world-update's rate.
	WorlUpdateRate int
	//The screen's height.
	ScreenHeight int
	//The screen's width.
	ScreenWidth int
	//The player's (or camera) field-of-view angle in Pie radian.
	PlayerFieldOfViewAngle float64
	//The player's (or camera) maximum's visibility.
	Visibility float64
	//The gradient-ray-sampler first distance upper-range value.
	GradientRSFirst float64
	//The gradient-ray-sampler distance exponential multiplicator.
	GradientRSMultiplicator float64
	//The gradient-ray-sampler distance maximum upper-range. After this range, the last gradient-color will be used until infinit.
	GradientRSLimit float64
	//The gradient-ray-sampler start-color value (closer color).
	GradientRSWallStartColor int
	//The gradient-ray-sampler end-color value (farest color).
	GradientRSWallEndColor int
	//The gradient-ray-sampler background-column-index ratio, from 0 to 1. The last value must be 1.0, and the values must be increasing
	GradientRSBackgroundRange []float32
	//The gradient-ray-sampler background-colors, which apply to the upper-range ratio of the row defined in GradientRSBackgroundRange.
	GradientRSBackgroundColors []int
}

//Engine represents an game engine which takes input and render on screen.
type Engine interface {
	StartEngine()
}

//Impl implements the Engine interface.
type Impl struct {
	screen              tcell.Screen
	player              environment.Character
	worldMap            environment.WorldMap
	bgRender            render.BackgroundRenderer
	consoleEventManager event.ConsoleEventManager
	quit                chan struct{}
	frameRate           int
	updateRate          int
}

//NewEngine provides a new engine.
func NewEngine(screen tcell.Screen, player environment.Character, worldMap environment.WorldMap, engineConfig *Configuration) *Impl {
	raySampler := render.CreateRaySamplerForAnsiColorTerminal(
		engineConfig.GradientRSFirst,
		engineConfig.GradientRSMultiplicator,
		engineConfig.GradientRSLimit,
		engineConfig.GradientRSWallStartColor,
		engineConfig.GradientRSWallEndColor,
		engineConfig.ScreenHeight,
		engineConfig.GradientRSBackgroundRange,
		engineConfig.GradientRSBackgroundColors)
	backgroundColumnRenderer := render.CreateBackgroundColumnRenderer(
		engineConfig.ScreenWidth,
		engineConfig.ScreenHeight,
		engineConfig.PlayerFieldOfViewAngle,
		engineConfig.Visibility,
		render.NewBackgroundRendererMathHelper(new(render.RayCasterImpl)),
		tcell.StyleDefault.Background(tcell.ColorBlueViolet),
		raySampler)
	bgRender := render.CreateBackgroundRenderer(engineConfig.ScreenWidth, backgroundColumnRenderer)
	quit := make(chan struct{})
	consoleEventManager := event.NewConsoleEventManager(screen, player, quit)
	engine := Impl{
		screen:              screen,
		player:              player,
		worldMap:            worldMap,
		bgRender:            bgRender,
		consoleEventManager: consoleEventManager,
		quit:                quit,
		frameRate:           engineConfig.FrameRate,
		updateRate:          engineConfig.WorlUpdateRate,
	}
	return &engine
}

//StartEngine initializes the required element and start the engine to render world's elements in pseudo-3D
func (engine *Impl) StartEngine() {
	engine.screen.Clear()
	engine.player.Start()
	go engine.consoleEventManager.Listen()
	frameUpdateTicker := time.NewTicker(time.Duration(1000/engine.frameRate) * time.Millisecond)
	worldUpdateTicker := time.NewTicker(time.Duration(1000/engine.updateRate) * time.Millisecond)
	for {
		select {
		case <-engine.quit:
			frameUpdateTicker.Stop()
			worldUpdateTicker.Stop()
			engine.screen.Fini()
			close(engine.player.GetQuitChannel())
			return
		case <-frameUpdateTicker.C:
			engine.bgRender.Render(engine.worldMap, engine.player, engine.screen)
		case updateWorldTickerTime := <-worldUpdateTicker.C:
			engine.player.GetUpdateChannel() <- updateWorldTickerTime
		}
	}

}
