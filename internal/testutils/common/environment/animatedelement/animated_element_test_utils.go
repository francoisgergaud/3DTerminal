package testanimatedelement

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/math/helper"

	"github.com/stretchr/testify/mock"
)

//MockAnimatedElementFactory mocks an animated-element factory
type MockAnimatedElementFactory struct {
	mock.Mock
}

//NewAnimatedElementWithState mocks the creation of new animated-element
func (mock *MockAnimatedElementFactory) NewAnimatedElementWithState(id string, animatedElementState *state.AnimatedElementState, world world.WorldMap, mathHelper helper.MathHelper) animatedelement.AnimatedElement {
	args := mock.Called(id, animatedElementState, world, mathHelper)
	return args.Get(0).(animatedelement.AnimatedElement)
}

//MockAnimatedElement mocks a world-element.
type MockAnimatedElement struct {
	mock.Mock
}

//Move mocks a method of the same name from a world-element.
func (mock *MockAnimatedElement) Move() {
	mock.Called()
}

//State mocks the animated-element's state.
func (mock *MockAnimatedElement) State() *state.AnimatedElementState {
	args := mock.Called()
	return args.Get(0).(*state.AnimatedElementState)
}

//SetState mocks the animated-element's state.
func (mock *MockAnimatedElement) SetState(state *state.AnimatedElementState) {
	mock.Called(state)
}

//ID mocks the animated-element's state.
func (mock *MockAnimatedElement) ID() string {
	args := mock.Called()
	return args.String(0)
}
