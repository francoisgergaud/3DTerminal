package impl

import (
	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/mock"
)

//MockRaySampler mocks a ray-sampler.
type MockRaySampler struct {
	mock.Mock
}

//GetBackgroundRune mocks the operation of the same name from the RaySampler interface.
func (mock *MockRaySampler) GetBackgroundRune(rowIndex int) rune {
	args := mock.Called(rowIndex)
	return args.Get(0).(rune)
}

//GetWallRune mocks the operation of the same name from the RaySampler interface.
func (mock *MockRaySampler) GetWallRune(rowIndex int) rune {
	args := mock.Called(rowIndex)
	return args.Get(0).(rune)
}

//GetBackgroundStyle mocks the operation of the same name from the RaySampler interface.
func (mock *MockRaySampler) GetBackgroundStyle(rowIndex int) tcell.Style {
	args := mock.Called(rowIndex)
	return args.Get(0).(tcell.Style)
}

//GetWallStyleFromDistance mocks the operation of the same name from the RaySampler interface.
func (mock *MockRaySampler) GetWallStyleFromDistance(distance float64) tcell.Style {
	args := mock.Called(distance)
	return args.Get(0).(tcell.Style)
}
