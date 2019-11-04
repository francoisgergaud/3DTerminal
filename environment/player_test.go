package environment

import (
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/internal/testutils"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
)

func TestActionKeyUp(t *testing.T) {
	worldMap := new(testutils.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	player := &PlayableCharacter{
		pos:             position,
		angle:           0.0,
		velocity:        1.0,
		stepAngle:       0.01,
		world:           worldMap,
		moveDirection:   None,
		rotateDirection: None,
		updateChannel:   nil,
		quitChannel:     nil,
	}
	worldMap.On("GetCellValue", 2, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyUp, 0, 0))
	player.move()
	assert.True(t, player.GetPosition().AlmostEquals(common.Point2D{X: 2, Y: 3}))
}

func TestActionKeyDown(t *testing.T) {
	worldMap := new(testutils.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	player := &PlayableCharacter{
		pos:             position,
		angle:           0.0,
		velocity:        1.0,
		stepAngle:       0.01,
		world:           worldMap,
		moveDirection:   None,
		rotateDirection: None,
		updateChannel:   nil,
		quitChannel:     nil,
	}
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyDown, 0, 0))
	player.move()
	assert.True(t, player.GetPosition().AlmostEquals(common.Point2D{X: 0, Y: 3}))
}

func TestActionKeyLeft(t *testing.T) {
	worldMap := new(testutils.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	player := &PlayableCharacter{
		pos:             position,
		angle:           0.0,
		velocity:        1.0,
		stepAngle:       0.01,
		world:           worldMap,
		moveDirection:   None,
		rotateDirection: None,
		updateChannel:   nil,
		quitChannel:     nil,
	}
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyLeft, 0, 0))
	player.move()
	assert.Equal(t, 1.99, player.GetAngle())
}

func TestActionKeyRight(t *testing.T) {
	worldMap := new(testutils.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	player := &PlayableCharacter{
		pos:             position,
		angle:           0.0,
		velocity:        1.0,
		stepAngle:       0.01,
		world:           worldMap,
		moveDirection:   None,
		rotateDirection: None,
		updateChannel:   nil,
		quitChannel:     nil,
	}
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyRight, 0, 0))
	player.move()
	assert.Equal(t, 0.01, player.GetAngle())
}

func TestNewPlayableCharacter(t *testing.T) {
	worldMap := new(testutils.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	initialAngle := 0.4
	playerConfiguration := &PlayableCharacterConfiguration{
		Velocity:  1.0,
		StepAngle: 0.01,
	}
	character := NewPlayableCharacter(position, initialAngle, playerConfiguration, worldMap)
	player, ok := character.(*PlayableCharacter)
	if !ok {
		assert.Fail(t, "Character is not type of PlayableCharacter.")
	} else {
		assert.Equal(t, position, player.pos)
		assert.Equal(t, initialAngle, player.angle)
		assert.Equal(t, 0.01, player.stepAngle)
		assert.Equal(t, 1.0, player.velocity)
		assert.Equal(t, worldMap, player.world)
		assert.Equal(t, None, player.moveDirection)
		assert.Equal(t, None, player.rotateDirection)
	}
}
