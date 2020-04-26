package impl

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/event"
	testanimatedelement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	"francoisgergaud/3dGame/server/bot"
	"francoisgergaud/3dGame/server/connector"
	"testing"

	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testconnector "francoisgergaud/3dGame/internal/testutils/server/connector"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewServer(t *testing.T) {
	quit := make(chan struct{})
	worldUpdateRate := 3
	server, error := NewServer(worldUpdateRate, quit)
	assert.Nil(t, error)
	serverImpl, ok := server.(*Impl)
	assert.True(t, ok)
	assert.Len(t, serverImpl.bots, 1)
	assert.Len(t, serverImpl.players, 0)
	assert.Equal(t, worldUpdateRate, serverImpl.botsUpdateRate)
}

func TestRegisterPlayer(t *testing.T) {
	quit := make(chan struct{})
	clientConnections := make(map[string]connector.ClientConnection)
	worldMap := new(testworld.MockWorldMap)
	serverPlayers := make(map[string]animatedelement.AnimatedElement)
	bots := make(map[string]bot.Bot)
	timeFrame := uint32(0)
	eventQueue := make(chan event.Event, 100)
	clientUpdateRate := 2
	botsUpdateRate := 3
	server := Impl{
		clientConnections: clientConnections,
		worldMap:          worldMap,
		players:           serverPlayers,
		bots:              bots,
		timeFrame:         timeFrame,
		eventQueue:        eventQueue,
		quit:              quit,
		clientUpdateRate:  clientUpdateRate,
		botsUpdateRate:    botsUpdateRate,
	}
	clientConnection := new(testconnector.MockClientConnection)
	worldMap.On("Clone").Return(worldMap)
	var eventCapture []event.Event
	clientConnection.On(
		"SendEventsToClient",
		server.timeFrame,
		mock.MatchedBy(
			func(events []event.Event) bool {
				eventCapture = events
				return true
			},
		),
	)
	server.RegisterPlayer(clientConnection)
	assert.NotEmpty(t, eventCapture[0].PlayerID)
	assert.Equal(t, worldMap, eventCapture[0].ExtraData["worldMap"])
	assert.NotNil(t, eventCapture[0].State)
	assert.Empty(t, eventCapture[0].ExtraData["otherPlayers"])
	assert.NotNil(t, serverPlayers[eventCapture[0].PlayerID])
	assert.Equal(t, 1, len(server.eventQueue))
	assert.Equal(t, timeFrame, server.timeFrame)
	newPlayerEvent := <-server.eventQueue
	assert.Equal(t, newPlayerEvent.Action, "join")
	assert.NotEmpty(t, newPlayerEvent.State)

}

func TestUnregisterClient(t *testing.T) {
	quit := make(chan struct{})
	clientConnections := make(map[string]connector.ClientConnection)
	worldMap := new(testworld.MockWorldMap)
	palyers := make(map[string]animatedelement.AnimatedElement)
	clientConnection := new(testconnector.MockClientConnection)
	playerID := "playerTest"
	clientConnections[playerID] = clientConnection
	palyers[playerID] = new(testanimatedelement.MockAnimatedElement)
	bots := make(map[string]bot.Bot)
	timeFrame := uint32(0)
	eventQueue := make(chan event.Event, 100)
	clientUpdateRate := 2
	botsUpdateRate := 3
	server := Impl{
		clientConnections: clientConnections,
		worldMap:          worldMap,
		players:           palyers,
		bots:              bots,
		timeFrame:         timeFrame,
		eventQueue:        eventQueue,
		quit:              quit,
		clientUpdateRate:  clientUpdateRate,
		botsUpdateRate:    botsUpdateRate,
	}
	server.UnregisterClient(playerID)
	assert.Equal(t, 1, len(server.eventQueue))
	assert.NotContains(t, playerID, server.clientConnections)
	assert.NotContains(t, playerID, server.players)
}

func TestReceiveEventFromClient(t *testing.T) {
	quit := make(chan struct{})
	clientConnections := make(map[string]connector.ClientConnection)
	worldMap := new(testworld.MockWorldMap)
	palyers := make(map[string]animatedelement.AnimatedElement)
	clientConnection := new(testconnector.MockClientConnection)
	playerID := "playerTest"
	clientConnections[playerID] = clientConnection
	mockAnimatedElement := new(testanimatedelement.MockAnimatedElement)
	palyers[playerID] = mockAnimatedElement
	bots := make(map[string]bot.Bot)
	timeFrame := uint32(0)
	eventQueue := make(chan event.Event, 100)
	clientUpdateRate := 2
	botsUpdateRate := 3
	server := Impl{
		clientConnections: clientConnections,
		worldMap:          worldMap,
		players:           palyers,
		bots:              bots,
		timeFrame:         timeFrame,
		eventQueue:        eventQueue,
		quit:              quit,
		clientUpdateRate:  clientUpdateRate,
		botsUpdateRate:    botsUpdateRate,
	}
	newPlayerState := state.AnimatedElementState{}
	eventFromPlayer := event.Event{
		Action:   "move",
		PlayerID: playerID,
		State:    &newPlayerState,
	}
	mockAnimatedElement.On("SetState", &newPlayerState)
	server.ReceiveEventFromClient(eventFromPlayer)
	assert.Equal(t, 1, len(server.eventQueue))
}

func TestSendEventsToClients(t *testing.T) {
	quit := make(chan struct{})
	clientConnections := make(map[string]connector.ClientConnection)
	worldMap := new(testworld.MockWorldMap)
	palyers := make(map[string]animatedelement.AnimatedElement)
	clientConnection := new(testconnector.MockClientConnection)
	playerID := "playerTest"
	clientConnections[playerID] = clientConnection
	palyers[playerID] = new(testanimatedelement.MockAnimatedElement)
	bots := make(map[string]bot.Bot)
	timeFrame := uint32(2)
	eventQueue := make(chan event.Event, 100)
	clientUpdateRate := 2
	botsUpdateRate := 3
	server := Impl{
		clientConnections: clientConnections,
		worldMap:          worldMap,
		players:           palyers,
		bots:              bots,
		timeFrame:         timeFrame,
		eventQueue:        eventQueue,
		quit:              quit,
		clientUpdateRate:  clientUpdateRate,
		botsUpdateRate:    botsUpdateRate,
	}
	//create an event
	newPlayerState := state.AnimatedElementState{}
	eventFromPlayer := event.Event{
		Action:   "move",
		PlayerID: "",
		State:    &newPlayerState,
	}
	eventQueue <- eventFromPlayer
	events := make([]event.Event, 1)
	events[0] = eventFromPlayer
	var eventCapture []event.Event
	clientConnection.On("SendEventsToClient", timeFrame, mock.MatchedBy(
		func(events []event.Event) bool {
			eventCapture = events
			return true
		},
	))
	server.sendEventsToClients()
	assert.Equal(t, timeFrame, eventCapture[0].TimeFrame)
	assert.Equal(t, &newPlayerState, eventCapture[0].State)
	assert.Equal(t, "move", eventCapture[0].Action)
}
