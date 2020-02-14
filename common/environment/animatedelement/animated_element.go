package animatedelement

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"time"
)

//AnimatedElement is the interface any animated-element should implement
type AnimatedElement interface {
	GetUpdateChannel() chan time.Time
	Start()
	Move()
	GetState() *state.AnimatedElementState
	SetState(state *state.AnimatedElementState)
	GetID() string
}
