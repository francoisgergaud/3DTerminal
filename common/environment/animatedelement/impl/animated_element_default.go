package impl

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	innerMath "francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"
	"math"

	"github.com/gdamore/tcell"
)

//NewAnimatedElement builds a new pointer to AnimatedElementImpl.
func NewAnimatedElement(initialPosition *innerMath.Point2D, initialAngle, velocity, stepAngle, size float64, moveDirection, rotateDirection state.Direction, style tcell.Style, world world.WorldMap, mathHelper helper.MathHelper) animatedelement.AnimatedElement {
	state := state.AnimatedElementState{
		Position:        initialPosition,
		Angle:           initialAngle,
		Velocity:        velocity,
		StepAngle:       stepAngle,
		Size:            size,
		Style:           style,
		MoveDirection:   moveDirection,
		RotateDirection: rotateDirection,
	}
	return NewAnimatedElementWithState(&state, world, mathHelper)
}

//NewAnimatedElementWithState builds a new pointer to AnimatedElementImpl.
func NewAnimatedElementWithState(animatedElementState *state.AnimatedElementState, world world.WorldMap, mathHelper helper.MathHelper) animatedelement.AnimatedElement {
	return &AnimatedElementImpl{
		state:      animatedElementState,
		world:      world,
		mathHelper: mathHelper,
	}
}

//AnimatedElementImpl is the implmenetation of AnimatedElement
type AnimatedElementImpl struct {
	state      *state.AnimatedElementState
	world      world.WorldMap
	mathHelper helper.MathHelper
	id         string
}

//State returns the animated-element's state.
func (animatedElement *AnimatedElementImpl) State() *state.AnimatedElementState {
	return animatedElement.state
}

//SetState update the animated-element's state.
func (animatedElement *AnimatedElementImpl) SetState(state *state.AnimatedElementState) {
	animatedElement.state = state
}

//Move updates the player's position depending on its moving and rotate Direction and the cell's value on the world-map
func (animatedElement *AnimatedElementImpl) Move() {
	if animatedElement.state.RotateDirection == state.Left {
		animatedElement.state.RotateDirection = state.Left
		animatedElement.state.Angle = animatedElement.state.Angle - animatedElement.state.StepAngle
		if animatedElement.state.Angle < 0 {
			animatedElement.state.Angle += 2
		}
	} else if animatedElement.state.RotateDirection == state.Right {
		animatedElement.state.Angle = animatedElement.state.Angle + animatedElement.state.StepAngle
		if animatedElement.state.Angle >= 2 {
			animatedElement.state.Angle -= 2
		}
	}
	if animatedElement.state.MoveDirection != state.None {
		newX := animatedElement.state.Position.X
		newY := animatedElement.state.Position.Y
		if animatedElement.state.MoveDirection == state.Forward {
			newX = animatedElement.state.Position.X + math.Cos(animatedElement.state.Angle*math.Pi)*animatedElement.state.Velocity
			newY = animatedElement.state.Position.Y + math.Sin(animatedElement.state.Angle*math.Pi)*animatedElement.state.Velocity
		} else if animatedElement.state.MoveDirection == state.Backward {
			newX = animatedElement.state.Position.X - math.Cos(animatedElement.state.Angle*math.Pi)*animatedElement.state.Velocity
			newY = animatedElement.state.Position.Y - math.Sin(animatedElement.state.Angle*math.Pi)*animatedElement.state.Velocity
		}
		if animatedElement.world.GetCellValue(int(newX), int(newY)) == 0 {
			animatedElement.state.Position.X = newX
			animatedElement.state.Position.Y = newY
		}
	}
}
