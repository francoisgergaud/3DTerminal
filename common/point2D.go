package common

import (
	"fmt"
	"math"
)

//Point2D is a 2-dimensional coordinate (can also be used as a vector).
type Point2D struct {
	X, Y float64
}

//Distance to another point
func (point *Point2D) Distance(otherPoint *Point2D) float64 {
	return math.Hypot(otherPoint.X-point.X, otherPoint.Y-point.Y)
}

func (point Point2D) String() string {
	return fmt.Sprintf("{X:%v, Y:%v}", point.X, point.Y)
}

const float64EqualityThreshold = 1e-3

//AlmostEquals checks if 2 points with floating-point precision could be considered as really closed.
func (point Point2D) AlmostEquals(otherPoint Point2D) bool {
	return math.Abs(otherPoint.X-point.X) <= float64EqualityThreshold && math.Abs(otherPoint.Y-point.Y) <= float64EqualityThreshold
}
