package render

import (
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/environment"
	"francoisgergaud/3dGame/internal/testutils"
	"math"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

type MockRayCaster struct {
	mock.Mock
}

func (mock *MockRayCaster) rayCast(origin *common.Point2D, world environment.WorldMap, angle float64, maxDistance float64) *common.Point2D {
	args := mock.Called(origin, world, angle, maxDistance)
	return args.Get(0).(*common.Point2D)
}

var world1 testutils.MockWorldMapWithGrid = testutils.MockWorldMapWithGrid{
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

func TestCalculateProjectionDistance(t *testing.T) {
	renderMathHelper := NewBackgroundRendererMathHelper(nil)
	distance1 := renderMathHelper.calculateProjectionDistance(&common.Point2D{X: 0, Y: 0}, &common.Point2D{X: 5, Y: 0}, 0.0)
	distance2 := renderMathHelper.calculateProjectionDistance(&common.Point2D{X: 0, Y: 0}, &common.Point2D{X: 5, Y: 5}, 0.25)
	float64EqualityThreshold := 0.001
	if math.Abs(distance1-distance2) > float64EqualityThreshold {
		t.Errorf("Projected distance incorrect, expected %v, got: %v.", distance1, distance2)
	}
}

// test if 3 consecutives rays casted at regular angle-step on a straight wall (no angle or empty wall)
// render 3 walls-slices with the same height difference between the 1st/2nd slice and 2nd/3rd slice.
func TestGetFillRowRange(t *testing.T) {
	raycaster := new(RayCasterImpl)
	renderMathHelper := NewBackgroundRendererMathHelper(raycaster)
	startPoint := common.Point2D{X: 1, Y: 2}
	visibility := 10.0
	screenHeight := 40.0
	startAngle := 1.7
	angleStep := 0.1
	impact1 := raycaster.rayCast(&startPoint, &world1, startAngle, visibility)
	distance1 := renderMathHelper.calculateProjectionDistance(&startPoint, impact1, 0)
	col1Start, _ := renderMathHelper.GetFillRowRange(distance1, screenHeight)
	impact2 := raycaster.rayCast(&startPoint, &world1, startAngle+angleStep, visibility)
	distance2 := renderMathHelper.calculateProjectionDistance(&startPoint, impact2, angleStep)
	col2Start, _ := renderMathHelper.GetFillRowRange(distance2, screenHeight)
	impact3 := raycaster.rayCast(&startPoint, &world1, startAngle+2*angleStep, visibility)
	distance3 := renderMathHelper.calculateProjectionDistance(&startPoint, impact3, 2*angleStep)
	col3Start, _ := renderMathHelper.GetFillRowRange(distance3, screenHeight)
	ratio1 := col1Start - col2Start
	ratio2 := col2Start - col3Start
	if ratio1-ratio2 < -1 || ratio1-ratio2 > 1 {
		t.Errorf("Incorrect distance projection ratio, expected %v, got: %v.", ratio1, ratio2)
	}
}

func TestIsWallAngle(t *testing.T) {
	renderMathHelper := new(BackgroundRendererMathHelperImpl)
	assert.True(t, renderMathHelper.isWallAngle(&common.Point2D{X: 0.01, Y: 0.05}))
}

func TestIsNotWallAngle(t *testing.T) {
	renderMathHelper := new(BackgroundRendererMathHelperImpl)
	assert.False(t, renderMathHelper.isWallAngle(&common.Point2D{X: 0.1, Y: 0.1}))
}

func TestCastRay(t *testing.T) {
	raycaster := new(MockRayCaster)
	renderMathHelper := NewBackgroundRendererMathHelper(raycaster)
	origin := &common.Point2D{}
	angle := 0.1
	visibility := 3.1
	raycaster.On("rayCast", origin, &world1, angle, visibility).Return(&common.Point2D{})
	renderMathHelper.castRay(origin, &world1, angle, visibility)
}

func TestGetRayTracingAngleForColumn(t *testing.T) {
	renderMathHelper := NewBackgroundRendererMathHelper(nil)
	assert.Equal(t, 1.76, renderMathHelper.getRayTracingAngleForColumn(0.0, 2, 100, 0.5))
	assert.Equal(t, 0.39, renderMathHelper.getRayTracingAngleForColumn(0.5, 28, 100, 0.5))
	assert.Equal(t, 0.125, renderMathHelper.getRayTracingAngleForColumn(1.875, 10, 10, 0.5))
}
