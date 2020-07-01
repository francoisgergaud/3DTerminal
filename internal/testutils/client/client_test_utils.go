package testclient

import (
	"francoisgergaud/3dGame/client/connector"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/projectile"
	"francoisgergaud/3dGame/common/event"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/mock"
)

//MockEngine is a mocked engine
type MockEngine struct {
	mock.Mock
}

//ReceiveEventsFromServer mocks the method of the same name
func (mock *MockEngine) ReceiveEventsFromServer(events []event.Event) {
	mock.Called(events)
}

//Player mocks the method of the same name
func (mock *MockEngine) Player() animatedelement.AnimatedElement {
	args := mock.Called()
	return args.Get(0).(animatedelement.AnimatedElement)
}

//Shutdown mocks the method of the same name
func (mock *MockEngine) Shutdown() {
	mock.Called()
}

//ConnectToServer mocks the method of the same name
func (mock *MockEngine) ConnectToServer(connectionToServer connector.ServerConnector) {
	mock.Called(connectionToServer)
}

//OtherPlayers mocks the method of the name
func (mock *MockEngine) OtherPlayers() map[string]animatedelement.AnimatedElement {
	args := mock.Called()
	return args.Get(0).(map[string]animatedelement.AnimatedElement)
}

//Action mocks the method of the name
func (mock *MockEngine) Action(eventKey *tcell.EventKey) {
	mock.Called(eventKey)
}

//ReceiveEvent mocks the method of the name
func (mock *MockEngine) ReceiveEvent(event event.Event) {
	mock.Called(event)
}

//Projectiles mocks the method of the name
func (mock *MockEngine) Projectiles() map[string]projectile.Projectile {
	args := mock.Called()
	return args.Get(0).(map[string]projectile.Projectile)
}
