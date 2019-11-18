package main

import (
	"francoisgergaud/3dGame/internal/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestInitGame(t *testing.T) {
	screen := new(testutils.MockScreen)
	game, err := InitGame(screen)
	assert.Nil(t, err)
	assert.NotNil(t, game)
}
