package bot

import (
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	"francoisgergaud/3dGame/server/bot/impl"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBot(t *testing.T) {
	id := "botID"
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	quit := make(chan interface{})
	bot := NewBot(id, worldMap, mathHelper, quit)
	assert.IsType(t, &impl.BotImpl{}, bot)
}
