package impl

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server/bot"
	"francoisgergaud/3dGame/server/connector"
	"testing"

	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testconnector "francoisgergaud/3dGame/internal/testutils/server/connector"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	quit := make(chan struct{})
	server, error := NewServer(quit)
	assert.Nil(t, error)
	serverImpl, ok := server.(*Impl)
	assert.True(t, ok)
	assert.Len(t, serverImpl.bots, 1)
	assert.Len(t, serverImpl.players, 0)
}

func TestRegisterPlayer(t *testing.T) {
	quit := make(chan struct{})
	clientConnections := make(map[string]connector.ClientConnection)
	worldMap := new(testworld.MockWorldMap)
	palyers := make(map[string]*state.AnimatedElementState)
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
	clientConnection := new(testconnector.MockClientConnection)
	worldMap.On("Clone").Return(worldMap)
	playerID, worldMapForPlayer, state, otherPlayers := server.RegisterPlayer(clientConnection)
	assert.NotEmpty(t, playerID)
	assert.Equal(t, worldMapForPlayer, worldMap)
	assert.NotNil(t, state)
	assert.NotEmpty(t, otherPlayers)
	assert.Equal(t, 1, len(server.eventQueue))
	assert.Equal(t, clientConnection, server.clientConnections[playerID])
	assert.Equal(t, state, *server.players[playerID])
	assert.Equal(t, worldMap, server.worldMap)
}

func TestUnregisterClient(t *testing.T) {
	quit := make(chan struct{})
	clientConnections := make(map[string]connector.ClientConnection)
	worldMap := new(testworld.MockWorldMap)
	palyers := make(map[string]*state.AnimatedElementState)
	clientConnection := new(testconnector.MockClientConnection)
	playerID := "playerTest"
	clientConnections[playerID] = clientConnection
	palyers[playerID] = new(state.AnimatedElementState)
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

func TestReceiveEventsFromClient(t *testing.T) {
	quit := make(chan struct{})
	clientConnections := make(map[string]connector.ClientConnection)
	worldMap := new(testworld.MockWorldMap)
	palyers := make(map[string]*state.AnimatedElementState)
	clientConnection := new(testconnector.MockClientConnection)
	playerID := "playerTest"
	clientConnections[playerID] = clientConnection
	palyers[playerID] = new(state.AnimatedElementState)
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
	eventsFromPlayer := make([]event.Event, 0)
	eventsFromPlayer = append(eventsFromPlayer, eventFromPlayer)
	server.ReceiveEventsFromClient(eventsFromPlayer)
	assert.Equal(t, newPlayerState, *server.players[playerID])
	assert.Equal(t, 1, len(server.eventQueue))
}

func TestSendEventsToClients(t *testing.T) {
	quit := make(chan struct{})
	clientConnections := make(map[string]connector.ClientConnection)
	worldMap := new(testworld.MockWorldMap)
	palyers := make(map[string]*state.AnimatedElementState)
	clientConnection := new(testconnector.MockClientConnection)
	playerID := "playerTest"
	clientConnections[playerID] = clientConnection
	palyers[playerID] = new(state.AnimatedElementState)
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
		PlayerID: "PlayerID",
		State:    &newPlayerState,
	}
	eventQueue <- eventFromPlayer
	events := make([]event.Event, 1)
	events[0] = eventFromPlayer
	clientConnection.On("SendEventsToClient", timeFrame, events)
	server.sendEventsToClients()
	clientConnection.AssertCalled(t, "SendEventsToClient", timeFrame, events)
}
