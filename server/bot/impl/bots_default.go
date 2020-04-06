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

type Bot interface {
	publisher.EventPublisher
	animatedelement.AnimatedElement
}

//NewBotImpl build a new bot implementation.
func NewBotImpl(id string, initialPosition *internalMath.Point2D, initialAngle, velocity, stepAngle, size float64, moveDirection, rotateDirection state.Direction, style tcell.Style, world world.WorldMap, mathHelper mathHelper.MathHelper, quit <-chan interface{}) bot.Bot {
	result := BotImpl{
		id:              id,
		AnimatedElement: animatedelementImpl.NewAnimatedElement(initialPosition, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, world, mathHelper),
		EventPublisher:  publisherImpl.NewEventPublisherImpl(),
		mathHelper:      mathHelper,
		world:           world,
	}
	return &result
}

//BotImpl is a bot implementation.
type BotImpl struct {
	id string
	animatedelement.AnimatedElement
	publisher.EventPublisher
	world      world.WorldMap
	mathHelper mathHelper.MathHelper
}

//Move the bot's position depending on the colision of walls
func (bot *BotImpl) Move() {
	state := bot.State()
	rayDestination := bot.mathHelper.CastRay(state.Position, bot.world, state.Angle, state.Velocity)
	if rayDestination != nil {
		//horizontal rebound
		if rayDestination.X-math.Floor(rayDestination.X) < 0.0001 && rayDestination.X-math.Floor(rayDestination.X) > -0.0001 {
			switch {
			case state.Angle <= 1.0:
				state.Angle = 1.0 - state.Angle
			case state.Angle <= 2.0:
				state.Angle = 3.0 - state.Angle
			}
		}
		//vertical rebound
		if rayDestination.Y-math.Floor(rayDestination.Y) < 0.0001 && rayDestination.Y-math.Floor(rayDestination.Y) > -0.0001 {
			switch {
			case state.Angle <= 1.0:
				state.Angle = 2.0 - state.Angle
			case state.Angle <= 2.0:
				state.Angle = 2.0 - state.Angle
			}
		}
		//distanceToWall := worldElement.GetState().Position.Distance(rayDestination)
		//state.Position.X = rayDestination.X + math.Cos(state.Angle*math.Pi)*(state.Velocity-distanceToWall)
		//state.Position.Y = rayDestination.Y + math.Sin(state.Angle*math.Pi)*(state.Velocity-distanceToWall)
		event := event.Event{
			PlayerID: bot.id,
			Action:   "move",
			State:    state,
		}
		bot.PublishEvent(event)
	} else {
		bot.AnimatedElement.Move()
	}

}
