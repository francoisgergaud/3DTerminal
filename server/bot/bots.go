package bot

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/event/publisher"
)

//Bot is a world-element with a position and a size. The shape and overall rendering and behavior
//will be defined by the implementation.
type Bot interface {
	animatedelement.AnimatedElement
	publisher.EventPublisher
}
