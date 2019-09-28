package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockEngine struct {
	mock.Mock
}

func (mock *MockEngine) StartEngine() {
	mock.Called()
}

func TestStart(t *testing.T) {
	engine := new(MockEngine)
	engine.On("StartEngine")
	game := Game{
		engine: engine,
	}
	game.Start()
}
