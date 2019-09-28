package environment

import (
	"math/rand"
)

//WorldMap is a world-map defining a gird a elements.
type WorldMap interface {
	GetCellValue(x, y int) int
}

// WorldMapImpl implements the WorldMap interface.
type WorldMapImpl struct {
	grid [][]int
}

// InitializeRandom : Initialize the map with random 1 or 0 values cells
func (w *WorldMapImpl) InitializeRandom(width, height int) {
	rand.Seed(86)
	w.grid = make([][]int, width)
	for rowIndex := range w.grid {
		w.grid[rowIndex] = make([]int, height)
		for columnIndex := range w.grid[rowIndex] {
			w.grid[rowIndex][columnIndex] = rand.Intn(2)
		}
	}
}

//NewWorldMap builds a new world-map from the input parameters.
func NewWorldMap(grid [][]int) WorldMap {
	return &WorldMapImpl{
		grid: grid,
	}
}

// GetCellValue : returns the Map's cell value. If coordinate are out of map, returns 0
func (w *WorldMapImpl) GetCellValue(x, y int) int {
	if x >= 0 && x < len(w.grid) && y >= 0 && y < len(w.grid[x]) {
		return w.grid[x][y]
	}
	return 0
}
