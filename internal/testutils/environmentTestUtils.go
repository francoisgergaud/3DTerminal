package testutils

import (
	"francoisgergaud/3dGame/common"

	"github.com/gdamore/tcell"

	"github.com/stretchr/testify/mock"
)

//MockCharacter mocks a character
type MockCharacter struct {
	mock.Mock
}

//GetPosition returns the player's position.
func (mock *MockCharacter) GetPosition() *common.Point2D {
	args := mock.Called()
	return args.Get(0).(*common.Point2D)
}

//GetAngle returns the player's orientation angle.
func (mock *MockCharacter) GetAngle() float64 {
	args := mock.Called()
	return args.Get(0).(float64)
}

//Action mocks the opration.
func (mock *MockCharacter) Action(eventKey *tcell.EventKey) {
	mock.Called(eventKey)
}

//MockWorldMap mocks a WorldMap
type MockWorldMap struct {
	mock.Mock
}

// GetCellValue : mocks the call to the GetCellValue method.
func (mock *MockWorldMap) GetCellValue(x, y int) int {
	args := mock.Called(x, y)
	return args.Int(0)
}

//MockWorldMapWithGrid mocks a WorldMap using a grid to return the cell's values.
type MockWorldMapWithGrid struct {
	Grid [][]int
}

// GetCellValue : returns the Map's cell value. If coordinate are out of map, returns 0
func (mock *MockWorldMapWithGrid) GetCellValue(x, y int) int {
	if x >= 0 && x < len(mock.Grid) && y >= 0 && y < len(mock.Grid[x]) {
		return mock.Grid[x][y]
	}
	return 0
}