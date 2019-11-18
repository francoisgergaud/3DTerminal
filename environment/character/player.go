package character

import (
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/environment/world"
	"github.com/gdamore/tcell"
	"math"
	"time"
)

//Direction is the direction type.
type Direction uint

//The Direction possible values.
const (
	None Direction = iota
	Left
	Right
	Forward
	Backward
)

//Character is an actionable character in the environment.
type Character interface {
	Action(eventKey *tcell.EventKey)
	GetPosition() *common.Point2D
	GetAngle() float64
	Start()
	GetUpdateChannel() chan time.Time
	GetQuitChannel() chan struct{}
}

//NewPlayableCharacter builds a new player from ithe input parameters.
func NewPlayableCharacter(initialPosition *common.Point2D, initialAngle, velocity, stepAngle float64, world world.WorldMap) Character {
	return &PlayableCharacter{
		pos:           initialPosition,
		angle:         initialAngle,
		velocity:      velocity,
		stepAngle:     stepAngle,
		world:         world,
		updateChannel: make(chan time.Time),
		quitChannel:   make(chan struct{}),
	}
}

// PlayableCharacter represents the player
type PlayableCharacter struct {
	pos             *common.Point2D
	angle           float64
	velocity        float64
	stepAngle       float64
	world           world.WorldMap
	moveDirection   Direction
	rotateDirection Direction
	updateChannel   chan time.Time
	quitChannel     chan struct{}
}

// Action the player according to the input key
func (p *PlayableCharacter) Action(eventKey *tcell.EventKey) {
	switch eventKey.Key() {
	case tcell.KeyUp:
		if p.moveDirection == Backward {
			p.moveDirection = None
		} else {
			p.moveDirection = Forward
		}
	case tcell.KeyDown:
		if p.moveDirection == Forward {
			p.moveDirection = None
		} else {
			p.moveDirection = Backward
		}
	case tcell.KeyLeft:
		if p.rotateDirection == Right {
			p.rotateDirection = None
		} else {
			p.rotateDirection = Left
		}
	case tcell.KeyRight:
		if p.rotateDirection == Left {
			p.rotateDirection = None
		} else {
			p.rotateDirection = Right
		}
	}
}

//GetPosition returns the player's position.
func (p *PlayableCharacter) GetPosition() *common.Point2D {
	return p.pos
}

//GetAngle returns the player's orientation angle.
func (p *PlayableCharacter) GetAngle() float64 {
	return p.angle
}

//GetUpdateChannel returns the channel used to listen to 'update' event.
func (p *PlayableCharacter) GetUpdateChannel() chan time.Time {
	return p.updateChannel
}

//GetQuitChannel returns the channel used to listen to 'quit' event.
func (p *PlayableCharacter) GetQuitChannel() chan struct{} {
	return p.quitChannel
}

//Start triggers the goroutine with the loop for player updates. The loop break when a message is received
//on the quit-channel.
func (p *PlayableCharacter) Start() {
	go func() {
		for {
			select {
			case <-p.updateChannel:
				p.move()
			case <-p.quitChannel:
				return
			}
		}
	}()
}

// Update the player's position depending on its moving and rotate Direction and the cell's value on the world-map
func (p *PlayableCharacter) move() {
	if p.rotateDirection == Left {
		p.rotateDirection = Left
		p.angle = p.angle - p.stepAngle
		if p.angle < 0 {
			p.angle += 2
		}
	} else if p.rotateDirection == Right {
		p.angle = p.angle + p.stepAngle
		if p.angle >= 2 {
			p.angle -= 2
		}
	}
	if p.moveDirection != None {
		newX := p.pos.X
		newY := p.pos.Y
		if p.moveDirection == Forward {
			newX = p.pos.X + math.Cos(p.angle*math.Pi)*p.velocity
			newY = p.pos.Y + math.Sin(p.angle*math.Pi)*p.velocity
		} else if p.moveDirection == Backward {
			newX = p.pos.X - math.Cos(p.angle*math.Pi)*p.velocity
			newY = p.pos.Y - math.Sin(p.angle*math.Pi)*p.velocity
		}
		if p.world.GetCellValue(int(newX), int(newY)) == 0 {
			p.pos.X = newX
			p.pos.Y = newY
		}
	}
}
