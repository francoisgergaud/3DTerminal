package testworld

import "github.com/stretchr/testify/mock"

import "francoisgergaud/3dGame/common/environment/world"

//MockWorldMap mocks a WorldMap
type MockWorldMap struct {
	mock.Mock
}

// GetCellValue : mocks the call to the GetCellValue method.
func (mock *MockWorldMap) GetCellValue(x, y int) int {
	args := mock.Called(x, y)
	return args.Int(0)
}

//Clone mocks the call to the Clone
func (mock *MockWorldMap) Clone() world.WorldMap {
	args := mock.Called()
	return args.Get(0).(world.WorldMap)
}

//MockWorldMapWithGrid mocks a WorldMap using a grid to return the cell's values.
type MockWorldMapWithGrid struct {
	mock.Mock
	Grid [][]int
}

// GetCellValue : returns the Map's cell value. If coordinate are out of map, returns 0
func (mock *MockWorldMapWithGrid) GetCellValue(x, y int) int {
	mock.Called(x, y)
	if x >= 0 && x < len(mock.Grid) && y >= 0 && y < len(mock.Grid[x]) {
		return mock.Grid[x][y]
	}
	return 0
}

//Clone mocks the call to the Clone
func (mock *MockWorldMapWithGrid) Clone() world.WorldMap {
	args := mock.Called()
	return args.Get(0).(world.WorldMap)
}
