package impl

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/math"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testmathhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
)

func TestNewWorldElementImpl(t *testing.T) {
	worldElementID := "worldElementID"
	position := &math.Point2D{}
	angle := 1.5
	velocity := 1.3
	size := 0.6
	stepAngle := 0.02
	moveDirection := state.Forward
	rotateDirection := state.Left
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmathhelper.MockMathHelper)
	worldElement := NewBotImpl(worldElementID, position, angle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper, nil)
	assert.NotNil(t, worldElement)
	state := worldElement.State()
	assert.Equal(t, position, state.Position)
	assert.Equal(t, angle, state.Angle)
	assert.Equal(t, velocity, state.Velocity)
	assert.Equal(t, stepAngle, state.StepAngle)
	assert.Equal(t, size, state.Size)
	assert.Equal(t, moveDirection, state.MoveDirection)
	assert.Equal(t, rotateDirection, state.RotateDirection)
	assert.Equal(t, style, state.Style)
}

func TestMove(t *testing.T) {
	worldElementID := "worldElementID"
	position := &math.Point2D{X: 0, Y: 0}
	angle := 1.5
	velocity := 1.3
	size := 0.6
	stepAngle := 0.02
	moveDirection := state.Forward
	rotateDirection := state.None
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	worldMap.On("GetCellValue", 0, -1).Return(0)
	mathHelper := new(testmathhelper.MockMathHelper)
	worldElement := NewBotImpl(worldElementID, position, angle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper, nil)
	mathHelper.On("CastRay", position, worldMap, angle, velocity).Return(nil)
	worldElement.Move()
	assert.True(t, worldElement.State().Position.AlmostEquals(&math.Point2D{X: 0, Y: -1.3}))
}

func TestMoveWithRightRebound(t *testing.T) {
	worldElementID := "worldElementID"
	position := &math.Point2D{X: 0, Y: 0}
	angle := 0.0
	velocity := 1.5
	size := 0.5
	stepAngle := 0.02
	moveDirection := state.Forward
	rotateDirection := state.None
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	worldMap.On("GetCellValue", 1, 0).Return(1)
	mathHelper := new(testmathhelper.MockMathHelper)
	worldElement := NewBotImpl(worldElementID, position, angle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper, nil)
	mathHelper.On("CastRay", position, worldMap, angle, velocity).Return(&math.Point2D{X: 1.0, Y: 0.0})
	worldElement.Move()
	//assert.True(t, worldElement.GetState().Position.AlmostEquals(&common.Point2D{X: 0.5, Y: 0}))
	assert.True(t, worldElement.State().Position.AlmostEquals(&math.Point2D{X: 0, Y: 0}))
}

func TestMoveWithLeftRebound(t *testing.T) {
	worldElementID := "worldElementID"
	position := &math.Point2D{X: 0, Y: 0}
	angle := 1.0
	velocity := 1.5
	size := 0.5
	stepAngle := 0.02
	moveDirection := state.Forward
	rotateDirection := state.None
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmathhelper.MockMathHelper)
	worldElement := NewBotImpl(worldElementID, position, angle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper, nil)
	mathHelper.On("CastRay", position, worldMap, angle, velocity).Return(&math.Point2D{X: -1.0, Y: 0.0})
	worldElement.Move()
	//assert.True(t, worldElement.GetState().Position.AlmostEquals(&common.Point2D{X: -0.5, Y: 0}))
	assert.True(t, worldElement.State().Position.AlmostEquals(&math.Point2D{X: 0, Y: 0}))
}

func TestMoveWitTopRebound(t *testing.T) {
	worldElementID := "worldElementID"
	position := &math.Point2D{X: 0, Y: 0}
	angle := 1.5
	velocity := 1.5
	size := 0.5
	stepAngle := 0.02
	moveDirection := state.Forward
	rotateDirection := state.Left
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmathhelper.MockMathHelper)
	worldElement := NewBotImpl(worldElementID, position, angle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper, nil)
	mathHelper.On("CastRay", position, worldMap, angle, velocity).Return(&math.Point2D{X: 0.0, Y: -1.0})
	worldElement.Move()
	//assert.True(t, worldElement.GetState().Position.AlmostEquals(&common.Point2D{X: 0.0, Y: -0.5}))
	assert.True(t, worldElement.State().Position.AlmostEquals(&math.Point2D{X: 0, Y: 0}))
}

func TestMoveWitBottomRebound(t *testing.T) {
	worldElementID := "worldElementID"
	position := &math.Point2D{X: 0, Y: 0}
	angle := 0.5
	velocity := 1.5
	size := 0.5
	stepAngle := 0.02
	moveDirection := state.Forward
	rotateDirection := state.Left
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testmathhelper.MockMathHelper)
	worldElement := NewBotImpl(worldElementID, position, angle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper, nil)
	mathHelper.On("CastRay", position, worldMap, angle, velocity).Return(&math.Point2D{X: 0.0, Y: 1.0})
	worldElement.Move()
	//assert.True(t, worldElement.GetState().Position.AlmostEquals(&common.Point2D{X: 0.0, Y: 0.5}))
	assert.True(t, worldElement.State().Position.AlmostEquals(&math.Point2D{X: 0.0, Y: 0.0}))
}
