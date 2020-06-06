package player

import (
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	animatedelementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/event/publisher"
	publisherImpl "francoisgergaud/3dGame/common/event/publisher/impl"
	"francoisgergaud/3dGame/common/math/helper"

	"github.com/gdamore/tcell"
)

//NewPlayer builds a new player from ithe input parameters.
func NewPlayer(playerState *state.AnimatedElementState, world world.WorldMap, mathHelper helper.MathHelper) player.Player {
	return &Impl{
		AnimatedElement: animatedelementImpl.NewAnimatedElementWithState(playerState, world, mathHelper),
		EventPublisher:  publisherImpl.NewEventPublisherImpl(),
	}
}

//Impl is the default implementation of a player
type Impl struct {
	animatedelement.AnimatedElement
	publisher.EventPublisher
}

// Action the player according to the input key
func (p *Impl) Action(eventKey *tcell.EventKey) {
	playerState := p.State()
	switch eventKey.Key() {
	case tcell.KeyUp:
		if playerState.MoveDirection == state.Backward {
			playerState.MoveDirection = state.None
		} else {
			playerState.MoveDirection = state.Forward
		}
	case tcell.KeyDown:
		if playerState.MoveDirection == state.Forward {
			playerState.MoveDirection = state.None
		} else {
			playerState.MoveDirection = state.Backward
		}
	case tcell.KeyLeft:
		if playerState.RotateDirection == state.Right {
			playerState.RotateDirection = state.None
		} else {
			playerState.RotateDirection = state.Left
		}
	case tcell.KeyRight:
		if playerState.RotateDirection == state.Left {
			playerState.RotateDirection = state.None
		} else {
			playerState.RotateDirection = state.Right
		}
	}
	p.PublishEvent(event.Event{Action: "move", State: p.State(), TimeFrame: 0})
}
