package player

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/event/publisher"
	"francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"
	"time"

	"github.com/gdamore/tcell"

	animatedelementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"

	eventPublisherImpl "francoisgergaud/3dGame/common/event/publisher/impl"
)

//NewPlayer creates a new player
func NewPlayer(id string, worldMap world.WorldMap, mathHelper helper.MathHelper, quit <-chan interface{}) animatedelement.AnimatedElement {
	animatedElementState := state.AnimatedElementState{
		Position:        &math.Point2D{X: 5, Y: 5},
		Angle:           0.0,
		Size:            0.5,
		Velocity:        0.1,
		StepAngle:       0.01,
		Style:           tcell.StyleDefault.Background(tcell.Color126),
		MoveDirection:   state.None,
		RotateDirection: state.None,
	}
	return animatedelementImpl.NewAnimatedElementWithState(id, &animatedElementState, worldMap, mathHelper)
}

//Spawner is in charge to spawn an animated-element
type Spawner interface {
	Spawn(string, state.Direction)
	publisher.EventPublisher
}

//NewStaticSpawner is a factory for static-spawner
func NewStaticSpawner(players map[string]animatedelement.AnimatedElement) *StaticSpawner {
	return &StaticSpawner{
		timeMs:                 time.Duration(2000),
		EventPublisher:         eventPublisherImpl.NewEventPublisherImpl(),
		players:                players,
		playersWaitingForSpawn: make(map[string]animatedelement.AnimatedElement),
	}
}

//StaticSpawner spawn an animated-element in a static state
type StaticSpawner struct {
	timeMs                 time.Duration
	players                map[string]animatedelement.AnimatedElement
	playersWaitingForSpawn map[string]animatedelement.AnimatedElement
	publisher.EventPublisher
}

//Spawn the animated-element
func (spawner *StaticSpawner) Spawn(animatedelementID string, moveDirection state.Direction) {
	timer := time.NewTimer(spawner.timeMs)
	animatedElement := spawner.players[animatedelementID]
	delete(spawner.players, animatedelementID)
	spawner.playersWaitingForSpawn[animatedelementID] = animatedElement
	go func() {
		select {
		case <-timer.C:
			animatedElementState := animatedElement.State()
			animatedElementState.Position = &math.Point2D{X: 5, Y: 5}
			animatedElementState.Angle = 0.0
			animatedElementState.MoveDirection = moveDirection
			animatedElementState.RotateDirection = state.None
			delete(spawner.playersWaitingForSpawn, animatedelementID)
			spawner.players[animatedelementID] = animatedElement
			spawner.PublishEvent(
				event.Event{
					Action:   "spawn",
					PlayerID: animatedelementID,
					State:    animatedElementState,
				},
			)
		}
	}()
}
