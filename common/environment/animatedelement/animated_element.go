package animatedelement

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
)

//AnimatedElement is the interface any animated-element should implement
type AnimatedElement interface {
	Move()
	State() *state.AnimatedElementState
	SetState(state *state.AnimatedElementState)
	ID() string
}
