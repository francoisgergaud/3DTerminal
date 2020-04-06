package testwebsocket

import "github.com/stretchr/testify/mock"

//MockWebsockeConnection mocks a WebsockeConnection
type MockWebsockeConnection struct {
	mock.Mock
}

//ReadJSON mcoks the method of the same name
func (wsConnecton *MockWebsockeConnection) ReadJSON(value interface{}) error {
	args := wsConnecton.Called(value)
	return args.Error(0)
}

//WriteJSON mcoks the method of the same name
func (wsConnecton *MockWebsockeConnection) WriteJSON(value interface{}) error {
	args := wsConnecton.Called(value)
	return args.Error(0)
}

//Close mcoks the method of the same name
func (wsConnecton *MockWebsockeConnection) Close() error {
	args := wsConnecton.Called()
	return args.Error(0)
}
