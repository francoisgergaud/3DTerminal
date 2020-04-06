package player

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/math"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	"testing"
	"time"

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
	player := NewPlayer(state, worldMap, mathHelper, nil)
	worldMap.On("GetCellValue", 2, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyUp, 0, 0))
	player.Move()
	assert.True(t, player.GetState().Position.AlmostEquals(&math.Point2D{X: 2, Y: 3}))
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
	player := NewPlayer(state, worldMap, mathHelper, nil)
	worldMap.On("GetCellValue", 2, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyUp, 0, 0))
	player.Move()
	assert.True(t, player.GetState().Position.AlmostEquals(&math.Point2D{X: 1, Y: 3}))
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
	player := NewPlayer(state, worldMap, mathHelper, nil)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyDown, 0, 0))
	player.Move()
	assert.True(t, player.GetState().Position.AlmostEquals(&math.Point2D{X: 0, Y: 3}))
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
	player := NewPlayer(state, worldMap, mathHelper, nil)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyDown, 0, 0))
	player.Move()
	assert.True(t, player.GetState().Position.AlmostEquals(&math.Point2D{X: 1, Y: 3}))
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
	player := NewPlayer(state, worldMap, mathHelper, nil)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyLeft, 0, 0))
	player.Move()
	assert.Equal(t, 1.99, player.GetState().Angle)
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
	player := NewPlayer(state, worldMap, mathHelper, nil)
	player.Action(tcell.NewEventKey(tcell.KeyLeft, 0, 0))
	player.Move()
	assert.Equal(t, 0.0, player.GetState().Angle)
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
	player := NewPlayer(state, worldMap, mathHelper, nil)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyRight, 0, 0))
	player.Move()
	assert.Equal(t, 0.01, player.GetState().Angle)
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
	player := NewPlayer(state, worldMap, mathHelper, nil)
	player.Action(tcell.NewEventKey(tcell.KeyRight, 0, 0))
	player.Move()
	assert.Equal(t, 0.0, player.GetState().Angle)
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
	player := NewPlayer(playerState, worldMap, mathHelper, nil)
	assert.Equal(t, position, player.GetState().Position)
	assert.Equal(t, angle, player.GetState().Angle)
	assert.Equal(t, stepAngle, player.GetState().StepAngle)
	assert.Equal(t, size, player.GetState().Size)
	assert.Equal(t, velocity, player.GetState().Velocity)
	assert.Equal(t, style, player.GetState().Style)
	assert.Equal(t, state.None, player.GetState().MoveDirection)
	assert.Equal(t, state.None, player.GetState().RotateDirection)
	assert.NotNil(t, player.GetUpdateChannel())
}

func TestPlayerStart(t *testing.T) {
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
	quit := make(chan struct{})
	player := NewPlayer(state, worldMap, mathHelper, quit)
	go func() {
		<-time.After(time.Millisecond * time.Duration(1000))
		close(quit)
	}()
	player.Start()
}
