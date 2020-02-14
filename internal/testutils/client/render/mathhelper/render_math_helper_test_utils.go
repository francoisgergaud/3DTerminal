package testmathhelper

import (
	"francoisgergaud/3dGame/common/math"

	"github.com/stretchr/testify/mock"
)

//MockRendererMathHelper mocks the calls to the RendererMathHelper interface.
type MockRendererMathHelper struct {
	mock.Mock
}

//CalculateProjectionDistance mocks the method of the same name
func (mock *MockRendererMathHelper) CalculateProjectionDistance(startPosition *math.Point2D, endPosition *math.Point2D, angle float64) float64 {
	args := mock.Called(startPosition, endPosition, angle)
	return args.Get(0).(float64)
}

//IsWallAngle mocks the method of the same name
func (mock *MockRendererMathHelper) IsWallAngle(point *math.Point2D) bool {
	args := mock.Called(point)
	return args.Bool(0)
}

//GetRayTracingAngleForColumn mocks the method of the same name
func (mock *MockRendererMathHelper) GetRayTracingAngleForColumn(angle float64, columnIndex, screenWidth int, viewAngle float64) float64 {
	args := mock.Called(angle, columnIndex, screenWidth, viewAngle)
	return args.Get(0).(float64)
}

//GetFillRowRange mocks the method of the same name
func (mock *MockRendererMathHelper) GetFillRowRange(distance, screenHeight float64) (int, int) {
	args := mock.Called(distance, screenHeight)
	return args.Int(0), args.Int(1)
}
