package impl

import (
	"francoisgergaud/3dGame/common/event"
	testeventpublisher "francoisgergaud/3dGame/internal/testutils/common/event/publisher"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestPublishEvent(t *testing.T) {
	eventPublisher := NewEventPublisherImpl()
	listener := new(testeventpublisher.MockEventListener)
	eventPublisher.RegisterListener(listener)
	eventPublished := event.Event{}
	listener.On("ReceiveEvent", eventPublished)
	eventPublisher.PublishEvent(eventPublished)
	mock.AssertExpectationsForObjects(t, listener)
}
