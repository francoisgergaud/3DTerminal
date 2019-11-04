package main

import (
	"fmt"
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/engine"
	"francoisgergaud/3dGame/environment"
	"os"

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

// NewScreen provides a new screen.
func NewScreen() *tcell.Screen {
	screen, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = screen.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	return &screen
}

//NewWorldMap provides a new world-map.
func NewWorldMap() environment.WorldMap {
	grid := [][]int{
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
	return environment.NewWorldMap(grid)
}

//NewPlayer provides a new player.
func NewPlayer(worldMap environment.WorldMap) environment.Character {
	player := environment.NewPlayableCharacter(
		&common.Point2D{X: 5, Y: 5},
		0,
		&environment.PlayableCharacterConfiguration{
			Velocity:  0.1,
			StepAngle: 0.01,
		},
		worldMap)
	return player
}

//InitGame initializes a game.
func InitGame() *Game {
	screen := NewScreen()
	worldMap := NewWorldMap()
	player := NewPlayer(worldMap)
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
		GradientRSBackgroundRange:  []float32{0.5, 0.55, 0.65, 1.0},
		GradientRSBackgroundColors: []int{63, 58, 64, 70},
	}
	return &Game{
		engine: engine.NewEngine(*screen, player, worldMap, engineConfiguration),
	}

}
