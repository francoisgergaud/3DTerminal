package testanimatedelement

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/math/helper"
	"time"

	"github.com/stretchr/testify/mock"
)

//MockAnimatedElementFactory mocks an animated-element factory
type MockAnimatedElementFactory struct {
	mock.Mock
}

//NewAnimatedElementWithState mocks the creation of new animated-element
func (mock *MockAnimatedElementFactory) NewAnimatedElementWithState(animatedElementState *state.AnimatedElementState, world world.WorldMap, mathHelper helper.MathHelper, quit chan struct{}) animatedelement.AnimatedElement {
	args := mock.Called(animatedElementState, world, mathHelper, quit)
	return args.Get(0).(animatedelement.AnimatedElement)
}

//MockAnimatedElement mocks a world-element.
type MockAnimatedElement struct {
	mock.Mock
	UpdateChannel chan time.Time
	QuitChannel   chan struct{}
}

//Start mocks a method of the same name from a world-element.
func (mock *MockAnimatedElement) Start() {
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

//Move mocks a method of the same name from a world-element.
func (mock *MockAnimatedElement) Move() {
	mock.Called()
}

//GetUpdateChannel mocks a method of the same name from a world-element.
func (mock *MockAnimatedElement) GetUpdateChannel() chan time.Time {
	mock.Called()
	return mock.UpdateChannel
}

//GetQuitChannel mocks a method of the same name from a world-element.
func (mock *MockAnimatedElement) GetQuitChannel() chan struct{} {
	mock.Called()
	return mock.QuitChannel
}

//GetState mocks the animated-element's state.
func (mock *MockAnimatedElement) GetState() *state.AnimatedElementState {
	args := mock.Called()
	return args.Get(0).(*state.AnimatedElementState)
}

//SetState mocks the animated-element's state.
func (mock *MockAnimatedElement) SetState(state *state.AnimatedElementState) {
	mock.Called(state)
}

//GetID mocks the identifier
func (mock *MockAnimatedElement) GetID() string {
	args := mock.Called()
	return args.String(0)
}
