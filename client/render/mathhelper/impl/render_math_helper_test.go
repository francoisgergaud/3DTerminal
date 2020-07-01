package impl

import (
	internalMath "francoisgergaud/3dGame/common/math"
	raycaster "francoisgergaud/3dGame/common/math/raycaster"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testmathhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateProjectionDistance(t *testing.T) {
	renderMathHelper := NewRendererMathHelper(nil)
	distance1 := renderMathHelper.CalculateProjectionDistance(&internalMath.Point2D{X: 0, Y: 0}, &internalMath.Point2D{X: 5, Y: 0}, 0.0)
	distance2 := renderMathHelper.CalculateProjectionDistance(&internalMath.Point2D{X: 0, Y: 0}, &internalMath.Point2D{X: 5, Y: 5}, 0.25)
	float64EqualityThreshold := 0.001
	if math.Abs(distance1-distance2) > float64EqualityThreshold {
		t.Errorf("Projected distance incorrect, expected %v, got: %v.", distance1, distance2)
	}
}

// test if 3 consecutives rays casted at regular angle-step on a straight wall (no angle or empty wall)
// render 3 walls-slices with the same height difference between the 1st/2nd slice and 2nd/3rd slice.
func TestGetFillRowRange(t *testing.T) {
	var world1 testworld.MockWorldMapWithGrid = testworld.MockWorldMapWithGrid{
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
	for x := 0; x < 10; x++ {
		world1.On("GetCellValue", x, 0)
		world1.On("GetCellValue", x, 1)
	}
	raycaster := new(raycaster.RayCasterImpl)
	renderMathHelper := NewRendererMathHelper(nil)
	startPoint := internalMath.Point2D{X: 1, Y: 2}
	visibility := 10.0
	screenHeight := 40
	startAngle := 1.7
	angleStep := 0.1
	height := 1.0
	impact1 := raycaster.CastRay(&startPoint, &world1, startAngle, visibility)
	distance1 := renderMathHelper.CalculateProjectionDistance(&startPoint, impact1, 0)
	col1Start, _ := renderMathHelper.GetFillRowRange(distance1, height, visibility, screenHeight)
	impact2 := raycaster.CastRay(&startPoint, &world1, startAngle+angleStep, visibility)
	distance2 := renderMathHelper.CalculateProjectionDistance(&startPoint, impact2, angleStep)
	col2Start, _ := renderMathHelper.GetFillRowRange(distance2, height, visibility, screenHeight)
	impact3 := raycaster.CastRay(&startPoint, &world1, startAngle+2*angleStep, visibility)
	distance3 := renderMathHelper.CalculateProjectionDistance(&startPoint, impact3, 2*angleStep)
	col3Start, _ := renderMathHelper.GetFillRowRange(distance3, height, visibility, screenHeight)
	ratio1 := col1Start - col2Start
	ratio2 := col2Start - col3Start
	if ratio1-ratio2 < -1 || ratio1-ratio2 > 1 {
		t.Errorf("Incorrect distance projection ratio, expected %v, got: %v.", ratio1, ratio2)
	}
}

func TestIsWallAngle(t *testing.T) {
	renderMathHelper := NewRendererMathHelper(nil)
	assert.True(t, renderMathHelper.IsWallAngle(&internalMath.Point2D{X: 0.01, Y: 0.05}))
}

func TestIsNotWallAngle(t *testing.T) {
	renderMathHelper := NewRendererMathHelper(nil)
	assert.False(t, renderMathHelper.IsWallAngle(&internalMath.Point2D{X: 0.1, Y: 0.1}))
}

func TestGetRayTracingAngleForColumn(t *testing.T) {
	mathHelper := new(testmathhelper.MockMathHelper)
	renderMathHelper := NewRendererMathHelper(mathHelper)
	mathHelper.On("NormalizeAngle", -0.24).Return(1.76)
	assert.Equal(t, 1.76, renderMathHelper.GetRayTracingAngleForColumn(0.0, 2, 100, 0.5))
	mathHelper.On("NormalizeAngle", 0.39).Return(0.39)
	assert.Equal(t, 0.39, renderMathHelper.GetRayTracingAngleForColumn(0.5, 28, 100, 0.5))
	mathHelper.On("NormalizeAngle", 2.125).Return(0.125)
	assert.Equal(t, 0.125, renderMathHelper.GetRayTracingAngleForColumn(1.875, 10, 10, 0.5))
	//this case happen during the running: float-imprecision lead a x%2 == 2
	mathHelper.On("NormalizeAngle", -8.673617379884035e-18).Return(2.0)
	assert.Equal(t, 2.0, renderMathHelper.GetRayTracingAngleForColumn(0.01, 57, 120, 0.4))
}
