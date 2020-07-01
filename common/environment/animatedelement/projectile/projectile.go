package projectile

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	animatedElementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/event/publisher"
	eventPublisherImpl "francoisgergaud/3dGame/common/event/publisher/impl"
	"francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"
	originalMath "math"

	"github.com/gdamore/tcell"
)

//Projectile is an animated-element which has a straight path until it impacts a wall or another-player
type Projectile interface {
	animatedelement.AnimatedElement
	publisher.EventPublisher
}

//NewProjectile is a factory for projectile
func NewProjectile(id string, position *math.Point2D, angle float64, world world.WorldMap, otherPlayers map[string]animatedelement.AnimatedElement, mathHelper helper.MathHelper) Projectile {
	projectileState := &state.AnimatedElementState{
		Velocity:      0.5,
		Position:      position,
		Angle:         angle,
		Size:          0.1,
		Style:         tcell.StyleDefault.Background(tcell.ColorDarkRed),
		MoveDirection: state.Forward,
	}
	return &ProjectileImpl{
		mathHelper:      mathHelper,
		world:           world,
		otherPlayers:    otherPlayers,
		AnimatedElement: animatedElementImpl.NewAnimatedElementWithState(id, projectileState, world, mathHelper),
		EventPublisher:  eventPublisherImpl.NewEventPublisherImpl(),
	}
}

//ProjectileImpl is the default-implementation of a Projectile
//TODO: avoid collision with the shoter itself
type ProjectileImpl struct {
	animatedelement.AnimatedElement
	publisher.EventPublisher
	world        world.WorldMap
	otherPlayers map[string]animatedelement.AnimatedElement
	mathHelper   helper.MathHelper
}

//Move moves the projectile on update
func (projectile *ProjectileImpl) Move() {
	projectileState := projectile.State()
	//check impacts with wall
	rayDestination := projectile.mathHelper.CastRay(projectileState.Position, projectile.world, projectileState.Angle, projectileState.Velocity)
	var endPosition *math.Point2D
	if rayDestination != nil {
		endPosition = rayDestination
	} else {
		endPosition = &math.Point2D{
			X: projectileState.Position.X + originalMath.Cos(projectileState.Angle*originalMath.Pi)*projectileState.Velocity,
			Y: projectileState.Position.Y + originalMath.Sin(projectileState.Angle*originalMath.Pi)*projectileState.Velocity,
		}
	}
	minImpactDistance := originalMath.Inf(1)
	var closestPlayer animatedelement.AnimatedElement
	var closestPlayerID string
	//checks impacts with other-players
	for otherPlayerID, otherPlayer := range projectile.otherPlayers {
		impact := detectImpactWithOtherPlayer(projectileState.Position, endPosition, otherPlayer.State().Position, otherPlayer.State().Size)
		if impact != nil {
			distanceToImpact := projectileState.Position.Distance(impact)
			if distanceToImpact < minImpactDistance {
				minImpactDistance = distanceToImpact
				closestPlayer = otherPlayer
				closestPlayerID = otherPlayerID
			}
		}
	}
	var eventToSend event.Event
	//if impact with wall
	if rayDestination != nil {

		if closestPlayer != nil && minImpactDistance < rayDestination.Distance(projectileState.Position) {
			//if there is another player in-between the player and the wall
			projectileState.MoveDirection = state.None
			eventToSend = event.Event{
				PlayerID: projectile.ID(),
				State:    projectileState,
				Action:   "projectilePlayerImpact",
				ExtraData: map[string]interface{}{
					"playerID": closestPlayerID,
				},
			}
			projectile.PublishEvent(eventToSend)
		} else {
			//if there no other player in-between the player and the wall
			projectileState.MoveDirection = state.None
			eventToSend = event.Event{
				PlayerID: projectile.ID(),
				State:    projectileState,
				Action:   "projectileWallImpact",
			}
			projectile.PublishEvent(eventToSend)
		}
	} else {
		if closestPlayer != nil {
			// if there is no impact with a wall, but there is an impact with another player
			projectileState.MoveDirection = state.None
			eventToSend = event.Event{
				PlayerID: projectile.ID(),
				State:    projectileState,
				Action:   "projectilePlayerImpact",
				ExtraData: map[string]interface{}{
					"playerID": closestPlayerID,
				},
			}
			projectile.PublishEvent(eventToSend)
		} else {
			// if there is no impact
			projectile.AnimatedElement.Move()
		}
	}
}

//detectImpactWithOtherPlayer calculate if there is an impact with another player. IF there is, it will return the point of impact, otherwise it returns nil
func detectImpactWithOtherPlayer(startPosition, endPosition, otherPalyerPosition *math.Point2D, otherPlayerSize float64) *math.Point2D {
	if endPosition.X != startPosition.X {
		slope := (endPosition.Y - startPosition.Y) / (endPosition.X - startPosition.X)
		xToProject := otherPalyerPosition.X - startPosition.X
		yToProject := otherPalyerPosition.Y - startPosition.Y
		translatedPlayerPosition := &math.Point2D{X: xToProject, Y: yToProject}
		xProjection := (xToProject + (slope * yToProject)) / (1 + (slope * slope))
		yProjection := ((slope * xToProject) + (slope * slope * yToProject)) / (1 + (slope * slope))
		projectedPoint := &math.Point2D{X: xProjection, Y: yProjection}
		projectPointGlobal := &math.Point2D{X: xProjection + startPosition.X, Y: yProjection + startPosition.Y}
		if projectedPoint.Distance(translatedPlayerPosition) < (otherPlayerSize / 2) {
			if endPosition.X-startPosition.X > 0 {
				if xProjection >= 0 && xProjection <= (endPosition.X-startPosition.X) {
					return projectPointGlobal
				}
			} else if endPosition.X-startPosition.X < 0 {
				if xProjection <= 0 && xProjection >= (endPosition.X-startPosition.X) {
					return projectPointGlobal
				}
			}
		}
	} else {
		xProjection := startPosition.X
		yProjection := otherPalyerPosition.Y
		projectedPoint := &math.Point2D{X: xProjection, Y: yProjection}
		distance := otherPalyerPosition.X - xProjection
		if distance < 0 {
			distance = -distance
		}
		if distance < (otherPlayerSize / 2) {
			// 3 points are on the same vertical line, uses the y coordinates to determine if the projection is between start and end
			if endPosition.Y-startPosition.Y > 0 {
				if yProjection >= startPosition.Y && yProjection <= endPosition.Y {
					return projectedPoint
				}
			} else {
				if yProjection >= endPosition.Y && yProjection <= startPosition.Y {
					return projectedPoint
				}
			}
		}
	}
	return nil
}
