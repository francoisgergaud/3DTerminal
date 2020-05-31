package impl

import (
	"fmt"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/event/publisher"
	"francoisgergaud/3dGame/common/math/helper"
	mathhelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/math/raycaster"
	"francoisgergaud/3dGame/common/runner"
	"francoisgergaud/3dGame/server/bot"
	"francoisgergaud/3dGame/server/connector"
	botgenerator "francoisgergaud/3dGame/server/impl/generator/bot"
	"francoisgergaud/3dGame/server/impl/generator/worldmap"
	"francoisgergaud/3dGame/server/impl/player"
	"time"

	"github.com/google/uuid"
)

//Impl is the default implementation for a server.
type Impl struct {
	worldMap          world.WorldMap
	players           map[string]animatedelement.AnimatedElement
	bots              map[string]bot.Bot
	quit              chan struct{}
	botsUpdateRate    int
	mathHelper        mathhelper.MathHelper
	clientEventSender clientEventSender
	runner            runner.Runner
	identifierFactory func() uuid.UUID
	worldMapFactory   func() world.WorldMap
	botFactory        func(id string, worldMap world.WorldMap, mathHelper mathhelper.MathHelper, quit chan struct{}) bot.Bot
	playerFactory     func(worldMap world.WorldMap, mathHelper helper.MathHelper, quit chan struct{}) animatedelement.AnimatedElement
}

//NewServer is a server factory
func NewServer(worldUpdateRate int, quit chan struct{}) (*Impl, error) {
	server := new(Impl)
	server.bots = make(map[string]bot.Bot)
	mathHelper, err := mathhelper.NewMathHelper(new(raycaster.RayCasterImpl))
	server.mathHelper = mathHelper
	if err != nil {
		return nil, fmt.Errorf("error while instantiating the math-helper: %w", err)
	}
	server.players = make(map[string]animatedelement.AnimatedElement)
	eventQueue := make(chan event.Event, 100)
	server.clientEventSender = &clientEventSenderImp{
		clientConnections: make(map[string]connector.ClientConnection),
		clientUpdateRate:  10,
		eventQueue:        eventQueue,
		quit:              quit,
		timeFrame:         0,
	}
	server.quit = quit
	server.botsUpdateRate = worldUpdateRate
	server.runner = &runner.AsyncRunner{}
	server.identifierFactory = uuid.New
	server.worldMapFactory = worldmap.NewWorldMap
	server.botFactory = botgenerator.NewBot
	server.playerFactory = player.NewPlayer
	return server, nil
}

//Start the server
func (server *Impl) Start() {
	//initialize the environment (world and bots)
	server.worldMap = server.worldMapFactory()
	botID := server.identifierFactory().String()
	server.bots[botID] = server.botFactory(botID, server.worldMap, server.mathHelper, server.quit)
	server.bots[botID].RegisterListener(server.clientEventSender)
	//start the asynchronous listeners
	server.runner.Start(server.clientEventSender)
	server.runner.Start(server)
}

//RegisterPlayer register a player and provide the environment
func (server *Impl) RegisterPlayer(clientConnection connector.ClientConnection) string {
	playerID := server.identifierFactory().String()
	server.clientEventSender.AddClient(playerID, clientConnection)
	player := server.playerFactory(server.worldMap, server.mathHelper, server.quit)
	server.players[playerID] = player
	newPlayerEvent := event.Event{
		PlayerID: playerID,
		State:    player.GetState(),
		Action:   "join",
	}
	worldMap := server.worldMap.Clone()
	server.clientEventSender.ReceiveEvent(newPlayerEvent)
	otherPlayers := make(map[string]state.AnimatedElementState)
	for id, player := range server.players {
		if id != playerID {
			otherPlayers[id] = player.GetState().Clone()
		}
	}
	for id, bot := range server.bots {
		otherPlayers[id] = bot.GetState().Clone()
	}
	extraData := make(map[string]interface{})
	extraData["worldMap"] = worldMap
	extraData["otherPlayers"] = otherPlayers
	newPlayerInitializationEvent := event.Event{
		Action:    "init",
		PlayerID:  playerID,
		State:     player.GetState(),
		ExtraData: extraData,
	}
	server.clientEventSender.SendEventToClient(playerID, newPlayerInitializationEvent)
	server.players[playerID].Start()
	return playerID
}

//UnregisterClient removes a player
func (server *Impl) UnregisterClient(playerID string) {
	delete(server.players, playerID)
	server.clientEventSender.RemoveClient(playerID)
	event := event.Event{
		PlayerID: playerID,
		Action:   "quit",
	}
	server.clientEventSender.ReceiveEvent(event)
}

//ReceiveEventFromClient manage an event received from a client
// as it is supposed to override the previous ones
func (server *Impl) ReceiveEventFromClient(event event.Event) {
	server.players[event.PlayerID].SetState(event.State)
	server.clientEventSender.ReceiveEvent(event)
}

//Run is a blocking loop using a ticket to update the environment
func (server *Impl) Run() {
	botsTicker := time.NewTicker(time.Duration(1000/server.botsUpdateRate) * time.Millisecond)
	for {
		select {
		case <-server.quit:
			botsTicker.Stop()
			return
		case <-botsTicker.C:
			for _, bot := range server.bots {
				bot.Move()
			}
			for _, player := range server.players {
				player.Move()
			}
		}
	}
}

type clientEventSender interface {
	runner.Runnable
	AddClient(playerID string, connectionToClient connector.ClientConnection)
	RemoveClient(playerID string)
	SendEventToClient(playerID string, eventToSend event.Event)
	publisher.EventListener
}

type clientEventSenderImp struct {
	clientConnections map[string]connector.ClientConnection
	clientUpdateRate  int
	timeFrame         uint32
	eventQueue        chan event.Event
	quit              chan struct{}
}

func (clientEventSender *clientEventSenderImp) Run() {
	clientUpdateTicker := time.NewTicker(time.Duration(1000/clientEventSender.clientUpdateRate) * time.Millisecond)
	for {
		select {
		case <-clientEventSender.quit:
			clientUpdateTicker.Stop()
			return
		case <-clientUpdateTicker.C:
			numberOfEvent := len(clientEventSender.eventQueue)
			if numberOfEvent > 0 {
				eventsToSend := make([]event.Event, len(clientEventSender.eventQueue))
				for i := 0; i < numberOfEvent; i++ {
					eventsToSend[i] = <-clientEventSender.eventQueue
					eventsToSend[i].TimeFrame = clientEventSender.timeFrame
				}
				for _, clientConnection := range clientEventSender.clientConnections {
					clientConnection.SendEventsToClient(eventsToSend)
				}
			}
			clientEventSender.timeFrame++
		}
	}
}

func (clientEventSender *clientEventSenderImp) AddClient(playerID string, connectionToClient connector.ClientConnection) {
	clientEventSender.clientConnections[playerID] = connectionToClient
}

func (clientEventSender *clientEventSenderImp) RemoveClient(playerID string) {
	delete(clientEventSender.clientConnections, playerID)
}

func (clientEventSender *clientEventSenderImp) SendEventToClient(playerID string, eventToSend event.Event) {
	eventToSend.TimeFrame = clientEventSender.timeFrame
	clientEventSender.clientConnections[playerID].SendEventsToClient([]event.Event{eventToSend})
}

func (clientEventSender *clientEventSenderImp) ReceiveEvent(event event.Event) {
	clientEventSender.eventQueue <- event
}
