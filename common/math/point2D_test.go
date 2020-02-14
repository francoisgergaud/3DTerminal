package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint2DClone(t *testing.T) {
	point1 := &Point2D{X: 0, Y: 0}
	point2 := point1.Clone()
	assert.True(t, point1 != point2)
	assert.Equal(t, point1.X, point2.X)
	assert.Equal(t, point1.Y, point2.Y)
}

func TestDistanceSamePoint(t *testing.T) {
	point1 := &Point2D{X: 0, Y: 0}
	point2 := point1.Clone()
	assert.Equal(t, 0.0, point1.Distance(point2))
}

func TestDistancePointSameX(t *testing.T) {
	point1 := &Point2D{X: 0, Y: 0}
	point2 := &Point2D{X: 0, Y: 1}
	assert.Equal(t, 1.0, point1.Distance(point2))
}

func TestDistancePointSameY(t *testing.T) {
	point1 := &Point2D{X: 0, Y: 0}
	point2 := &Point2D{X: 1, Y: 0}
	assert.Equal(t, 1.0, point1.Distance(point2))
}

func TestDistance(t *testing.T) {
	point1 := &Point2D{X: 0, Y: 0}
	point2 := &Point2D{X: 1, Y: 1}
	assert.True(t, AreFloatAlmostEquals(1.414, point1.Distance(point2), 0.001))
}

func AreFloatAlmostEquals(f1, f2 float64, precision float64) bool {
	return math.Abs(f1-f2) < precision
}
