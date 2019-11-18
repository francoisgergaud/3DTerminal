package engine

import (
	"fmt"
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/engine/event"
	"francoisgergaud/3dGame/engine/render"
	"francoisgergaud/3dGame/environment/character"
	"francoisgergaud/3dGame/environment/world"
	"francoisgergaud/3dGame/environment/worldelement"
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
	//The player's configuration
	PlayerConfiguration *PlayerConfiguration
	//the world-map
	WorldMap [][]int
	//the world-elements configurations.
	WorldElementConfigurations []WorldElementConfiguration
}

//PlayerConfiguration contains the required information to create a player.
type PlayerConfiguration struct {
	//inital position.
	InitialPosition *common.Point2D
	//initial angle
	InitialAngle float64
	//velocity
	Velocity float64
	//rotation velocity
	StepAngle float64
}

//WorldElementConfiguration contains the required information to create a world-element.
type WorldElementConfiguration struct {
	//inital position.
	InitialPosition *common.Point2D
	//initial angle
	InitialAngle float64
	//velocity
	Velocity float64
	//the size
	Size float64
	//the style
	Style tcell.Style
}

//Engine represents an game engine which takes input and render on screen.
type Engine interface {
	StartEngine()
}

//Impl implements the Engine interface.
type Impl struct {
	screen              tcell.Screen
	player              character.Character
	worldMap            world.WorldMap
	worldElements       []worldelement.WorldElement
	bgRender            render.Renderer
	consoleEventManager event.ConsoleEventManager
	quit                chan struct{}
	frameRate           int
	updateRate          int
}

//NewEngine provides a new engine.
func NewEngine(screen tcell.Screen, engineConfig *Configuration) (*Impl, error) {
	if engineConfig.WorldMap == nil {
		return nil, fmt.Errorf("world-map cannot be nil")
	}
	worldMap := world.NewWorldMap(engineConfig.WorldMap)
	if engineConfig.PlayerConfiguration == nil {
		return nil, fmt.Errorf("player-configuration cannot be nil")
	}
	player := character.NewPlayableCharacter(
		engineConfig.PlayerConfiguration.InitialPosition,
		engineConfig.PlayerConfiguration.InitialAngle,
		engineConfig.PlayerConfiguration.Velocity,
		engineConfig.PlayerConfiguration.StepAngle,
		worldMap)
	raySampler, err := render.CreateRaySamplerForAnsiColorTerminal(
		engineConfig.GradientRSFirst,
		engineConfig.GradientRSMultiplicator,
		engineConfig.GradientRSLimit,
		engineConfig.GradientRSWallStartColor,
		engineConfig.GradientRSWallEndColor,
		engineConfig.ScreenHeight,
		engineConfig.GradientRSBackgroundRange,
		engineConfig.GradientRSBackgroundColors)
	if err != nil {
		return nil, fmt.Errorf("error while instantiating the ray-sampler: %w", err)
	}
	mathHelper, err := common.NewMathHelper(new(common.RayCasterImpl))
	if err != nil {
		return nil, fmt.Errorf("error while instantiating the math-helper: %w", err)
	}
	renderMathHelper := render.NewRendererMathHelper(mathHelper)
	wallRendererProducer := render.CreateWallRendererProducer(
		engineConfig.ScreenWidth,
		engineConfig.ScreenHeight,
		engineConfig.PlayerFieldOfViewAngle,
		engineConfig.Visibility,
		mathHelper,
		renderMathHelper,
		tcell.StyleDefault.Background(tcell.ColorBlueViolet),
		raySampler)
	worldElementRendererProducer := render.CreateWorldElementRendererProducer(mathHelper, renderMathHelper, engineConfig.ScreenHeight, engineConfig.ScreenWidth)
	bgRender := render.CreateRenderer(engineConfig.ScreenWidth, engineConfig.ScreenHeight, renderMathHelper, engineConfig.PlayerFieldOfViewAngle, wallRendererProducer, worldElementRendererProducer)
	quit := make(chan struct{})
	consoleEventManager := event.NewConsoleEventManager(screen, player, quit)
	worldElements := make([]worldelement.WorldElement, len(engineConfig.WorldElementConfigurations))
	for i := 0; i < len(engineConfig.WorldElementConfigurations); i++ {
		worldElements[i] = worldelement.NewWorldElementImpl(
			engineConfig.WorldElementConfigurations[i].InitialPosition,
			engineConfig.WorldElementConfigurations[i].InitialAngle,
			engineConfig.WorldElementConfigurations[i].Velocity,
			engineConfig.WorldElementConfigurations[i].Size,
			engineConfig.WorldElementConfigurations[i].Style,
			worldMap,
			mathHelper,
		)
	}
	engine := Impl{
		screen:              screen,
		player:              player,
		worldMap:            worldMap,
		worldElements:       worldElements,
		bgRender:            bgRender,
		consoleEventManager: consoleEventManager,
		quit:                quit,
		frameRate:           engineConfig.FrameRate,
		updateRate:          engineConfig.WorlUpdateRate,
	}
	return &engine, nil
}

//StartEngine initializes the required element and start the engine to render world's elements in pseudo-3D
func (engine *Impl) StartEngine() {
	engine.screen.Clear()
	engine.player.Start()
	go engine.consoleEventManager.Listen()
	frameUpdateTicker := time.NewTicker(time.Duration(1000/engine.frameRate) * time.Millisecond)
	worldUpdateTicker := time.NewTicker(time.Duration(1000/engine.updateRate) * time.Millisecond)
	for _, worldelement := range engine.worldElements {
		worldelement.Start()
	}
	for {
		select {
		case <-engine.quit:
			frameUpdateTicker.Stop()
			worldUpdateTicker.Stop()
			engine.screen.SetStyle(tcell.StyleDefault)
			engine.screen.Clear()
			engine.screen.Fini()
			close(engine.player.GetQuitChannel())
			for _, worldelement := range engine.worldElements {
				close(worldelement.GetQuitChannel())
			}
			return
		case <-frameUpdateTicker.C:
			engine.bgRender.Render(engine.worldMap, engine.player, engine.worldElements, engine.screen)
		case updateWorldTickerTime := <-worldUpdateTicker.C:
			engine.player.GetUpdateChannel() <- updateWorldTickerTime
			for _, worldelement := range engine.worldElements {
				worldelement.GetUpdateChannel() <- updateWorldTickerTime
			}
		}
	}
}
