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
	Grid [][]int
}

// InitializeRandom : Initialize the map with random 1 or 0 values cells
func (w *WorldMapImpl) InitializeRandom(width, height int) {
	rand.Seed(86)
	w.Grid = make([][]int, width)
	for rowIndex := range w.Grid {
		w.Grid[rowIndex] = make([]int, height)
		for columnIndex := range w.Grid[rowIndex] {
			w.Grid[rowIndex][columnIndex] = rand.Intn(2)
		}
	}
}

//NewWorldMap builds a new world-map from the input parameters.
func NewWorldMap(grid [][]int) WorldMap {
	return &WorldMapImpl{
		Grid: grid,
	}
}

// GetCellValue : returns the Map's cell value. If coordinate are out of map, returns 0
func (w *WorldMapImpl) GetCellValue(x, y int) int {
	if y >= 0 && y < len(w.Grid) && x >= 0 && x < len(w.Grid[y]) {
		return w.Grid[y][x]
	}
	return 0
}

//Clone creates a deep-copy.
func (w *WorldMapImpl) Clone() WorldMap {
	if len(w.Grid) > 0 {
		grid := make([][]int, len(w.Grid))
		for i := 0; i < len(w.Grid); i++ {
			grid[i] = make([]int, len(w.Grid[i]))
			for j := 0; j < len(w.Grid[i]); j++ {
				grid[i][j] = w.Grid[i][j]

			}
		}
		return &WorldMapImpl{
			Grid: grid,
		}
	} else {
		return &WorldMapImpl{
			Grid: make([][]int, 0),
		}
	}
}
