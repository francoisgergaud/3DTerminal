package projectile

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	animatedElementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/event/publisher"
	"francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testeventpublisher "francoisgergaud/3dGame/internal/testutils/common/event/publisher"
	testhelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewProjectile(t *testing.T) {
	projectileID := "idTest"
	projectileStartPosition := &math.Point2D{X: 0.5, Y: 0.0}
	angle := 1.75
	world := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	otherPlayers := make(map[string]animatedelement.AnimatedElement)

	projectile := NewProjectile(projectileID, projectileStartPosition, angle, world, otherPlayers, mathHelper)

	assert.Equal(t, projectileID, projectile.ID())
	assert.Equal(t, angle, projectile.State().Angle)
	assert.Equal(t, projectileStartPosition, projectile.State().Position)
}

func TestDetectImpactWithOtherPlayer1(t *testing.T) {
	projectileStartPosition := &math.Point2D{X: 0.0, Y: 0.0}
	projectileEndPosition := &math.Point2D{X: 5.0, Y: 0.0}
	otherPlayerPosition := &math.Point2D{X: 4.0, Y: 0.25}
	playerSize := 0.6
	impact := detectImpactWithOtherPlayer(projectileStartPosition, projectileEndPosition, otherPlayerPosition, playerSize)
	assert.Equal(t, &math.Point2D{X: 4.0, Y: 0.0}, impact)
}

func TestDetectImpactWithOtherPlayer2(t *testing.T) {
	projectileStartPosition := &math.Point2D{X: 5.0, Y: 0.0}
	projectileEndPosition := &math.Point2D{X: 0.0, Y: 0.0}
	otherPlayerPosition := &math.Point2D{X: 4.0, Y: 0.25}
	playerSize := 0.6
	impact := detectImpactWithOtherPlayer(projectileStartPosition, projectileEndPosition, otherPlayerPosition, playerSize)
	assert.Equal(t, &math.Point2D{X: 4.0, Y: 0.0}, impact)
}

func TestDetectImpactWithOtherPlayer3(t *testing.T) {
	projectileStartPosition := &math.Point2D{X: 5.0, Y: 0.0}
	projectileEndPosition := &math.Point2D{X: 0.0, Y: 0.0}
	otherPlayerPosition := &math.Point2D{X: 6.0, Y: 0.25}
	playerSize := 0.6
	impact := detectImpactWithOtherPlayer(projectileStartPosition, projectileEndPosition, otherPlayerPosition, playerSize)
	assert.Nil(t, impact)
}

func TestDetectImpactWithOtherPlayer4(t *testing.T) {
	projectileStartPosition := &math.Point2D{X: 0.0, Y: 0.0}
	projectileEndPosition := &math.Point2D{X: 0.0, Y: 3.0}
	otherPlayerPosition := &math.Point2D{X: 0.25, Y: 0.5}
	playerSize := 0.6
	impact := detectImpactWithOtherPlayer(projectileStartPosition, projectileEndPosition, otherPlayerPosition, playerSize)
	assert.Equal(t, &math.Point2D{X: 0.0, Y: 0.5}, impact)
}

func TestDetectImpactWithOtherPlayer5(t *testing.T) {
	projectileStartPosition := &math.Point2D{X: 0.0, Y: 0.0}
	projectileEndPosition := &math.Point2D{X: 0.0, Y: 3.0}
	otherPlayerPosition := &math.Point2D{X: 0.2, Y: 5.0}
	playerSize := 0.6
	impact := detectImpactWithOtherPlayer(projectileStartPosition, projectileEndPosition, otherPlayerPosition, playerSize)
	assert.Nil(t, impact)
}

func TestDetectImpactWithOtherPlayer6(t *testing.T) {
	projectileStartPosition := &math.Point2D{X: 0.0, Y: 3.0}
	projectileEndPosition := &math.Point2D{X: 0.0, Y: 0.0}
	otherPlayerPosition := &math.Point2D{X: -0.2, Y: 2.0}
	playerSize := 0.6
	impact := detectImpactWithOtherPlayer(projectileStartPosition, projectileEndPosition, otherPlayerPosition, playerSize)
	assert.Equal(t, &math.Point2D{X: 0.0, Y: 2.0}, impact)
}

func TestMoveWithNoPlayerNoWall(t *testing.T) {
	startPosition := &math.Point2D{X: 0.0, Y: 0.0}
	world := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	velocity := 1.0
	angle := 0.0
	eventPublisher := new(testeventpublisher.MockEventPublisher)
	projectile := createProjectForMoveAction(
		"projectileID",
		startPosition,
		velocity,
		angle,
		map[string]*math.Point2D{
			"otherPlayerID": {X: 2.0, Y: 0.0},
		},
		0.6,
		world,
		mathHelper,
		eventPublisher)
	var wallImpact *math.Point2D
	mathHelper.On("CastRay", startPosition, world, angle, velocity).Return(wallImpact)
	world.On("GetCellValue", 1, 0).Return(0)
	projectile.Move()
	mock.AssertExpectationsForObjects(t, mathHelper, world, eventPublisher)
}

