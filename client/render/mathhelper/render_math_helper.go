package mathhelper

import (
	"francoisgergaud/3dGame/common/math"
)

//RendererMathHelper provides maths for background render. It is an internal struct. Its purpose
//is to make the renderer's code more modular for testing purpose.
// This struct is stateless (no internal state) but usefull for dependency-injection and testing.
type RendererMathHelper interface {
	CalculateProjectionDistance(startPosition *math.Point2D, endPosition *math.Point2D, angle float64) float64
	IsWallAngle(point *math.Point2D) bool
	GetRayTracingAngleForColumn(playerAngle float64, columnIndex, screenWidth int, viewAngle float64) float64
	GetFillRowRange(distance, screenHeight float64) (int, int)
}
