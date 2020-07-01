package impl

import (
	animatedelement "francoisgergaud/3dGame/common/environment/animatedelement"
	animatedelementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	publisher "francoisgergaud/3dGame/common/event/publisher"
	publisherImpl "francoisgergaud/3dGame/common/event/publisher/impl"
	internalMath "francoisgergaud/3dGame/common/math"
	mathHelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/server/bot"
	"math"

	"github.com/gdamore/tcell"
)

//Bot is an NPC animated-element
type Bot interface {
	publisher.EventPublisher
	animatedelement.AnimatedElement
}

//NewBotImpl build a new bot implementation.
func NewBotImpl(id string, initialPosition *internalMath.Point2D, initialAngle, velocity, stepAngle, size float64, moveDirection, rotateDirection state.Direction, style tcell.Style, world world.WorldMap, mathHelper mathHelper.MathHelper, quit <-chan interface{}) bot.Bot {
	result := BotImpl{
		AnimatedElement: animatedelementImpl.NewAnimatedElement(id, initialPosition, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, world, mathHelper),
		EventPublisher:  publisherImpl.NewEventPublisherImpl(),
		mathHelper:      mathHelper,
		world:           world,
	}
	return &result
}

//BotImpl is a bot implementation.
type BotImpl struct {
	animatedelement.AnimatedElement
	publisher.EventPublisher
	world      world.WorldMap
	mathHelper mathHelper.MathHelper
}

//Move the bot's position depending on the colision of walls
func (bot *BotImpl) Move() {
	botState := bot.State()
	rayDestination := bot.mathHelper.CastRay(botState.Position, bot.world, botState.Angle, botState.Velocity)
	if rayDestination != nil {
		//horizontal rebound
		if rayDestination.X-math.Floor(rayDestination.X) < 0.0001 && rayDestination.X-math.Floor(rayDestination.X) > -0.0001 {
			switch {
			case botState.Angle <= 1.0:
				botState.Angle = 1.0 - botState.Angle
			case botState.Angle <= 2.0:
				botState.Angle = 3.0 - botState.Angle
			}
		}
		//vertical rebound
		if rayDestination.Y-math.Floor(rayDestination.Y) < 0.0001 && rayDestination.Y-math.Floor(rayDestination.Y) > -0.0001 {
			switch {
			case botState.Angle <= 1.0:
				botState.Angle = 2.0 - botState.Angle
			case botState.Angle <= 2.0:
				botState.Angle = 2.0 - botState.Angle
			}
		}
		//distanceToWall := worldElement.GetState().Position.Distance(rayDestination)
		//state.Position.X = rayDestination.X + math.Cos(state.Angle*math.Pi)*(state.Velocity-distanceToWall)
		//state.Position.Y = rayDestination.Y + math.Sin(state.Angle*math.Pi)*(state.Velocity-distanceToWall)
		event := event.Event{
			PlayerID: bot.ID(),
			Action:   "move",
			State:    botState,
		}
		bot.PublishEvent(event)
	} else {
		bot.AnimatedElement.Move()
	}

}
