package world

import (
	"math/rand"
)

//WorldMap is a world-map defining a gird a elements.
type WorldMap interface {
	GetCellValue(x, y int) int
	Clone() WorldMap
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
	if y >= 0 && y < len(w.grid) && x >= 0 && x < len(w.grid[y]) {
		return w.grid[y][x]
	}
	return 0
}

//Clone creates a deep-copy.
func (w *WorldMapImpl) Clone() WorldMap {
	if len(w.grid) > 0 {
		grid := make([][]int, len(w.grid))
		for i := 0; i < len(w.grid); i++ {
			grid[i] = make([]int, len(w.grid[i]))
			for j := 0; j < len(w.grid[i]); j++ {
				grid[i][j] = w.grid[i][j]

			}
		}
		return &WorldMapImpl{
			grid: grid,
		}
	} else {
		return &WorldMapImpl{
			grid: make([][]int, 0),
		}
	}
}
