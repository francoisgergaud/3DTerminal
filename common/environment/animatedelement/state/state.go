package state

import (
	"francoisgergaud/3dGame/common/math"

	"github.com/gdamore/tcell"
)

//AnimatedElementState provides a base implmentation for AnimatedElement.
//TODO: create interface and mocks for tests
type AnimatedElementState struct {
	Position        *math.Point2D
	Angle           float64
	StepAngle       float64
	Size            float64
	Velocity        float64
	Style           tcell.Style
	MoveDirection   Direction
	RotateDirection Direction
}

//Clone creates a copy.
func (a *AnimatedElementState) Clone() *AnimatedElementState {
	var position *math.Point2D
	if a.Position != nil {
		position = a.Position.Clone()
	}
	return &AnimatedElementState{
		Position:        position,
		Angle:           a.Angle,
		StepAngle:       a.StepAngle,
		Size:            a.Size,
		Velocity:        a.Velocity,
		Style:           a.Style,
		MoveDirection:   a.MoveDirection,
		RotateDirection: a.RotateDirection,
	}
}

//Direction is the direction type.
type Direction uint

//The Direction possible values.
const (
	None Direction = iota
	Left
	Right
	Forward
	Backward
)
