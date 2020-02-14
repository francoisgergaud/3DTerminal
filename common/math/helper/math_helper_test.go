package helper

import (
	"francoisgergaud/3dGame/common/environment/world"
	innerMath "francoisgergaud/3dGame/common/math"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	"math"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

type MockRayCaster struct {
	mock.Mock
}

func (mock *MockRayCaster) CastRay(origin *innerMath.Point2D, world world.WorldMap, angle float64, maxDistance float64) *innerMath.Point2D {
	args := mock.Called(origin, world, angle, maxDistance)
	return args.Get(0).(*innerMath.Point2D)
}

var world2 testworld.MockWorldMapWithGrid = testworld.MockWorldMapWithGrid{
	Grid: [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	},
}

func TestNewMathHelperWithRayCaster(t *testing.T) {
	mathHelper, err := NewMathHelper(nil)
	assert.Nil(t, mathHelper)
	assert.NotNil(t, err)
}

func TestNormalizeAngle(t *testing.T) {
	mathHelper, _ := NewMathHelper(new(MockRayCaster))
	assert.True(t, AreFloatAlmostEquals(0.1, mathHelper.NormalizeAngle(2.1), 0.0001))
	assert.True(t, AreFloatAlmostEquals(1.9, mathHelper.NormalizeAngle(-0.1), 0.0001))
}

func TestCastRay(t *testing.T) {
	raycaster := new(MockRayCaster)
	mathHelper, _ := NewMathHelper(raycaster)
	origin := &innerMath.Point2D{}
	angle := 0.1
	visibility := 3.1
	raycaster.On("CastRay", origin, &world2, angle, visibility).Return(&innerMath.Point2D{})
	mathHelper.CastRay(origin, &world2, angle, visibility)
}

func TestGetWorldElementProjectionLeftInsideRightInsideFov(t *testing.T) {
	raycaster := new(MockRayCaster)
	renderMathHelper, _ := NewMathHelper(raycaster)
	playerPosition := &innerMath.Point2D{X: 0, Y: 0}
	viewAngle := 0.5
	fov := 0.5
	worlElementPosition := &innerMath.Point2D{X: 5, Y: 5}
	worldElementSize := 0.2
	isVisible, startColumnRatio, startOffset, endColumnRatio, endOffset := renderMathHelper.GetWorldElementProjection(playerPosition, viewAngle, fov, worlElementPosition, worldElementSize)
	assert.True(t, isVisible)
	assert.True(t, AreFloatAlmostEquals(0.0, startColumnRatio, 0.001), "wrong value for startColumnRatio: %v, expected: %v", startColumnRatio, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.5, startOffset, 0.001), "wrong value for startOffset: %v, expected: %v", startOffset, 0.5)
	assert.True(t, AreFloatAlmostEquals(0.0180, endColumnRatio, 0.001), "wrong value for endColumnRatio: %v, expected: %v", endColumnRatio, 0.0180)
	assert.True(t, AreFloatAlmostEquals(1.0, endOffset, 0.001), "wrong value for endOffset: %v, expected: %v", endOffset, 1.0)
}

func TestGetWorldElementProjectionLeftOutsideRightInsideFov(t *testing.T) {
	raycaster := new(MockRayCaster)
	renderMathHelper, _ := NewMathHelper(raycaster)
	playerPosition := &innerMath.Point2D{X: 5, Y: 5}
	viewAngle := 0.0
	fov := 0.5
	worlElementPosition := &innerMath.Point2D{X: 7, Y: 2.5}
	worldElementSize := 0.5
	isVisible, startColumnRatio, startOffset, endColumnRatio, endOffset := renderMathHelper.GetWorldElementProjection(playerPosition, viewAngle, fov, worlElementPosition, worldElementSize)
	assert.True(t, isVisible)
	assert.True(t, AreFloatAlmostEquals(0.0, startColumnRatio, 0.001), "wrong value for startColumnRatio: %v, expected: %v", startColumnRatio, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.8571, startOffset, 0.001), "wrong value for startOffset: %v, expected: %v", startOffset, 0.8571)
	assert.True(t, AreFloatAlmostEquals(0.0281, endColumnRatio, 0.001), "wrong value for endColumnRatio: %v, expected: %v", endColumnRatio, 0.0281)
	assert.True(t, AreFloatAlmostEquals(1.0, endOffset, 0.001), "wrong value for endOffset: %v, expected: %v", endOffset, 1.0)
}

