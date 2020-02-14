package state

import (
	"francoisgergaud/3dGame/common/math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gdamore/tcell"
)

func TestAnimatedElementStateClone(t *testing.T) {
	position := &math.Point2D{X: 0, Y: 0}
	state1 := &AnimatedElementState{
		Position:        position,
		Angle:           0.1,
		MoveDirection:   Forward,
		RotateDirection: Left,
		Size:            1.0,
		StepAngle:       0.01,
		Style:           tcell.StyleDefault,
		Velocity:        3.0,
	}
	state2 := state1.Clone()
	assert.True(t, state1.Position != state2.Position)
	assert.True(t, state1.Position.AlmostEquals(state2.Position))
	assert.Equal(t, state1.Angle, state2.Angle)
	assert.Equal(t, state1.MoveDirection, state2.MoveDirection)
	assert.Equal(t, state1.RotateDirection, state2.RotateDirection)
	assert.Equal(t, state1.Size, state2.Size)
	assert.Equal(t, state1.StepAngle, state2.StepAngle)
	assert.Equal(t, state1.Style, state2.Style)
	assert.Equal(t, state1.Velocity, state2.Velocity)
}
