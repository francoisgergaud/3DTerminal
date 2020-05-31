package worldmap

import (
	"francoisgergaud/3dGame/common/environment/world"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBot(t *testing.T) {
	worldMap := NewWorldMap()
	assert.IsType(t, &world.WorldMapImpl{}, worldMap)
}
