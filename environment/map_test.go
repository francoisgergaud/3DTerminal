package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var grid = [][]int{
	{1, 1},
	{1, 0},
}

func TestGetCellValue(t *testing.T) {
	worldMap := NewWorldMap(grid)
	assert.Equal(t, worldMap.GetCellValue(0, 0), 1)
	assert.Equal(t, worldMap.GetCellValue(1, 1), 0)
}

func TestGetCellValueOutOfGrid(t *testing.T) {
	worldMap := NewWorldMap(grid)
	assert.Equal(t, worldMap.GetCellValue(10, 10), 0)
	assert.Equal(t, worldMap.GetCellValue(-1, -1), 0)
}
