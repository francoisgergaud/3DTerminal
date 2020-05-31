package player

import (
	animatedelementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBot(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	quit := make(chan struct{})
	player := NewPlayer(worldMap, mathHelper, quit)
	assert.IsType(t, &animatedelementImpl.AnimatedElementImpl{}, player)
}