func TestGetWorldElementProjectionLeftInsideRightOutsideFov(t *testing.T) {
	raycaster := new(MockRayCaster)
	renderMathHelper, _ := NewMathHelper(raycaster)
	playerPosition := &innerMath.Point2D{X: 5, Y: 5}
	viewAngle := 1.5
	fov := 0.5
	worlElementPosition := &innerMath.Point2D{X: 7, Y: 3.5}
	worldElementSize := 0.5
	isVisible, startColumnRatio, startOffset, endColumnRatio, endOffset := renderMathHelper.GetWorldElementProjection(playerPosition, viewAngle, fov, worlElementPosition, worldElementSize)
	assert.True(t, isVisible)
	assert.True(t, AreFloatAlmostEquals(0.9646, startColumnRatio, 0.001), "wrong value for startColumnRatio: %v, expected: %v", startColumnRatio, 0.9646)
	assert.True(t, AreFloatAlmostEquals(0.0, startOffset, 0.001), "wrong value for startOffset: %v, expected: %v", startOffset, 0.0)
	assert.True(t, AreFloatAlmostEquals(1.0, endColumnRatio, 0.001), "wrong value for endColumnRatio: %v, expected: %v", endColumnRatio, 1.0)
	assert.True(t, AreFloatAlmostEquals(0.1405, endOffset, 0.001), "wrong value for endOffset: %v, expected: %v", endOffset, 0.1405)
	viewAngle = 1.4
	isVisible, startColumnRatio, startOffset, endColumnRatio, endOffset = renderMathHelper.GetWorldElementProjection(playerPosition, viewAngle, fov, worlElementPosition, worldElementSize)
	assert.False(t, isVisible)
	assert.True(t, AreFloatAlmostEquals(0.0, startColumnRatio, 0.001), "wrong value for startColumnRatio: %v, expected: %v", startColumnRatio, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.0, startOffset, 0.001), "wrong value for startOffset: %v, expected: %v", startOffset, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.0, endColumnRatio, 0.001), "wrong value for endColumnRatio: %v, expected: %v", endColumnRatio, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.0, endOffset, 0.001), "wrong value for endOffset: %v, expected: %v", endOffset, 0.0)
}

func TestGetWorldElementProjectionLeftOutsideRightOutsideFov(t *testing.T) {
	raycaster := new(MockRayCaster)
	renderMathHelper, _ := NewMathHelper(raycaster)
	playerPosition := &innerMath.Point2D{X: 5, Y: 5}
	viewAngle := 1.5
	fov := 0.5
	worlElementPosition := &innerMath.Point2D{X: 5, Y: 4}
	worldElementSize := 1.2
	isVisible, startColumnRatio, startOffset, endColumnRatio, endOffset := renderMathHelper.GetWorldElementProjection(playerPosition, viewAngle, fov, worlElementPosition, worldElementSize)
	assert.True(t, isVisible)
	assert.True(t, AreFloatAlmostEquals(0.0, startColumnRatio, 0.001), "wrong value for startColumnRatio: %v, expected: %v", startColumnRatio, 0.695)
	assert.True(t, AreFloatAlmostEquals(0.051, startOffset, 0.001), "wrong value for startOffset: %v, expected: %v", startOffset, 0.051)
	assert.True(t, AreFloatAlmostEquals(1.0, endColumnRatio, 0.001), "wrong value for endColumnRatio: %v, expected: %v", endColumnRatio, 1.0)
	assert.True(t, AreFloatAlmostEquals(0.948, endOffset, 0.001), "wrong value for endOffset: %v, expected: %v", endOffset, 0.948)
}

