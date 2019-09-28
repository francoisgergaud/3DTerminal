package environment

import (
	"francoisgergaud/3dGame/common"
	"math"

	"github.com/gdamore/tcell"
)

//Character is an actionable character in the environment.
type Character interface {
	Action(eventKey *tcell.EventKey)
	GetPosition() *common.Point2D
	GetAngle() float64
}

//NewPlayer builds a new player from ithe input parameters.
func NewPlayer(initialPosition *common.Point2D, initialAngle float64, initialVelocity float64, world WorldMap) Character {
	return &Player{
		pos:      initialPosition,
		angle:    initialAngle,
		velocity: initialVelocity,
		world:    world,
	}
}

// Player represents the player
type Player struct {
	pos      *common.Point2D
	angle    float64
	velocity float64
	world    WorldMap
}

// Action the player according to the input key
func (p *Player) Action(eventKey *tcell.EventKey) {
	stepAngle := 0.05
	switch eventKey.Key() {
	case tcell.KeyUp:
		p.move(true)
	case tcell.KeyDown:
		p.move(false)
	case tcell.KeyLeft:
		p.angle = p.angle - stepAngle
		if p.angle < 0 {
			p.angle += 2
		}
	case tcell.KeyRight:
		p.angle = p.angle + stepAngle
		if p.angle > 2 {
			p.angle -= 2
		}
	}
}

//GetPosition returns the player's position.
func (p *Player) GetPosition() *common.Point2D {
	return p.pos
}

//GetAngle returns the player's orientation angle.
func (p *Player) GetAngle() float64 {
	return p.angle
}

// Move the player forward ('backForward' true) or backward ('backForward' false)
func (p *Player) move(backForward bool) {
	var newX, newY float64
	if backForward {
		newX = p.pos.X + math.Cos(p.angle*math.Pi)*p.velocity
		newY = p.pos.Y + math.Sin(p.angle*math.Pi)*p.velocity
	} else {
		newX = p.pos.X - math.Cos(p.angle*math.Pi)*p.velocity
		newY = p.pos.Y - math.Sin(p.angle*math.Pi)*p.velocity
	}
	if p.world.GetCellValue(int(newX), int(newY)) == 0 {
		p.pos.X = newX
		p.pos.Y = newY
	}
}
