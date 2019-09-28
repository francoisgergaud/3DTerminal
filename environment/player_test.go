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
	player := NewPlayer(position, 0.0, 1.0, worldMap)
	worldMap.On("GetCellValue", 2, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyUp, 0, 0))
	assert.True(t, player.GetPosition().AlmostEquals(common.Point2D{X: 2, Y: 3}))
}

func TestActionKeyDown(t *testing.T) {
	worldMap := new(testutils.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	player := NewPlayer(position, 0.0, 1.0, worldMap)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyDown, 0, 0))
	assert.True(t, player.GetPosition().AlmostEquals(common.Point2D{X: 0, Y: 3}))
}

func TestActionKeyLeft(t *testing.T) {
	worldMap := new(testutils.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	player := NewPlayer(position, 0.0, 1.0, worldMap)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyLeft, 0, 0))
	assert.Equal(t, player.GetAngle(), 1.95)
}

func TestActionKeyRight(t *testing.T) {
	worldMap := new(testutils.MockWorldMap)
	position := &common.Point2D{X: 1, Y: 3}
	player := NewPlayer(position, 0.0, 1.0, worldMap)
	worldMap.On("GetCellValue", 0, 3).Return(0)
	player.Action(tcell.NewEventKey(tcell.KeyRight, 0, 0))
	assert.Equal(t, player.GetAngle(), 0.05)
}
