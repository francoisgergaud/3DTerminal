package character

import (
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/internal/testutils/environment/worldmap"
	"testing"
	"time"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
)

func TestActionKeyUp(t *testing.T) {
	worldMap := new(worldmap.MockWorldMap)
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

func TestActionKeyUpWhenMovingBackward(t *testing.T) {
	worldMap := new(worldmap.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	player := &PlayableCharacter{
		pos:             position,
		angle:           0.0,
		velocity:        1.0,
		stepAngle:       0.01,
		world:           worldMap,
		moveDirection:   Backward,
		rotateDirection: None,
		updateChannel:   nil,
		quitChannel:     nil,
	}
	worldMap.On("GetCellValue", 2, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyUp, 0, 0))
	player.move()
	assert.True(t, player.GetPosition().AlmostEquals(common.Point2D{X: 1, Y: 3}))
}

func TestActionKeyDown(t *testing.T) {
	worldMap := new(worldmap.MockWorldMap)
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

func TestActionKeyDownWhenMovingForward(t *testing.T) {
	worldMap := new(worldmap.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	player := &PlayableCharacter{
		pos:             position,
		angle:           0.0,
		velocity:        1.0,
		stepAngle:       0.01,
		world:           worldMap,
		moveDirection:   Forward,
		rotateDirection: None,
		updateChannel:   nil,
		quitChannel:     nil,
	}
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyDown, 0, 0))
	player.move()
	assert.True(t, player.GetPosition().AlmostEquals(common.Point2D{X: 1, Y: 3}))
}

func TestActionKeyLeft(t *testing.T) {
	worldMap := new(worldmap.MockWorldMap)
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

func TestActionKeyLeftWhenRotatingRight(t *testing.T) {
	worldMap := new(worldmap.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	intiialAngle := 0.0
	player := &PlayableCharacter{
		pos:             position,
		angle:           intiialAngle,
		velocity:        1.0,
		stepAngle:       0.01,
		world:           worldMap,
		moveDirection:   None,
		rotateDirection: Right,
		updateChannel:   nil,
		quitChannel:     nil,
	}
	player.Action(tcell.NewEventKey(tcell.KeyLeft, 0, 0))
	player.move()
	assert.Equal(t, intiialAngle, player.GetAngle())
}

func TestActionKeyRight(t *testing.T) {
	worldMap := new(worldmap.MockWorldMap)
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

func TestActionKeyRightWhenRotatingLeft(t *testing.T) {
	worldMap := new(worldmap.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	intiialAngle := 0.0
	player := &PlayableCharacter{
		pos:             position,
		angle:           intiialAngle,
		velocity:        1.0,
		stepAngle:       0.01,
		world:           worldMap,
		moveDirection:   None,
		rotateDirection: Left,
		updateChannel:   nil,
		quitChannel:     nil,
	}
	player.Action(tcell.NewEventKey(tcell.KeyRight, 0, 0))
	player.move()
	assert.Equal(t, intiialAngle, player.GetAngle())
}

func TestNewPlayableCharacter(t *testing.T) {
	worldMap := new(worldmap.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	initialAngle := 0.4
	velocity := 1.0
	stepAngle := 0.01
	character := NewPlayableCharacter(position, initialAngle, velocity, stepAngle, worldMap)
	player, ok := character.(*PlayableCharacter)
	if !ok {
		assert.Fail(t, "Character is not type of PlayableCharacter.")
	} else {
		assert.Equal(t, position, player.GetPosition())
		assert.Equal(t, initialAngle, player.GetAngle())
		assert.Equal(t, 0.01, player.stepAngle)
		assert.Equal(t, 1.0, player.velocity)
		assert.Equal(t, worldMap, player.world)
		assert.Equal(t, None, player.moveDirection)
		assert.Equal(t, None, player.rotateDirection)
		assert.NotNil(t, player.GetQuitChannel())
		assert.NotNil(t, player.GetUpdateChannel())
	}
}

func TestStart(t *testing.T) {
	worldMap := new(worldmap.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	intiialAngle := 0.0
	updateChannel := make(chan time.Time)
	quitChannel := make(chan struct{})
	player := &PlayableCharacter{
		pos:             position,
		angle:           intiialAngle,
		velocity:        1.0,
		stepAngle:       0.01,
		world:           worldMap,
		moveDirection:   None,
		rotateDirection: Left,
		updateChannel:   updateChannel,
		quitChannel:     quitChannel,
	}
	go func() {
		<-time.After(time.Millisecond * time.Duration(1000))
		close(quitChannel)
	}()
	player.Start()
}
