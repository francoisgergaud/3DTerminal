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
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

var info = log.New(os.Stderr, "INFO ", 0)

//Impl is the default implementation for a server.
type Impl struct {
	worldMap          world.WorldMap
	players           map[string]animatedelement.AnimatedElement
	bots              map[string]bot.Bot
	quit              chan interface{}
	botsUpdateRate    int
	mathHelper        mathhelper.MathHelper
	clientEventSender clientEventSender
	runner            runner.Runner
	identifierFactory func() uuid.UUID
	worldMapFactory   func() world.WorldMap
	botFactory        func(id string, worldMap world.WorldMap, mathHelper mathhelper.MathHelper, quit <-chan interface{}) bot.Bot
	playerFactory     func(worldMap world.WorldMap, mathHelper helper.MathHelper, quit <-chan interface{}) animatedelement.AnimatedElement
}

//NewServer is a server factory
func NewServer(worldUpdateRate int, quit chan interface{}) (*Impl, error) {
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
		shutdownCompleted: make(chan interface{}),
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
	info.Print("starting server...")
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
	info.Printf("register new player with id %v", playerID)
	server.clientEventSender.addClient(playerID, clientConnection)
	player := server.playerFactory(server.worldMap, server.mathHelper, server.quit)
	server.players[playerID] = player
	newPlayerEvent := event.Event{
		PlayerID: playerID,
		State:    player.State(),
		Action:   "join",
	}
	worldMap := server.worldMap.Clone()
	server.clientEventSender.ReceiveEvent(newPlayerEvent)
	otherPlayers := make(map[string]state.AnimatedElementState)
	for id, player := range server.players {
		if id != playerID {
			otherPlayers[id] = player.State().Clone()
		}
	}
	for id, bot := range server.bots {
		otherPlayers[id] = bot.State().Clone()
	}
	extraData := make(map[string]interface{})
	extraData["worldMap"] = worldMap
	extraData["otherPlayers"] = otherPlayers
	newPlayerInitializationEvent := event.Event{
		Action:    "init",
		PlayerID:  playerID,
		State:     player.State(),
		ExtraData: extraData,
	}
	server.clientEventSender.sendEventToClient(playerID, newPlayerInitializationEvent)
	return playerID
}

//UnregisterClient removes a player
func (server *Impl) UnregisterClient(playerID string) {
	info.Printf("unregister new player with id %v", playerID)
	delete(server.players, playerID)
	server.clientEventSender.removeClient(playerID)
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
func (server *Impl) Run() error {
	environmentTicker := time.NewTicker(time.Duration(1000/server.botsUpdateRate) * time.Millisecond)
	for {
		select {
		case <-server.quit:
			environmentTicker.Stop()
			return nil
		case <-environmentTicker.C:
			for _, bot := range server.bots {
				bot.Move()
			}
			for _, player := range server.players {
				player.Move()
			}
		}
	}
}

//Shutdown waits for the gracefull shutdown to complete
func (server *Impl) Shutdown() {
	server.clientEventSender.shutdown()
}

type clientEventSender interface {
	runner.Runnable
	addClient(playerID string, connectionToClient connector.ClientConnection)
	removeClient(playerID string)
	sendEventToClient(playerID string, eventToSend event.Event)
	publisher.EventListener
	close()
	shutdown()
}

type clientEventSenderImp struct {
	clientConnections map[string]connector.ClientConnection
	clientUpdateRate  int
	timeFrame         uint32
	eventQueue        chan event.Event
	quit              <-chan interface{}
	shutdownCompleted chan interface{}
}

func (clientEventSender *clientEventSenderImp) Run() error {
	clientUpdateTicker := time.NewTicker(time.Duration(1000/clientEventSender.clientUpdateRate) * time.Millisecond)
	for {
		select {
		case <-clientEventSender.quit:
			clientUpdateTicker.Stop()
			clientEventSender.close()
			return nil
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

func (clientEventSender *clientEventSenderImp) addClient(playerID string, connectionToClient connector.ClientConnection) {
	clientEventSender.clientConnections[playerID] = connectionToClient
}

func (clientEventSender *clientEventSenderImp) removeClient(playerID string) {
	clientEventSender.clientConnections[playerID].Close()
	delete(clientEventSender.clientConnections, playerID)
}

func (clientEventSender *clientEventSenderImp) sendEventToClient(playerID string, eventToSend event.Event) {
	eventToSend.TimeFrame = clientEventSender.timeFrame
	clientEventSender.clientConnections[playerID].SendEventsToClient([]event.Event{eventToSend})
}

func (clientEventSender *clientEventSenderImp) ReceiveEvent(event event.Event) {
	clientEventSender.eventQueue <- event
}

func (clientEventSender *clientEventSenderImp) close() {
	for _, clientConnection := range clientEventSender.clientConnections {
		clientConnection.Close()
	}
	close(clientEventSender.shutdownCompleted)
}

func (clientEventSender *clientEventSenderImp) shutdown() {
	<-clientEventSender.shutdownCompleted
}
