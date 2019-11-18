package main

import (
	"fmt"
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/engine"

	"github.com/gdamore/tcell"
)

//Game defines the game entity.
type Game struct {
	engine engine.Engine
}

//Start the game.
func (game *Game) Start() {
	game.engine.StartEngine()
}

//NewWorldMap provides a new world-map.
func NewWorldMap() [][]int {
	return [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
}

//InitGame initializes a game.
func InitGame(screen tcell.Screen) (*Game, error) {
	worldMap := NewWorldMap()
	playerConfiguration := engine.PlayerConfiguration{
		InitialPosition: &common.Point2D{X: 5, Y: 5},
		InitialAngle:    0.0,
		Velocity:        0.1,
		StepAngle:       0.01,
	}
	worlElementConfigurations := make([]engine.WorldElementConfiguration, 1)
	worlElementConfigurations[0] = engine.WorldElementConfiguration{
		InitialPosition: &common.Point2D{X: 9, Y: 12},
		InitialAngle:    0.3,
		Velocity:        0.02,
		Size:            0.3,
		Style:           tcell.StyleDefault.Background(tcell.ColorDarkBlue),
	}
	engineConfiguration := &engine.Configuration{
		FrameRate:                  20,
		WorlUpdateRate:             40,
		ScreenHeight:               40,
		ScreenWidth:                120,
		PlayerFieldOfViewAngle:     0.4,
		Visibility:                 20.0,
		GradientRSFirst:            1.0,
		GradientRSMultiplicator:    2.0,
		GradientRSLimit:            10.0,
		GradientRSWallStartColor:   255,
		GradientRSWallEndColor:     240,
		GradientRSBackgroundRange:  []float32{0.5, 0.55, 0.65},
		GradientRSBackgroundColors: []int{63, 58, 64, 70},
		WorldMap:                   worldMap,
		PlayerConfiguration:        &playerConfiguration,
		WorldElementConfigurations: worlElementConfigurations,
	}
	engine, err := engine.NewEngine(screen, engineConfiguration)
	if err != nil {
		return nil, fmt.Errorf("Error while initializing engine: %w", err)
	}
	return &Game{
		engine: engine,
	}, nil

}
