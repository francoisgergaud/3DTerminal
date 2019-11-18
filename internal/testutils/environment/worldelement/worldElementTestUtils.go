package worldelement

import (
	"francoisgergaud/3dGame/common"
	"time"

	"github.com/gdamore/tcell"

	"github.com/stretchr/testify/mock"
)

//MockWorldElement mocks a world-element.
type MockWorldElement struct {
	mock.Mock
	UpdateChannel chan time.Time
	QuitChannel   chan struct{}
}

//GetPosition mocks a method of the same name from aworld-element.
func (mock *MockWorldElement) GetPosition() *common.Point2D {
	args := mock.Called()
	return args.Get(0).(*common.Point2D)
}

//GetSize mocks a method of the same name from aworld-element.
func (mock *MockWorldElement) GetSize() float64 {
	args := mock.Called()
	return args.Get(0).(float64)
}

//GetStyle mocks a method of the same name from a world-element.
func (mock *MockWorldElement) GetStyle() tcell.Style {
	args := mock.Called()
	return args.Get(0).(tcell.Style)
}

//Start mocks a method of the same name from a world-element.
func (mock *MockWorldElement) Start() {
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

//GetUpdateChannel mocks a method of the same name from a world-element.
func (mock *MockWorldElement) GetUpdateChannel() chan time.Time {
	mock.Called()
	return mock.UpdateChannel
}

//GetQuitChannel mocks a method of the same name from a world-element.
func (mock *MockWorldElement) GetQuitChannel() chan struct{} {
	mock.Called()
	return mock.QuitChannel
}
