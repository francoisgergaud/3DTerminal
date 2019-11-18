package worldelement

import (
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/environment/world"
	"github.com/gdamore/tcell"
	"math"
	"time"
)

//WorldElement is a world-element with a position and a size. The shape and overall rendering and behavior
//will be defined by the implementation.
type WorldElement interface {
	GetPosition() *common.Point2D
	GetSize() float64
	GetStyle() tcell.Style
	Start()
	GetUpdateChannel() chan time.Time
	GetQuitChannel() chan struct{}
}

//NewWorldElementImpl build a new world-element implementation.
func NewWorldElementImpl(position *common.Point2D, angle, velocity, size float64, style tcell.Style, world world.WorldMap, mathHelper common.MathHelper) WorldElement {
	return &WorldElementImpl{
		position:      position,
		size:          size,
		style:         style,
		updateChannel: make(chan time.Time),
		quitChannel:   make(chan struct{}),
		velocity:      velocity,
		angle:         angle,
		world:         world,
		mathHelper:    mathHelper,
	}
}

//WorldElementImpl is a world-element implementation which renders as a diamond.
type WorldElementImpl struct {
	position      *common.Point2D
	size          float64
	angle         float64
	velocity      float64
	world         world.WorldMap
	style         tcell.Style
	updateChannel chan time.Time
	quitChannel   chan struct{}
	mathHelper    common.MathHelper
}

//GetPosition returns the position of the world-element.
func (worldElement *WorldElementImpl) GetPosition() *common.Point2D {
	return worldElement.position
}

//GetSize returns the size of the world-element.
func (worldElement *WorldElementImpl) GetSize() float64 {
	return worldElement.size
}

//GetStyle returns the style of the world-element.
func (worldElement *WorldElementImpl) GetStyle() tcell.Style {
	return worldElement.style
}

//GetUpdateChannel returns the channel used to listen to 'update' event.
func (worldElement *WorldElementImpl) GetUpdateChannel() chan time.Time {
	return worldElement.updateChannel
}

//GetQuitChannel returns the channel used to listen to 'quit' event.
func (worldElement *WorldElementImpl) GetQuitChannel() chan struct{} {
	return worldElement.quitChannel
}

//Start triggers the animation of the world-element.
func (worldElement *WorldElementImpl) Start() {
	go func() {
		for {
			select {
			case <-worldElement.updateChannel:
				worldElement.move()
			case <-worldElement.quitChannel:
				return
			}
		}
	}()
}

// Update the world-element's position depending on the colision of walls
func (worldElement *WorldElementImpl) move() {
	rayDestination := worldElement.mathHelper.CastRay(worldElement.position, worldElement.world, worldElement.angle, worldElement.velocity)
	if rayDestination == nil {
		worldElement.position.X = worldElement.position.X + math.Cos(worldElement.angle*math.Pi)*worldElement.velocity
		worldElement.position.Y = worldElement.position.Y + math.Sin(worldElement.angle*math.Pi)*worldElement.velocity
	} else {
		//horizontal rebound
		if rayDestination.X-math.Floor(rayDestination.X) < 0.0001 && rayDestination.X-math.Floor(rayDestination.X) > -0.0001 {
			switch {
			case worldElement.angle <= 1.0:
				worldElement.angle = 1.0 - worldElement.angle
			case worldElement.angle <= 2.0:
				worldElement.angle = 3.0 - worldElement.angle
			}
		}
		//vertical rebound
		if rayDestination.Y-math.Floor(rayDestination.Y) < 0.0001 && rayDestination.Y-math.Floor(rayDestination.Y) > -0.0001 {
			switch {
			case worldElement.angle <= 1.0:
				worldElement.angle = 2.0 - worldElement.angle
			case worldElement.angle <= 2.0:
				worldElement.angle = 2.0 - worldElement.angle
			}
		}
		distanceToWall := worldElement.position.Distance(rayDestination)
		worldElement.position.X = rayDestination.X + math.Cos(worldElement.angle*math.Pi)*(worldElement.velocity-distanceToWall)
		worldElement.position.Y = rayDestination.Y + math.Sin(worldElement.angle*math.Pi)*(worldElement.velocity-distanceToWall)
	}
}
