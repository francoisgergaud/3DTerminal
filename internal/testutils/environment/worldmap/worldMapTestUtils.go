package worldmap

import "github.com/stretchr/testify/mock"

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
