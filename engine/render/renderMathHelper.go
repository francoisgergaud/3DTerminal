package render

import (
	"francoisgergaud/3dGame/common"
	"math"
)

//RendererMathHelper provides maths for background render. It is an internal struct. Its purpose
//is to make the renderer's code more modular for testing purpose.
// This struct is stateless (no internal state) but usefull for dependency-injection and testing.
type RendererMathHelper interface {
	calculateProjectionDistance(startPosition *common.Point2D, endPosition *common.Point2D, angle float64) float64
	isWallAngle(point *common.Point2D) bool
	getRayTracingAngleForColumn(playerAngle float64, columnIndex, screenWidth int, viewAngle float64) float64
	getFillRowRange(distance, screenHeight float64) (int, int)
}

//RendererMathHelperImpl implements the RendererMathHelper interface.
type RendererMathHelperImpl struct {
	mathHelper common.MathHelper
}

//NewRendererMathHelper builds a new BackgroundRendererMathHelper.
func NewRendererMathHelper(mathHelper common.MathHelper) RendererMathHelper {
	return &RendererMathHelperImpl{
		mathHelper: mathHelper,
	}
}

// calculateProjectionDistance returns the distance of the projection from the destination to the projected
// point of the distance using the angle. The distance is not taken as it to avoid the "fish-eye" effect.
func (rendererMathHelper *RendererMathHelperImpl) calculateProjectionDistance(startPosition *common.Point2D, endPosition *common.Point2D, angle float64) float64 {
	cosAngle := math.Cos(angle * math.Pi)
	return endPosition.Distance(startPosition) * cosAngle
}

//isWallAngle checks if a point is close enough to a grid-point (the world-map is using a grid where the wall are
//using a whole cell). It returns a true if the input-point is close enough to be considered as a wall-edge.
func (rendererMathHelper *RendererMathHelperImpl) isWallAngle(point *common.Point2D) bool {
	distanceToWallAngle := math.Hypot(point.X-float64(math.Round(point.X)), point.Y-float64(math.Round(point.Y)))
	if distanceToWallAngle < 0.1 {
		return true
	}
	return false
}

//getRayTracingAngleForColumn returns the ray-tracing's angle from an user position, a column on the screen to be renderer and the player´s view-angle.
func (rendererMathHelper *RendererMathHelperImpl) getRayTracingAngleForColumn(playerAngle float64, columnIndex, screenWidth int, viewAngle float64) float64 {
	stepAngle := viewAngle / float64(screenWidth)
	rayTracingAngleToPlayer := -viewAngle/2 + (stepAngle * float64(columnIndex))
	return rendererMathHelper.mathHelper.NormalizeAngle(playerAngle + rayTracingAngleToPlayer)
}

//GetFillRowRange returns the start and end rows for a given obstable distance
func (rendererMathHelper *RendererMathHelperImpl) getFillRowRange(distance, screenHeight float64) (int, int) {
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
