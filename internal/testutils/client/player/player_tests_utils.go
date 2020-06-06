package testplayer

import (
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/math/helper"
	testanimatedelement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	testeventpublisher "francoisgergaud/3dGame/internal/testutils/common/event/publisher"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/mock"
)

//MockPlayerFactory
type MockPlayerFactory struct {
	mock.Mock
}

//NewPlayer mocks the player factory
func (mock *MockPlayerFactory) NewPlayer(playerState *state.AnimatedElementState, world world.WorldMap, mathHelper helper.MathHelper) player.Player {
	args := mock.Called(playerState, world, mathHelper)
	return args.Get(0).(player.Player)
}

//MockPlayer mocks a character
type MockPlayer struct {
	mock.Mock
	testanimatedelement.MockAnimatedElement
	testeventpublisher.MockEventPublisher
}

//Action mocks the method of the same name
func (mock *MockPlayer) Action(eventKey *tcell.EventKey) {
	mock.Called(eventKey)
}
