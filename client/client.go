package client

import (
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
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
	//the player's identifier
	PlayerID string
	//The world-map
	//WorldMap world.WorldMap
	//The player's configuration
	//PlayerConfiguration state.AnimatedElementState
	//the other-players.
	//OtherPlayerConfigurations map[string]state.AnimatedElementState
	// the quit-channel
	QuitChannel chan struct{}
}

//Engine represents an game engine which takes input and render on screen.
type Engine interface {
	GetPlayer() player.Player
	StartEngine()
	Initialize(playerID string, playerState state.AnimatedElementState, worldMap world.WorldMap, otherPlayers map[string]state.AnimatedElementState, serverTimeFramce uint32)
	ReceiveEventsFromServer(events []event.Event)
	GetShutdown() <-chan interface{}
}
