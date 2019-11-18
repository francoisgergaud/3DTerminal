package character

import (
	"francoisgergaud/3dGame/common"

	"github.com/stretchr/testify/mock"
)

type MockCharacter struct {
	mock.Mock
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
