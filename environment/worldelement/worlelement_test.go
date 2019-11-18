package worldelement

import (
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/internal/testutils"
	"francoisgergaud/3dGame/internal/testutils/environment/worldmap"
	"testing"
	"time"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
)

func TestNewWorldElementImpl(t *testing.T) {
	position := &common.Point2D{}
	angle := 1.5
	velocity := 1.3
	size := 0.6
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(worldmap.MockWorldMap)
	mathHelper := new(testutils.MockMathHelper)
	worldElement := NewWorldElementImpl(position, angle, velocity, size, style, worldMap, mathHelper)
	assert.NotNil(t, worldElement)
	assert.Equal(t, position, worldElement.GetPosition())
	assert.Equal(t, size, worldElement.GetSize())
	assert.Equal(t, style, worldElement.GetStyle())
	assert.NotNil(t, worldElement.GetQuitChannel())
	assert.NotNil(t, worldElement.GetUpdateChannel())
}

func TestStart(t *testing.T) {
	position := &common.Point2D{}
	angle := 1.5
	velocity := 1.3
	size := 0.6
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(worldmap.MockWorldMap)
	mathHelper := new(testutils.MockMathHelper)
	updateChannel := make(chan time.Time)
	quitChannel := make(chan struct{})
	worldElement := &WorldElementImpl{
		position:      position,
		size:          size,
		angle:         angle,
		velocity:      velocity,
		world:         worldMap,
		style:         style,
		updateChannel: updateChannel,
		quitChannel:   quitChannel,
		mathHelper:    mathHelper,
	}
	go func() {
		<-time.After(time.Millisecond * 1)
		close(quitChannel)
	}()
	worldElement.Start()
}

func TestMove(t *testing.T) {
	position := &common.Point2D{X: 0, Y: 0}
	angle := 1.5
	velocity := 1.3
	size := 0.6
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(worldmap.MockWorldMap)
	mathHelper := new(testutils.MockMathHelper)
	updateChannel := make(chan time.Time)
	quitChannel := make(chan struct{})
	worldElement := &WorldElementImpl{
		position:      position,
		size:          size,
		angle:         angle,
		velocity:      velocity,
		world:         worldMap,
		style:         style,
		updateChannel: updateChannel,
		quitChannel:   quitChannel,
		mathHelper:    mathHelper,
	}
	mathHelper.On("CastRay", position, worldMap, angle, velocity).Return(nil)
	worldElement.move()
	assert.True(t, worldElement.position.AlmostEquals(common.Point2D{X: 0, Y: -1.3}))
}

func TestMoveWithRightRebound(t *testing.T) {
	position := &common.Point2D{X: 0, Y: 0}
	angle := 0.0
	velocity := 1.5
	size := 0.5
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(worldmap.MockWorldMap)
	mathHelper := new(testutils.MockMathHelper)
	updateChannel := make(chan time.Time)
	quitChannel := make(chan struct{})
	worldElement := &WorldElementImpl{
		position:      position,
		size:          size,
		angle:         angle,
		velocity:      velocity,
		world:         worldMap,
		style:         style,
		updateChannel: updateChannel,
		quitChannel:   quitChannel,
		mathHelper:    mathHelper,
	}
	mathHelper.On("CastRay", position, worldMap, angle, velocity).Return(&common.Point2D{X: 1.0, Y: 0.0})
	worldElement.move()
	assert.True(t, worldElement.position.AlmostEquals(common.Point2D{X: 0.5, Y: 0}))
}

func TestMoveWithLeftRebound(t *testing.T) {
	position := &common.Point2D{X: 0, Y: 0}
	angle := 1.0
	velocity := 1.5
	size := 0.5
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(worldmap.MockWorldMap)
	mathHelper := new(testutils.MockMathHelper)
	updateChannel := make(chan time.Time)
	quitChannel := make(chan struct{})
	worldElement := &WorldElementImpl{
		position:      position,
		size:          size,
		angle:         angle,
		velocity:      velocity,
		world:         worldMap,
		style:         style,
		updateChannel: updateChannel,
		quitChannel:   quitChannel,
		mathHelper:    mathHelper,
	}
	mathHelper.On("CastRay", position, worldMap, angle, velocity).Return(&common.Point2D{X: -1.0, Y: 0.0})
	worldElement.move()
	assert.True(t, worldElement.position.AlmostEquals(common.Point2D{X: -0.5, Y: 0}))
}

func TestMoveWitTopRebound(t *testing.T) {
	position := &common.Point2D{X: 0, Y: 0}
	angle := 1.5
	velocity := 1.5
	size := 0.5
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(worldmap.MockWorldMap)
	mathHelper := new(testutils.MockMathHelper)
	updateChannel := make(chan time.Time)
	quitChannel := make(chan struct{})
	worldElement := &WorldElementImpl{
		position:      position,
		size:          size,
		angle:         angle,
		velocity:      velocity,
		world:         worldMap,
		style:         style,
		updateChannel: updateChannel,
		quitChannel:   quitChannel,
		mathHelper:    mathHelper,
	}
	mathHelper.On("CastRay", position, worldMap, angle, velocity).Return(&common.Point2D{X: 0.0, Y: -1.0})
	worldElement.move()
	assert.True(t, worldElement.position.AlmostEquals(common.Point2D{X: 0.0, Y: -0.5}))
}

func TestMoveWitBottomRebound(t *testing.T) {
	position := &common.Point2D{X: 0, Y: 0}
	angle := 0.5
	velocity := 1.5
	size := 0.5
	style := tcell.StyleDefault.Background(tcell.Color104)
	worldMap := new(worldmap.MockWorldMap)
	mathHelper := new(testutils.MockMathHelper)
	updateChannel := make(chan time.Time)
	quitChannel := make(chan struct{})
	worldElement := &WorldElementImpl{
		position:      position,
		size:          size,
		angle:         angle,
		velocity:      velocity,
		world:         worldMap,
		style:         style,
		updateChannel: updateChannel,
		quitChannel:   quitChannel,
		mathHelper:    mathHelper,
	}
	mathHelper.On("CastRay", position, worldMap, angle, velocity).Return(&common.Point2D{X: 0.0, Y: 1.0})
	worldElement.move()
	assert.True(t, worldElement.position.AlmostEquals(common.Point2D{X: 0.0, Y: 0.5}))
}
