package bot

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	internalmath "francoisgergaud/3dGame/common/math"
	mathhelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/server/bot"
	botImpl "francoisgergaud/3dGame/server/bot/impl"

	"github.com/gdamore/tcell"
)

//NewBot creates a bot
func NewBot(id string, worldMap world.WorldMap, mathHelper mathhelper.MathHelper, quit chan struct{}) bot.Bot {
	position := &internalmath.Point2D{X: 9, Y: 12}
	initialAngle := 0.3
	velocity := 0.02
	size := 0.3
	stepAngle := 0.0
	moveDirection := state.Forward
	rotateDirection := state.None
	style := tcell.StyleDefault.Background(tcell.ColorDarkBlue)
	return botImpl.NewBotImpl(id, position, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper, quit)
}
