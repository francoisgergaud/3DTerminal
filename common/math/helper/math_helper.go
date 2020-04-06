package helper

import (
	"fmt"
	"francoisgergaud/3dGame/common/environment/world"
	innerMath "francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/raycaster"
	"math"
)

//MathHelper provides maths for background render. It is an internal struct. Its purpose
//is to make the renderer's code more modular for testing purpose.
// This struct is stateless (no internal state) but usefull for dependency-injection and testing.
type MathHelper interface {
	CastRay(origin *innerMath.Point2D, worldMap world.WorldMap, rayAngle, visibility float64) *innerMath.Point2D
	GetWorldElementProjection(
		playerPosition *innerMath.Point2D,
		viewAngle float64,
		fov float64,
		worlElementPosition *innerMath.Point2D,
		worldElementSize float64) (isVisible bool, startColumnNumber float64, startOffset float64, endColumnNumber float64, endOffset float64)
	NormalizeAngle(angle float64) float64
}

//NewMathHelper builds a new MathHelper.
func NewMathHelper(rayCaster raycaster.RayCaster) (*MathHelperImpl, error) {
	if rayCaster == nil {
		return nil, fmt.Errorf("Math-helper 'raycaster' cannot be nil")
	}
	return &MathHelperImpl{
		raycaster: rayCaster,
	}, nil
}

//MathHelperImpl implements the IBackgroundRendererMathHelper interface.
type MathHelperImpl struct {
	raycaster raycaster.RayCaster
}

//NormalizeAngle nromalize and angle float-value by keeping it between 0 and 2.
func (mathHelper *MathHelperImpl) NormalizeAngle(angle float64) float64 {
	if angle >= 2 {
		angle -= 2
	} else if angle < 0 {
		angle += 2
	}
	return angle
}

//CastRay casts a ray from:
// - an origin-position
// - a world-map (which contains the walls for collision with the ray)
// - a ray's angle
// -a visibility, the max distance a ray can be. If the ray does not encounter a wall with this distance, the Raycast returns null (infinite ray)
func (mathHelper *MathHelperImpl) CastRay(origin *innerMath.Point2D, worldMap world.WorldMap, angle, maxDistance float64) *innerMath.Point2D {
	return mathHelper.raycaster.CastRay(origin, worldMap, angle, maxDistance)
}

//GetWorldElementProjection returns the projection data of a world element given:
//- the player's position
//- the player's angle
//- the field of view angle
//- the world-elemnt's position
//- the world-element's size
// it returns:
//- the start-ratio on the screen-width from which the world-element draw must start
//- the world-element's width-offset at the start-index
//- the end-ratio on the screen-width to which the world-element draw must end
//- the world-element's width-offset at the end-index
func (mathHelper *MathHelperImpl) GetWorldElementProjection(
	playerPosition *innerMath.Point2D,
	viewAngle float64,
	fov float64,
	worlElementPosition *innerMath.Point2D,
	worldElementSize float64) (isVisible bool, startScreenRatio float64, startOffset float64, endScreenNumber float64, endOffset float64) {
	worldElementVectorFromPlayer := &innerMath.Point2D{
		X: worlElementPosition.X - playerPosition.X,
		Y: worlElementPosition.Y - playerPosition.Y,
	}
	// var worldElementAngle float64
	// if worldElementVectorFromPlayer.X > 0.001 || worldElementVectorFromPlayer.X < -0.001 {
	// 	worldElementAngle = normalizeAngle(math.Atan(worldElementVectorFromPlayer.Y/worldElementVectorFromPlayer.X) / math.Pi)
	// } else if worldElementVectorFromPlayer.Y > 0 {
	// 	worldElementAngle = 0.5
	// } else {
	// 	worldElementAngle = 1.5
	// }
	worldElementAngle := mathHelper.NormalizeAngle(math.Atan2(worldElementVectorFromPlayer.Y, worldElementVectorFromPlayer.X) / math.Pi)
	worlelementHitBoxBorderDeltaAngle := mathHelper.NormalizeAngle(math.Atan(worldElementSize/playerPosition.Distance(worlElementPosition)) / math.Pi)
	//worldElementStartAngle and worldElementEndAngle are and can only be clockwise.
	worldElementStartAngle := mathHelper.NormalizeAngle(worldElementAngle - worlelementHitBoxBorderDeltaAngle)
	worldElementEndAngle := mathHelper.NormalizeAngle(worldElementAngle + worlelementHitBoxBorderDeltaAngle)
	// check the start and end angles anf the plaer-angle +/- field-of-view
	fovStartAngle := mathHelper.NormalizeAngle(viewAngle - (fov / 2))
	fovEndAngle := mathHelper.NormalizeAngle(viewAngle + (fov / 2))
	if isAngleBetween(fovEndAngle, fovStartAngle, worldElementStartAngle) && isAngleBetween(fovEndAngle, fovStartAngle, worldElementEndAngle) {
		if isAngleBetween(worldElementStartAngle, worldElementEndAngle, viewAngle) {
			isVisible = true
			startScreenRatio = 0
			endScreenNumber = 1
			worldElementTotalAngle := mathHelper.NormalizeAngle(worldElementEndAngle - worldElementStartAngle)
			startOffset = mathHelper.NormalizeAngle(fovStartAngle-worldElementStartAngle) / worldElementTotalAngle
			endOffset = mathHelper.NormalizeAngle(fovEndAngle-worldElementStartAngle) / mathHelper.NormalizeAngle(worldElementEndAngle-worldElementStartAngle)
		} else {
			isVisible = false //both worldElementStartAngle and worldElementEndAngle are outside fov, and worldElementStartAngle and worldElementEndAngle are and can only be clockwise, NOT VISIBLE
		}
	} else {
		isVisible = true
		if isAngleBetween(fovStartAngle, fovEndAngle, worldElementStartAngle) {
			startScreenRatio = mathHelper.NormalizeAngle(worldElementStartAngle-fovStartAngle) / fov
			startOffset = 0
		} else {
			startScreenRatio = 0
			startOffset = mathHelper.NormalizeAngle(fovStartAngle-worldElementStartAngle) / mathHelper.NormalizeAngle(worldElementEndAngle-worldElementStartAngle)
		}
		if isAngleBetween(fovStartAngle, fovEndAngle, worldElementEndAngle) {
			endScreenNumber = mathHelper.NormalizeAngle(worldElementEndAngle-fovStartAngle) / fov
			endOffset = 1
		} else {
			endScreenNumber = 1
			endOffset = mathHelper.NormalizeAngle(fovEndAngle-worldElementStartAngle) / mathHelper.NormalizeAngle(worldElementEndAngle-worldElementStartAngle)
		}
	}
	return
}

//return true if an angle is between start-angle and end-angle clockwise.
func isAngleBetween(start, end, angle float64) bool {
	result := false
	if start <= end {
		if angle >= start && angle <= end {
			result = true
		}
	} else {
		if angle <= end || angle >= start {
			result = true
		}
	}
	return result
}
