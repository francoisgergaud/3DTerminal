package testutils

import (
	"github.com/stretchr/testify/mock"
)

//MockConsoleEventManager mocks the calls to the ConsoleEventManager interface.
type MockConsoleEventManager struct {
	mock.Mock
}

//Listen mocks the call to the the method of the same name.
func (mock *MockConsoleEventManager) Listen() {
	mock.Called()
}
