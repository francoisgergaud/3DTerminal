package raycaster

import (
	"francoisgergaud/3dGame/common/environment/world"
	innerMath "francoisgergaud/3dGame/common/math"
	"math"
)

//RayCaster provides the function to cast a ray.
type RayCaster interface {
	CastRay(origin *innerMath.Point2D, world world.WorldMap, angle float64, maxDistance float64) *innerMath.Point2D
}

//RayCasterImpl implements the RayCast interface.
type RayCasterImpl struct {
}

//CastRay casts a ray from the origin, with a given angle, on the worldmap, until a wall is found, or the ray's length is greater than visibility.
func (raycaster *RayCasterImpl) CastRay(origin *innerMath.Point2D, world world.WorldMap, angle float64, maxDistance float64) *innerMath.Point2D {
	//Digital Diferential Analyzer
	// get the point's coordinate of the ray and the first obstacle on the map.
	// on vertical-intersection with the grid: verticalIntersectStep.Y is how much is increased Y everytime verticalIntersectStep.X is incremented/decremented by 1
	// on horizontal-intersection with the grid, horizontalIntersectStep.X is how much is increased X everytime horizontalIntersectStep.Y is incremented/decremented by 1
	var verticalIntersectStep, horizontalIntersectStep innerMath.Point2D
	// calculate the ray's steps from its angle
	switch {
	case angle < 0.25:
		verticalIntersectStep.Y = math.Tan(angle * math.Pi)
		verticalIntersectStep.X = 1
		horizontalIntersectStep.X = 1 / verticalIntersectStep.Y
		horizontalIntersectStep.Y = 1
	case angle < 0.75:
		horizontalIntersectStep.X = math.Tan((0.5 - angle) * math.Pi)
		horizontalIntersectStep.Y = 1
		if angle < 0.5 {
			verticalIntersectStep.X = 1
			verticalIntersectStep.Y = 1 / horizontalIntersectStep.X
		} else {
			verticalIntersectStep.X = -1
			verticalIntersectStep.Y = -1 / horizontalIntersectStep.X
		}
	case angle < 1.25:
		verticalIntersectStep.Y = math.Tan((1 - angle) * math.Pi)
		verticalIntersectStep.X = -1
		if angle < 1 {
			horizontalIntersectStep.Y = 1
			horizontalIntersectStep.X = -1 / verticalIntersectStep.Y
		} else {
			horizontalIntersectStep.Y = -1
			horizontalIntersectStep.X = 1 / verticalIntersectStep.Y
		}
	case angle < 1.75:
		horizontalIntersectStep.X = math.Tan((angle - 1.5) * math.Pi)
		horizontalIntersectStep.Y = -1
		if angle < 1.5 {
			verticalIntersectStep.X = -1
			verticalIntersectStep.Y = 1 / horizontalIntersectStep.X
		} else {
			verticalIntersectStep.X = 1
			verticalIntersectStep.Y = -1 / horizontalIntersectStep.X
		}
	case angle <= 2:
		verticalIntersectStep.Y = math.Tan(angle * math.Pi)
		verticalIntersectStep.X = 1
		horizontalIntersectStep.X = -1 / verticalIntersectStep.Y
		horizontalIntersectStep.Y = -1
	}
	//calculate the first horizontal an vertical intersections
	//TODO: double-check what happens when origin is negative
	var verticalIntersect, horizontalIntersect innerMath.Point2D
	verticalIntersect.X = float64(int(origin.X))
	if verticalIntersectStep.X > 0 {
		verticalIntersect.X++
	}
	verticalIntersect.Y = ((verticalIntersect.X - origin.X) * math.Tan(angle*math.Pi)) + origin.Y

	horizontalIntersect.Y = float64(int(origin.Y))
	if verticalIntersectStep.Y > 0 {
		horizontalIntersect.Y++
	}
	horizontalIntersect.X = ((horizontalIntersect.Y - origin.Y) / math.Tan(angle*math.Pi)) + origin.X

	var rayLength float64
	var result *innerMath.Point2D
	for rayLength < maxDistance {
		verticalIntersectDistance := origin.Distance(&verticalIntersect)
		horizontalIntersectDistance := origin.Distance(&horizontalIntersect)
		if verticalIntersectDistance > horizontalIntersectDistance {
			if raycaster.checkHorizontalCollision(world, &horizontalIntersect, horizontalIntersectStep.Y) {
				result = &horizontalIntersect
				break
			} else {
				horizontalIntersect.X += horizontalIntersectStep.X
				horizontalIntersect.Y += horizontalIntersectStep.Y
				rayLength = horizontalIntersectDistance
			}
		} else {
			if raycaster.checkVerticalCollision(world, &verticalIntersect, verticalIntersectStep.X) {
				result = &verticalIntersect
				break
			} else {
				verticalIntersect.X += verticalIntersectStep.X
				verticalIntersect.Y += verticalIntersectStep.Y
				rayLength = verticalIntersectDistance
			}
		}
	}
	return result
}

//checkHorizontalCollision check if a point on a horizontal-line on the grid is hitting a wall given the ray's vertical-direction.
func (*RayCasterImpl) checkHorizontalCollision(world world.WorldMap, horizontalIntersect *innerMath.Point2D, horizontalIntersectStepY float64) bool {
	result := false
	if horizontalIntersectStepY > 0 {
		if world.GetCellValue(int(horizontalIntersect.X), int(horizontalIntersect.Y)) == 1 {
			result = true
		}
	} else {
		if world.GetCellValue(int(horizontalIntersect.X), int(horizontalIntersect.Y)-1) == 1 {
			result = true
		}
	}
	return result
}

//checkVerticalCollision check if a point on a vertical-line on the grid is hitting a wall given the ray's horizontal-direction.
func (*RayCasterImpl) checkVerticalCollision(world world.WorldMap, verticalIntersect *innerMath.Point2D, verticalIntersectStepX float64) bool {
	result := false
	if verticalIntersectStepX > 0 {
		if world.GetCellValue(int(verticalIntersect.X), int(verticalIntersect.Y)) == 1 {
			result = true
		}
	} else {
		if world.GetCellValue(int(verticalIntersect.X)-1, int(verticalIntersect.Y)) == 1 {
			result = true
		}
	}
	return result
}
