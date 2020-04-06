package player

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"

	"github.com/gdamore/tcell"

	animatedelementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
)

//NewPlayer creates a new player
func NewPlayer(worldMap world.WorldMap, mathHelper helper.MathHelper, quit <-chan interface{}) animatedelement.AnimatedElement {
	animatedElementState := state.AnimatedElementState{
		Position:  &math.Point2D{X: 5, Y: 5},
		Angle:     0.0,
		Size:      0.5,
		Velocity:  0.1,
		StepAngle: 0.01,
		Style:     tcell.StyleDefault.Background(tcell.Color126),
	}
	return animatedelementImpl.NewAnimatedElementWithState(&animatedElementState, worldMap, mathHelper)
}
