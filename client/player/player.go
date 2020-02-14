package player

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/event/publisher"

	"github.com/gdamore/tcell"
)

//Player is an actionable character in the environment.
type Player interface {
	animatedelement.AnimatedElement
	publisher.EventPublisher
	Action(eventKey *tcell.EventKey)
}
