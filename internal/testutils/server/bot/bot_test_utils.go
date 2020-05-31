package testbot

import (
	testanimatedelement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	testeventpublisher "francoisgergaud/3dGame/internal/testutils/common/event/publisher"
)

//MockBot mocks a bot
type MockBot struct {
	testeventpublisher.MockEventPublisher
	testanimatedelement.MockAnimatedElement
}
