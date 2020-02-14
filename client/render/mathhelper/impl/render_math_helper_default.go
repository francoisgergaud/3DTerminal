package impl

import (
	renderMathHelper "francoisgergaud/3dGame/client/render/mathhelper"
	internalMath "francoisgergaud/3dGame/common/math"
	internalMathHelper "francoisgergaud/3dGame/common/math/helper"
	"math"
)

//RendererMathHelperImpl implements the RendererMathHelper interface.
type RendererMathHelperImpl struct {
	mathHelper internalMathHelper.MathHelper
}

//NewRendererMathHelper builds a new BackgroundRendererMathHelper.
func NewRendererMathHelper(mathHelper internalMathHelper.MathHelper) renderMathHelper.RendererMathHelper {
	return &RendererMathHelperImpl{
		mathHelper: mathHelper,
	}
}

//CalculateProjectionDistance returns the distance of the projection from the destination to the projected
// point of the distance using the angle. The distance is not taken as it to avoid the "fish-eye" effect.
func (rendererMathHelper *RendererMathHelperImpl) CalculateProjectionDistance(startPosition *internalMath.Point2D, endPosition *internalMath.Point2D, angle float64) float64 {
	cosAngle := math.Cos(angle * math.Pi)
	return endPosition.Distance(startPosition) * cosAngle
}

//IsWallAngle checks if a point is close enough to a grid-point (the world-map is using a grid where the wall are
//using a whole cell). It returns a true if the input-point is close enough to be considered as a wall-edge.
func (rendererMathHelper *RendererMathHelperImpl) IsWallAngle(point *internalMath.Point2D) bool {
	distanceToWallAngle := math.Hypot(point.X-float64(math.Round(point.X)), point.Y-float64(math.Round(point.Y)))
	if distanceToWallAngle < 0.1 {
		return true
	}
	return false
}

//GetRayTracingAngleForColumn returns the ray-tracing's angle from an user position, a column on the screen to be renderer and the player´s view-angle.
func (rendererMathHelper *RendererMathHelperImpl) GetRayTracingAngleForColumn(playerAngle float64, columnIndex, screenWidth int, viewAngle float64) float64 {
	stepAngle := viewAngle / float64(screenWidth)
	rayTracingAngleToPlayer := -viewAngle/2 + (stepAngle * float64(columnIndex))
	return rendererMathHelper.mathHelper.NormalizeAngle(playerAngle + rayTracingAngleToPlayer)
}

//GetFillRowRange returns the start and end rows for a given obstable distance
func (rendererMathHelper *RendererMathHelperImpl) GetFillRowRange(distance, screenHeight float64) (int, int) {
	//if distance = verticalFieldOfView, startRow = 0, endRow = screenHeight
	//if distance = visibility, startRow=(screenHeight/2)-1, endRow=(screenHeight/2)+1
	verticalFieldOfView := 1.0
	if distance < verticalFieldOfView {
		distance = verticalFieldOfView
	}
	startRow := int(screenHeight/2.0 - screenHeight/(2.0*distance))
	endRow := int(screenHeight) - startRow
	return startRow, endRow
}