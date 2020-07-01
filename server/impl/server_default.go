package impl

import (
	"fmt"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/projectile"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"
	mathhelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/math/raycaster"
	"francoisgergaud/3dGame/common/runner"
	"francoisgergaud/3dGame/server/bot"
	"francoisgergaud/3dGame/server/connector"
	botgenerator "francoisgergaud/3dGame/server/impl/generator/bot"
	"francoisgergaud/3dGame/server/impl/generator/player"
	"francoisgergaud/3dGame/server/impl/generator/worldmap"
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
	projectiles       map[string]projectile.Projectile
	botIDs            []string
	quit              chan interface{}
	botsUpdateRate    int
	mathHelper        mathhelper.MathHelper
	clientEventSender clientEventSender
	runner            runner.Runner
	identifierFactory func() uuid.UUID
	worldMapFactory   func() world.WorldMap
	botFactory        func(id string, worldMap world.WorldMap, mathHelper mathhelper.MathHelper, quit <-chan interface{}) bot.Bot
	playerFactory     func(wid string, orldMap world.WorldMap, mathHelper helper.MathHelper, quit <-chan interface{}) animatedelement.AnimatedElement
	projectileFactory func(id string, position *math.Point2D, angle float64, world world.WorldMap, otherPlayers map[string]animatedelement.AnimatedElement, mathHelper helper.MathHelper) projectile.Projectile
	spawner           player.Spawner
}

//NewServer is a server factory
func NewServer(worldUpdateRate int, quit chan interface{}) (*Impl, error) {
	server := new(Impl)
	server.botIDs = make([]string, 0)
	mathHelper, err := mathhelper.NewMathHelper(new(raycaster.RayCasterImpl))
	server.mathHelper = mathHelper
	if err != nil {
		return nil, fmt.Errorf("error while instantiating the math-helper: %w", err)
	}
	server.players = make(map[string]animatedelement.AnimatedElement)
	server.projectiles = make(map[string]projectile.Projectile)
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
	server.projectileFactory = projectile.NewProjectile
	server.spawner = player.NewStaticSpawner(server.players)
	server.spawner.RegisterListener(server)
	return server, nil
}

//Start the server
func (server *Impl) Start() {
	info.Print("starting server...")
	//initialize the environment (world and bots)
	server.worldMap = server.worldMapFactory()
	botID := server.identifierFactory().String()
	bot := server.botFactory(botID, server.worldMap, server.mathHelper, server.quit)
	bot.RegisterListener(server)
	server.players[botID] = bot
	server.botIDs = append(server.botIDs, botID)
	//start the asynchronous listeners
	server.runner.Start(server.clientEventSender)
	server.runner.Start(server)
}

//RegisterPlayer register a player and provide the environment
func (server *Impl) RegisterPlayer(clientConnection connector.ClientConnection) string {
	playerID := server.identifierFactory().String()
	info.Printf("register new player with id %v", playerID)
	server.clientEventSender.addClient(playerID, clientConnection)
	player := server.playerFactory(playerID, server.worldMap, server.mathHelper, server.quit)
	server.players[playerID] = player
	newPlayerEvent := event.Event{
		PlayerID: playerID,
		State:    player.State(),
		Action:   "join",
	}
	server.clientEventSender.sendEventToAllClients(newPlayerEvent)
	otherPlayers := make(map[string]*state.AnimatedElementState)
	for id, player := range server.players {
		if id != playerID {
			otherPlayers[id] = player.State()
		}
	}
	extraData := make(map[string]interface{})
	extraData["worldMap"] = server.worldMap
	extraData["otherPlayers"] = otherPlayers
	projectilesStates := make(map[string]*state.AnimatedElementState)
	for id, projectile := range server.projectiles {
		projectilesStates[id] = projectile.State()
	}
	extraData["projectiles"] = projectilesStates
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
	server.clientEventSender.sendEventToAllClients(event)
}

//ReceiveEventFromClient manage an event received from a client
// as it is supposed to override the previous ones
func (server *Impl) ReceiveEventFromClient(event event.Event) {
	if event.Action == "fire" {
		projectileID := event.ExtraData["projectileID"].(string)
		projectile := server.projectileFactory(projectileID, event.State.Position, event.State.Angle, server.worldMap, server.players, server.mathHelper)
		server.projectiles[projectileID] = projectile
		server.projectiles[projectileID].RegisterListener(server)
		server.clientEventSender.sendEventToAllClients(event)
	} else if event.Action == "move" {
		server.players[event.PlayerID].SetState(event.State)
		server.clientEventSender.sendEventToAllClients(event)
	}
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
			for _, player := range server.players {
				player.Move()
			}
			for _, projectile := range server.projectiles {
				projectile.Move()
			}
		}
	}
}

//ReceiveEvent receives event the server subscribed for
func (server *Impl) ReceiveEvent(eventReceived event.Event) {
	if eventReceived.Action == "projectileWallImpact" {
		delete(server.projectiles, eventReceived.PlayerID)
		eventReceived.Action = "projectileImpact"
		server.clientEventSender.sendEventToAllClients(eventReceived)
	} else if eventReceived.Action == "projectilePlayerImpact" {
		delete(server.projectiles, eventReceived.PlayerID)
		playerKilledID := eventReceived.ExtraData["playerID"].(string)
		eventReceived.Action = "projectileImpact"
		server.clientEventSender.sendEventToAllClients(eventReceived)
		killEvent := event.Event{
			Action:   "kill",
			PlayerID: playerKilledID,
		}
		server.clientEventSender.sendEventToAllClients(killEvent)
		//if the player killed is a bot, the server has to make it move forward
		moveDirection := state.None
		for _, botID := range server.botIDs {
			if botID == playerKilledID {
				moveDirection = state.Forward
				break
			}
		}
		server.spawner.Spawn(playerKilledID, moveDirection)
	} else if eventReceived.Action == "move" {
		server.players[eventReceived.PlayerID].SetState(eventReceived.State)
		server.clientEventSender.sendEventToAllClients(eventReceived)
	} else if eventReceived.Action == "spawn" {
		server.clientEventSender.sendEventToAllClients(eventReceived)
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
	sendEventToAllClients(eventToSend event.Event)
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

func (clientEventSender *clientEventSenderImp) sendEventToAllClients(event event.Event) {
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
