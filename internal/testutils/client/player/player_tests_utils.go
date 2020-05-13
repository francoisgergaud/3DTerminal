package testplayer

import (
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"
	"time"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/mock"
)

//MockPlayerFactory
type MockPlayerFactory struct {
	mock.Mock
}

//NewPlayer mocks the player factory
func (mock *MockPlayerFactory) NewPlayer(playerState *state.AnimatedElementState, world world.WorldMap, mathHelper helper.MathHelper, quit chan struct{}) player.Player {
	args := mock.Called(playerState, world, mathHelper, quit)
	return args.Get(0).(player.Player)
}

//MockPlayer mocks a character
type MockPlayer struct {
	mock.Mock
	UpdateChannel chan time.Time
	QuitChannel   chan struct{}
}

//GetPosition returns the player's position.
func (mock *MockPlayer) GetPosition() *math.Point2D {
	args := mock.Called()
	return args.Get(0).(*math.Point2D)
}

//GetAngle returns the player's orientation angle.
func (mock *MockPlayer) GetAngle() float64 {
	args := mock.Called()
	return args.Get(0).(float64)
}

//Action mocks the method of the same name
func (mock *MockPlayer) Action(eventKey *tcell.EventKey) {
	mock.Called(eventKey)
}

//Start mocks the method of the same name
func (mock *MockPlayer) Start() {
	mock.Called()
	go func() {
		for {
			select {
			case <-mock.UpdateChannel:
			case <-mock.QuitChannel:
				break
			}
		}
	}()
}

//Move mocks the method of the same name
func (mock *MockPlayer) Move() {
	mock.Called()
}

//GetUpdateChannel mocks the method of the same name
func (mock *MockPlayer) GetUpdateChannel() chan time.Time {
	mock.Called()
	return mock.UpdateChannel
}

//GetQuitChannel mocks the method of the same name
func (mock *MockPlayer) GetQuitChannel() chan struct{} {
	mock.Called()
	return mock.QuitChannel
}

//GetState mocks the method of the same name
func (mock *MockPlayer) GetState() *state.AnimatedElementState {
	args := mock.Called()
	return args.Get(0).(*state.AnimatedElementState)
}

//SetState mocks the method of the same name
func (mock *MockPlayer) SetState(state *state.AnimatedElementState) {
	mock.Called(state)
}

//GetID mocks the method of the same name
func (mock *MockPlayer) GetID() string {
	args := mock.Called()
	return args.String(0)
}

//PublishEvent mocks the method of the same name
func (mock *MockPlayer) PublishEvent(event event.Event) {
	mock.Called(event)
}

//RegisterListener mocks the method of the same name
func (mock *MockPlayer) RegisterListener(eventChanel chan<- event.Event) {
	mock.Called(eventChanel)
}
