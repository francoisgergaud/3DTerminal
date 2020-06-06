package player

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/math"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
)

func TestActionKeyUp(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	state := state.AnimatedElementState{
		Position:        &math.Point2D{X: 1, Y: 3},
		Angle:           0.0,
		StepAngle:       0.01,
		Size:            0.5,
		Velocity:        1.0,
		Style:           tcell.StyleDefault.Background(tcell.Color111),
		MoveDirection:   state.None,
		RotateDirection: state.None,
	}
	player := NewPlayer(&state, worldMap, mathHelper)
	worldMap.On("GetCellValue", 2, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyUp, 0, 0))
	player.Move()
	assert.True(t, player.State().Position.AlmostEquals(&math.Point2D{X: 2, Y: 3}))
}

func TestActionKeyUpWhenMovingBackward(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	state := state.AnimatedElementState{
		Position:        &math.Point2D{X: 1, Y: 3},
		Angle:           0.0,
		StepAngle:       0.01,
		Size:            0.5,
		Velocity:        1.0,
		Style:           tcell.StyleDefault.Background(tcell.Color111),
		MoveDirection:   state.Backward,
		RotateDirection: state.None,
	}
	player := NewPlayer(&state, worldMap, mathHelper)
	worldMap.On("GetCellValue", 2, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyUp, 0, 0))
	player.Move()
	assert.True(t, player.State().Position.AlmostEquals(&math.Point2D{X: 1, Y: 3}))
}

func TestActionKeyDown(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	state := state.AnimatedElementState{
		Position:        &math.Point2D{X: 1, Y: 3},
		Angle:           0.0,
		StepAngle:       0.01,
		Size:            0.5,
		Velocity:        1.0,
		Style:           tcell.StyleDefault.Background(tcell.Color111),
		MoveDirection:   state.None,
		RotateDirection: state.None,
	}
	player := NewPlayer(&state, worldMap, mathHelper)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyDown, 0, 0))
	player.Move()
	assert.True(t, player.State().Position.AlmostEquals(&math.Point2D{X: 0, Y: 3}))
}

func TestActionKeyDownWhenMovingForward(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	state := state.AnimatedElementState{
		Position:        &math.Point2D{X: 1, Y: 3},
		Angle:           0.0,
		StepAngle:       0.01,
		Size:            0.5,
		Velocity:        1.0,
		Style:           tcell.StyleDefault.Background(tcell.Color111),
		MoveDirection:   state.Forward,
		RotateDirection: state.None,
	}
	player := NewPlayer(&state, worldMap, mathHelper)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyDown, 0, 0))
	player.Move()
	assert.True(t, player.State().Position.AlmostEquals(&math.Point2D{X: 1, Y: 3}))
}

func TestActionKeyLeft(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	state := state.AnimatedElementState{
		Position:        &math.Point2D{X: 1, Y: 3},
		Angle:           0.0,
		StepAngle:       0.01,
		Size:            0.5,
		Velocity:        1.0,
		Style:           tcell.StyleDefault.Background(tcell.Color111),
		MoveDirection:   state.None,
		RotateDirection: state.None,
	}
	player := NewPlayer(&state, worldMap, mathHelper)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyLeft, 0, 0))
	player.Move()
	assert.Equal(t, 1.99, player.State().Angle)
}

func TestActionKeyLeftWhenRotatingRight(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	state := state.AnimatedElementState{
		Position:        &math.Point2D{X: 1, Y: 3},
		Angle:           0.0,
		StepAngle:       0.01,
		Size:            0.5,
		Velocity:        1.0,
		Style:           tcell.StyleDefault.Background(tcell.Color111),
		MoveDirection:   state.None,
		RotateDirection: state.Right,
	}
	player := NewPlayer(&state, worldMap, mathHelper)
	player.Action(tcell.NewEventKey(tcell.KeyLeft, 0, 0))
	player.Move()
	assert.Equal(t, 0.0, player.State().Angle)
}

func TestActionKeyRight(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	state := state.AnimatedElementState{
		Position:        &math.Point2D{X: 1, Y: 3},
		Angle:           0.0,
		StepAngle:       0.01,
		Size:            0.5,
		Velocity:        1.0,
		Style:           tcell.StyleDefault.Background(tcell.Color111),
		MoveDirection:   state.None,
		RotateDirection: state.None,
	}
	player := NewPlayer(&state, worldMap, mathHelper)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyRight, 0, 0))
	player.Move()
	assert.Equal(t, 0.01, player.State().Angle)
}

func TestActionKeyRightWhenRotatingLeft(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	state := state.AnimatedElementState{
		Position:        &math.Point2D{X: 1, Y: 3},
		Angle:           0.0,
		StepAngle:       0.01,
		Size:            0.5,
		Velocity:        1.0,
		Style:           tcell.StyleDefault.Background(tcell.Color111),
		MoveDirection:   state.None,
		RotateDirection: state.Left,
	}
	player := NewPlayer(&state, worldMap, mathHelper)
	player.Action(tcell.NewEventKey(tcell.KeyRight, 0, 0))
	player.Move()
	assert.Equal(t, 0.0, player.State().Angle)
}

func TestNewPlayableCharacter(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	position := &math.Point2D{X: 1, Y: 3}
	angle := 0.0
	stepAngle := 0.01
	size := 0.5
	velocity := 1.0
	style := tcell.StyleDefault.Background(tcell.Color111)
	moveDirection := state.None
	rotateDirection := state.None
	playerState := state.AnimatedElementState{
		Position:        position,
		Angle:           angle,
		StepAngle:       stepAngle,
		Size:            size,
		Velocity:        velocity,
		Style:           style,
		MoveDirection:   moveDirection,
		RotateDirection: rotateDirection,
	}
	player := NewPlayer(&playerState, worldMap, mathHelper)
	playerCreatedState := player.State()
	assert.Equal(t, position, playerCreatedState.Position)
	assert.Equal(t, angle, playerCreatedState.Angle)
	assert.Equal(t, stepAngle, playerCreatedState.StepAngle)
	assert.Equal(t, size, playerCreatedState.Size)
	assert.Equal(t, velocity, playerCreatedState.Velocity)
	assert.Equal(t, style, playerCreatedState.Style)
	assert.Equal(t, state.None, playerCreatedState.MoveDirection)
	assert.Equal(t, state.None, playerCreatedState.RotateDirection)
}
