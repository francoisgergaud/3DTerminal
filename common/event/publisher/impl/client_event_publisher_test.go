package impl

import (
	"francoisgergaud/3dGame/common/event"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublishEvent(t *testing.T) {
	eventPublisher := NewEventPublisherImpl()
	listener := make(chan event.Event)
	eventPublisher.RegisterListener(listener)
	eventPublished := event.Event{}
	go eventPublisher.PublishEvent(eventPublished)
	eventReceived := <-listener
	assert.Equal(t, eventReceived, eventPublished)
}
