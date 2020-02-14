package testhelper

import (
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/math"

	"github.com/stretchr/testify/mock"
)

//MockMathHelper mocks the calls to the MathHelper interface.
type MockMathHelper struct {
	mock.Mock
}

//CastRay mocks the call to the method of the same name.
func (mock *MockMathHelper) CastRay(origin *math.Point2D, worldMap world.WorldMap, rayAngle, visibility float64) *math.Point2D {
	args := mock.Called(origin, worldMap, rayAngle, visibility)
	if args.Get(0) != nil {
		return args.Get(0).(*math.Point2D)
	}
	return nil
}

//GetWorldElementProjection mocks the call to the method of the same name.
func (mock *MockMathHelper) GetWorldElementProjection(
	playerPosition *math.Point2D,
	viewAngle float64,
	fov float64,
	worlElementPosition *math.Point2D,
	worldElementSize float64) (isVisible bool, startColumnNumber float64, startOffset float64, endColumnNumber float64, endOffset float64) {
	args := mock.Called(playerPosition, viewAngle, fov, worlElementPosition, worldElementSize)
	return args.Bool(0), args.Get(1).(float64), args.Get(2).(float64), args.Get(3).(float64), args.Get(4).(float64)
}

//NormalizeAngle mocks the call to the method of the same name.
func (mock *MockMathHelper) NormalizeAngle(angle float64) float64 {
	args := mock.Called(angle)
	return args.Get(0).(float64)
}
