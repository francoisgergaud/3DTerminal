package testprojectile

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/projectile"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"
	testanimatedelement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	testeventpublisher "francoisgergaud/3dGame/internal/testutils/common/event/publisher"

	"github.com/stretchr/testify/mock"
)

//MockProjectileFactory provides feature to mock the projectile factories
type MockProjectileFactory struct {
	mock.Mock
}

//CreateProjectile mock the factory
func (factory *MockProjectileFactory) CreateProjectile(id string, position *math.Point2D, angle float64, world world.WorldMap, otherPlayers map[string]animatedelement.AnimatedElement, mathHelper helper.MathHelper) projectile.Projectile {
	args := factory.Called(id, position, angle, world, otherPlayers, mathHelper)
	return args.Get(0).(projectile.Projectile)
}

//MockProjectile mocks a projectile
type MockProjectile struct {
	testanimatedelement.MockAnimatedElement
	testeventpublisher.MockEventPublisher
	mock.Mock
}

//GetID mocks the method of the name
func (mock *MockProjectile) GetID() string {
	args := mock.Called()
	return args.String(0)
}