func TestGetWorldElementProjectionBehindFov(t *testing.T) {
	raycaster := new(MockRayCaster)
	renderMathHelper, _ := NewMathHelper(raycaster)
	playerPosition := &innerMath.Point2D{X: 5, Y: 5}
	viewAngle := 1.0
	fov := 0.5
	worlElementPosition := &innerMath.Point2D{X: 6, Y: 4}
	worldElementSize := 1.2
	isVisible, startColumnRatio, startOffset, endColumnRatio, endOffset := renderMathHelper.GetWorldElementProjection(playerPosition, viewAngle, fov, worlElementPosition, worldElementSize)
	assert.False(t, isVisible)
	assert.True(t, AreFloatAlmostEquals(0.0, startColumnRatio, 0.001), "wrong value for startColumnRatio: %v, expected: %v", startColumnRatio, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.0, startOffset, 0.001), "wrong value for startOffset: %v, expected: %v", startOffset, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.0, endColumnRatio, 0.001), "wrong value for endColumnRatio: %v, expected: %v", endColumnRatio, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.0, endOffset, 0.001), "wrong value for endOffset: %v, expected: %v", endOffset, 0.0)
}

func TestGetWorldElementProjectionInsideFov(t *testing.T) {
	raycaster := new(MockRayCaster)
	renderMathHelper, _ := NewMathHelper(raycaster)
	playerPosition := &innerMath.Point2D{X: 0, Y: 5}
	viewAngle := 0.0
	fov := 0.5
	worlElementPosition := &innerMath.Point2D{X: 5, Y: 5}
	worldElementSize := 0.5
	isVisible, startColumnRatio, startOffset, endColumnRatio, endOffset := renderMathHelper.GetWorldElementProjection(playerPosition, viewAngle, fov, worlElementPosition, worldElementSize)
	assert.True(t, isVisible)
	assert.True(t, AreFloatAlmostEquals(0.4365, startColumnRatio, 0.001), "wrong value for startColumnRatio: %v, expected: %v", startColumnRatio, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.0, startOffset, 0.001), "wrong value for startOffset: %v, expected: %v", startOffset, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.5634, endColumnRatio, 0.001), "wrong value for endColumnRatio: %v, expected: %v", endColumnRatio, 0.0)
	assert.True(t, AreFloatAlmostEquals(1.0, endOffset, 0.001), "wrong value for endOffset: %v, expected: %v", endOffset, 0.0)
	//th object is closer, and still fully visible
	worlElementPosition2 := &innerMath.Point2D{X: 2, Y: 5}
	isVisible2, startColumnRatio2, startOffset2, endColumnRatio2, endOffset2 := renderMathHelper.GetWorldElementProjection(playerPosition, viewAngle, fov, worlElementPosition2, worldElementSize)
	assert.True(t, isVisible2)
	assert.True(t, AreFloatAlmostEquals(0.3440, startColumnRatio2, 0.001), "wrong value for startColumnRatio: %v, expected: %v", startColumnRatio, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.0, startOffset2, 0.001), "wrong value for startOffset: %v, expected: %v", startOffset, 0.0)
	assert.True(t, AreFloatAlmostEquals(0.656, endColumnRatio2, 0.001), "wrong value for endColumnRatio: %v, expected: %v", endColumnRatio, 0.0)
	assert.True(t, AreFloatAlmostEquals(1.0, endOffset2, 0.001), "wrong value for endOffset: %v, expected: %v", endOffset, 0.0)
}

func AreFloatAlmostEquals(f1, f2 float64, precision float64) bool {
	return math.Abs(f1-f2) < precision
}
