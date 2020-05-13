package testrunner

import (
	"francoisgergaud/3dGame/common/runner"

	"github.com/stretchr/testify/mock"
)

//MockRunner mocks the runner interface
type MockRunner struct {
	mock.Mock
}

//Start mocks the method of the same name
func (runner *MockRunner) Start(runnable runner.Runnable) {
	runner.Called(runnable)
}
