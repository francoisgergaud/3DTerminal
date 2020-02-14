package world

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
	assert.Equal(t, worldMap.GetCellValue(-1, 10), 0)
	assert.Equal(t, worldMap.GetCellValue(10, -1), 0)
}

func TestInitializeRandom(t *testing.T) {
	width := 5
	height := 4
	worldMap := new(WorldMapImpl)
	worldMap.InitializeRandom(width, height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			assert.Contains(t, [...]int{0, 1}, worldMap.GetCellValue(i, j))
		}
	}
}

func TestCloneEmptyMap(t *testing.T) {
	width := 0
	height := 0
	worldMap := new(WorldMapImpl)
	worldMap.InitializeRandom(width, height)
	worldMapCloned := worldMap.Clone()
	assert.NotNil(t, worldMapCloned)
}

func TestCloneNonEmptyMap(t *testing.T) {
	width := 5
	height := 4
	worldMap := new(WorldMapImpl)
	worldMap.InitializeRandom(width, height)
	worldMapCloned := worldMap.Clone()
	assert.True(t, worldMap != worldMapCloned)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			assert.Equal(t, worldMap.GetCellValue(i, j), worldMapCloned.GetCellValue(i, j))
		}
	}
}
