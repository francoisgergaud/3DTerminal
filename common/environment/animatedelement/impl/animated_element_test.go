package impl

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	innerMath "francoisgergaud/3dGame/common/math"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testmath "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	"math"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
)

func TestNewWorldElementImpl(t *testing.T) {
	position := &innerMath.Point2D{}
	initialAngle := 1.5
	velocity := 1.3
	size := 0.6
	stepAngle := 0.02
	moveDirection := state.Forward
	rotateDirection := state.Left
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmath.MockMathHelper)
	animatedElement := NewAnimatedElement(position, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper)
	assert.NotNil(t, animatedElement)
	state := animatedElement.State()
	assert.Equal(t, position, state.Position)
	assert.Equal(t, initialAngle, state.Angle)
	assert.Equal(t, velocity, state.Velocity)
	assert.Equal(t, stepAngle, state.StepAngle)
	assert.Equal(t, size, state.Size)
	assert.Equal(t, moveDirection, state.MoveDirection)
	assert.Equal(t, rotateDirection, state.RotateDirection)
	assert.Equal(t, style, state.Style)
}

func TestMoveLeft(t *testing.T) {
	position := &innerMath.Point2D{X: 1, Y: 1}
	initialAngle := 1.5
	velocity := 1.3
	size := 0.6
	stepAngle := 0.1
	moveDirection := state.None
	rotateDirection := state.Left
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmath.MockMathHelper)
	animatedElement := NewAnimatedElement(position, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper)
	animatedElement.Move()
	assert.Equal(t, 1.4, animatedElement.State().Angle)
	assert.True(t, innerMath.Point2D{X: 1, Y: 1}.AlmostEquals(animatedElement.State().Position))
}

func TestMoveLeftAngleReset(t *testing.T) {
	position := &innerMath.Point2D{X: 1, Y: 1}
	initialAngle := 0.1
	velocity := 1.3
	size := 0.6
	stepAngle := 0.2
	moveDirection := state.None
	rotateDirection := state.Left
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmath.MockMathHelper)
	animatedElement := NewAnimatedElement(position, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper)
	animatedElement.Move()
	assert.Equal(t, 1.9, animatedElement.State().Angle)
	assert.True(t, innerMath.Point2D{X: 1, Y: 1}.AlmostEquals(animatedElement.State().Position))
}

func TestMoveRight(t *testing.T) {
	position := &innerMath.Point2D{X: 1, Y: 1}
	initialAngle := 1.5
	velocity := 1.3
	size := 0.6
	stepAngle := 0.1
	moveDirection := state.None
	rotateDirection := state.Right
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmath.MockMathHelper)
	animatedElement := NewAnimatedElement(position, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper)
	animatedElement.Move()
	assert.Equal(t, 1.6, animatedElement.State().Angle)
	assert.True(t, innerMath.Point2D{X: 1, Y: 1}.AlmostEquals(animatedElement.State().Position))
}

func TestMoveRightAngleReset(t *testing.T) {
	position := &innerMath.Point2D{X: 1, Y: 1}
	initialAngle := 1.9
	velocity := 1.3
	size := 0.6
	stepAngle := 0.2
	moveDirection := state.None
	rotateDirection := state.Right
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmath.MockMathHelper)
	animatedElement := NewAnimatedElement(position, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper)
	animatedElement.Move()
	assert.True(t, AreFloatAlmostEquals(0.1, animatedElement.State().Angle, 0.001))
	assert.True(t, innerMath.Point2D{X: 1, Y: 1}.AlmostEquals(animatedElement.State().Position))
}

func AreFloatAlmostEquals(f1, f2 float64, precision float64) bool {
	return math.Abs(f1-f2) < precision
}

func TestMoveForward(t *testing.T) {
	position := &innerMath.Point2D{X: 1, Y: 1}
	initialAngle := 0.0
	velocity := 0.1
	size := 0.6
	stepAngle := 0.1
	moveDirection := state.Forward
	rotateDirection := state.None
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmath.MockMathHelper)
	worldMap.On("GetCellValue", 1, 1).Return(0)
	animatedElement := NewAnimatedElement(position, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper)
	animatedElement.Move()
	assert.True(t, innerMath.Point2D{X: 1.1, Y: 1}.AlmostEquals(animatedElement.State().Position))
}

func TestMoveBackward(t *testing.T) {
	position := &innerMath.Point2D{X: 1, Y: 1}
	initialAngle := 0.0
	velocity := 0.1
	size := 0.6
	stepAngle := 0.1
	moveDirection := state.Backward
	rotateDirection := state.None
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmath.MockMathHelper)
	worldMap.On("GetCellValue", 0, 1).Return(0)
	animatedElement := NewAnimatedElement(position, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper)
	animatedElement.Move()
	assert.True(t, innerMath.Point2D{X: 0.9, Y: 1}.AlmostEquals(animatedElement.State().Position))
}

func TestMoveForwardWithWall(t *testing.T) {
	position := &innerMath.Point2D{X: 1, Y: 1}
	initialAngle := 0.0
	velocity := 0.1
	size := 0.6
	stepAngle := 0.1
	moveDirection := state.Forward
	rotateDirection := state.None
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmath.MockMathHelper)
	worldMap.On("GetCellValue", 1, 1).Return(1)
	animatedElement := NewAnimatedElement(position, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper)
	animatedElement.Move()
	assert.True(t, innerMath.Point2D{X: 1, Y: 1}.AlmostEquals(animatedElement.State().Position))
}

func TestSetState(t *testing.T) {
	newState := &state.AnimatedElementState{}
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmath.MockMathHelper)
	worldMap.On("GetCellValue", 1, 1).Return(1)
	animatedElement := NewAnimatedElementWithState(&state.AnimatedElementState{}, worldMap, mathHelper)
	animatedElement.SetState(newState)
	assert.Equal(t, newState, animatedElement.State())
}
