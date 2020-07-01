package player

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	animatedelementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/math"
	testanimatedelement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testeventpublisher "francoisgergaud/3dGame/internal/testutils/common/event/publisher"
	testhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewBot(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	quit := make(chan interface{})
	player := NewPlayer("id", worldMap, mathHelper, quit)
	assert.IsType(t, &animatedelementImpl.AnimatedElementImpl{}, player)
}

func TestStaticSpawnerSpawn(t *testing.T) {
	animatedElementID := "idtest"
	eventPublisher := new(testeventpublisher.MockEventPublisher)
	players := make(map[string]animatedelement.AnimatedElement)
	spawner := &StaticSpawner{
		timeMs:                 time.Duration(2),
		EventPublisher:         eventPublisher,
		playersWaitingForSpawn: make(map[string]animatedelement.AnimatedElement),
		players:                players,
	}
	animatedElement := new(testanimatedelement.MockAnimatedElement)
	animatedElementState := new(state.AnimatedElementState)
	animatedElement.On("State").Return(animatedElementState)
	players[animatedElementID] = animatedElement
	var eventPublished event.Event
	eventPublisher.On("PublishEvent", mock.MatchedBy(
		func(eventParaemter event.Event) bool {
			eventPublished = eventParaemter
			return true
		},
	))
	spawner.Spawn(animatedElementID, state.Forward)
	assert.Contains(t, spawner.playersWaitingForSpawn, animatedElementID)
	assert.NotContains(t, spawner.players, animatedElementID)
	time.Sleep(time.Duration(5 * time.Millisecond))
	assert.Contains(t, spawner.players, animatedElementID)
	assert.NotContains(t, spawner.playersWaitingForSpawn, animatedElementID)
	assert.Equal(t, animatedElementState.Style, eventPublished.State.Style)
	assert.Equal(t, animatedElementState.Velocity, eventPublished.State.Velocity)
	assert.Equal(t, animatedElementState.StepAngle, eventPublished.State.StepAngle)
	assert.Equal(t, animatedElementState.Size, eventPublished.State.Size)
	assert.Equal(t, state.Forward, eventPublished.State.MoveDirection)
	assert.Equal(t, &math.Point2D{X: 5, Y: 5}, eventPublished.State.Position)
	assert.Equal(t, 0.0, eventPublished.State.Angle)
	assert.Equal(t, animatedElementID, eventPublished.PlayerID)
	assert.Equal(t, "spawn", eventPublished.Action)
	mock.AssertExpectationsForObjects(t, animatedElement, eventPublisher)
}

func TestNewNewStaticSpawner(t *testing.T) {
	players := make(map[string]animatedelement.AnimatedElement)
	spawner := NewStaticSpawner(players)
	assert.Equal(t, time.Duration(2000), spawner.timeMs, players)
	assert.NotNil(t, spawner.EventPublisher)
	assert.Equal(t, spawner.players, players)
	assert.NotNil(t, spawner.playersWaitingForSpawn)
}