func TestMoveWithOnePlayerNoWall(t *testing.T) {
	startPosition := &math.Point2D{X: 0.0, Y: 0.0}
	world := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	velocity := 2.0
	angle := 0.0
	eventPublisher := new(testeventpublisher.MockEventPublisher)
	projectile := createProjectForMoveAction(
		"projectileID",
		startPosition,
		velocity,
		angle,
		map[string]*math.Point2D{
			"otherPlayerID": {X: 2.0, Y: 0.0},
		},
		0.6,
		world,
		mathHelper,
		eventPublisher)
	var wallImpact *math.Point2D
	mathHelper.On("CastRay", startPosition, world, angle, velocity).Return(wallImpact)
	eventPublisher.On(
		"PublishEvent",
		mock.MatchedBy(
			func(ev event.Event) bool {
				return ev.Action == "projectilePlayerImpact" && ev.ExtraData["playerID"] == "otherPlayerID"
			},
		),
	)
	projectile.Move()
	mock.AssertExpectationsForObjects(t, mathHelper, world, eventPublisher)
}

func TestMoveWithOnePlayerWithWall(t *testing.T) {
	startPosition := &math.Point2D{X: 0.0, Y: 0.0}
	world := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	velocity := 2.0
	angle := 0.0
	eventPublisher := new(testeventpublisher.MockEventPublisher)
	projectile := createProjectForMoveAction(
		"projectileID",
		startPosition,
		velocity,
		angle,
		map[string]*math.Point2D{
			"otherPlayerID": {X: 1.5, Y: 0.0},
		},
		0.6,
		world,
		mathHelper,
		eventPublisher)
	wallImpact := &math.Point2D{X: 2.0, Y: 0.0}
	mathHelper.On("CastRay", startPosition, world, angle, velocity).Return(wallImpact)
	eventPublisher.On(
		"PublishEvent",
		mock.MatchedBy(
			func(ev event.Event) bool {
				return ev.Action == "projectilePlayerImpact" && ev.ExtraData["playerID"] == "otherPlayerID"
			},
		),
	)
	projectile.Move()
	mock.AssertExpectationsForObjects(t, mathHelper, world, eventPublisher)
}

func TestMoveWithOnePlayerBehindWall(t *testing.T) {
	startPosition := &math.Point2D{X: 0.0, Y: 0.0}
	world := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	velocity := 2.0
	angle := 0.0
	eventPublisher := new(testeventpublisher.MockEventPublisher)
	projectile := createProjectForMoveAction(
		"projectileID",
		startPosition,
		velocity,
		angle,
		map[string]*math.Point2D{
			"otherPlayerID": {X: 2.0, Y: 0.0},
		},
		0.6,
		world,
		mathHelper,
		eventPublisher)
	wallImpact := &math.Point2D{X: 1.0, Y: 0.0}
	mathHelper.On("CastRay", startPosition, world, angle, velocity).Return(wallImpact)
	eventPublisher.On(
		"PublishEvent",
		mock.MatchedBy(
			func(ev event.Event) bool {
				return ev.Action == "projectileWallImpact"
			},
		),
	)
	projectile.Move()
	mock.AssertExpectationsForObjects(t, mathHelper, world, eventPublisher)
}

func TestMoveWithtwoPlayersNoWall(t *testing.T) {
	world := new(testworld.MockWorldMap)
	mathHelper := new(testhelper.MockMathHelper)
	eventPublisher := new(testeventpublisher.MockEventPublisher)
	startPosition := &math.Point2D{X: 0.0, Y: 0.0}
	angle := 0.0
	velocity := 5.0
	projectile := createProjectForMoveAction(
		"projectileID",
		startPosition,
		velocity,
		angle,
		map[string]*math.Point2D{
			"otherPlayerID1": {X: 3.0, Y: 0.0},
			"otherPlayerID2": {X: 2.0, Y: 0.0},
		},
		0.6,
		world,
		mathHelper,
		eventPublisher)
	var wallImpactPosition *math.Point2D
	mathHelper.On("CastRay", startPosition, world, angle, velocity).Return(wallImpactPosition)
	eventPublisher.On(
		"PublishEvent",
		mock.MatchedBy(
			func(ev event.Event) bool {
				return ev.Action == "projectilePlayerImpact" && ev.ExtraData["playerID"] == "otherPlayerID2"
			},
		),
	)
	projectile.Move()
	mock.AssertExpectationsForObjects(t, mathHelper, world, eventPublisher)
}

func createProjectForMoveAction(
	id string,
	startPosition *math.Point2D,
	velocity, angle float64,
	players map[string]*math.Point2D,
	playerSize float64,
	world world.WorldMap,
	mathHelper helper.MathHelper,
	publisher publisher.EventPublisher) *ProjectileImpl {
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	for playerID, playerPosition := range players {
		otherPlayers[playerID] = &animatedElementImpl.AnimatedElementImpl{}
		otherPlayers[playerID].SetState(&state.AnimatedElementState{
			Position: playerPosition,
			Size:     playerSize,
		})
	}
	projectileState := &state.AnimatedElementState{
		Position:      startPosition,
		Size:          playerSize,
		MoveDirection: state.Forward,
		Velocity:      velocity,
		Angle:         angle,
	}
	projectile := &ProjectileImpl{
		mathHelper:      mathHelper,
		world:           world,
		otherPlayers:    otherPlayers,
		AnimatedElement: animatedElementImpl.NewAnimatedElementWithState(id, projectileState, world, mathHelper),
		EventPublisher:  publisher,
	}
	return projectile

}
