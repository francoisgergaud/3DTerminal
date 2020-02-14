package testanimatedelement

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"time"

	"github.com/stretchr/testify/mock"
)

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
