package character

import (
	"francoisgergaud/3dGame/common"
	"time"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/mock"
)

//MockCharacter mocks a character
type MockCharacter struct {
	mock.Mock
	UpdateChannel chan time.Time
	QuitChannel   chan struct{}
}

//GetPosition returns the player's position.
func (mock *MockCharacter) GetPosition() *common.Point2D {
	args := mock.Called()
	return args.Get(0).(*common.Point2D)
}

//GetAngle returns the player's orientation angle.
func (mock *MockCharacter) GetAngle() float64 {
	args := mock.Called()
	return args.Get(0).(float64)
}

//Action mocks the operation.
func (mock *MockCharacter) Action(eventKey *tcell.EventKey) {
	mock.Called(eventKey)
}

//Start mocks the operation.
func (mock *MockCharacter) Start() {
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

//GetUpdateChannel mocks the operation.
func (mock *MockCharacter) GetUpdateChannel() chan time.Time {
	mock.Called()
	return mock.UpdateChannel
}

//GetQuitChannel mocks the operation.
func (mock *MockCharacter) GetQuitChannel() chan struct{} {
	mock.Called()
	return mock.QuitChannel
}
