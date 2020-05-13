package impl

import (
	"fmt"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	animatedelementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	internalmath "francoisgergaud/3dGame/common/math"
	mathhelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/math/raycaster"
	"francoisgergaud/3dGame/server"
	"francoisgergaud/3dGame/server/bot"
	botImpl "francoisgergaud/3dGame/server/bot/impl"
	"francoisgergaud/3dGame/server/connector"
	"time"

	"github.com/gdamore/tcell"
	"github.com/google/uuid"
)

//Impl is the default implementation for a server.
type Impl struct {
	clientConnections map[string]connector.ClientConnection
	worldMap          world.WorldMap
	players           map[string]animatedelement.AnimatedElement
	bots              map[string]bot.Bot
	timeFrame         uint32
	eventQueue        chan event.Event
	quit              chan struct{}
	clientUpdateRate  int
	botsUpdateRate    int
	mathHelper        mathhelper.MathHelper
}

//NewWorldMap provides a new world-map.
func NewWorldMap() [][]int {
	return [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
}

//NewBot creates a bot
func NewBot(id string, worldMap world.WorldMap, mathHelper mathhelper.MathHelper, quit chan struct{}) bot.Bot {
	position := &internalmath.Point2D{X: 9, Y: 12}
	initialAngle := 0.3
	velocity := 0.02
	size := 0.3
	stepAngle := 0.0
	moveDirection := state.Forward
	rotateDirection := state.None
	style := tcell.StyleDefault.Background(tcell.ColorDarkBlue)
	return botImpl.NewBotImpl(id, position, initialAngle, velocity, stepAngle, size, moveDirection, rotateDirection, style, worldMap, mathHelper, quit)
}

//NewServer is a server factory
func NewServer(worldUpdateRate int, quit chan struct{}) (server.Server, error) {
	server := new(Impl)
	worldMap := world.NewWorldMap(NewWorldMap())
	server.worldMap = worldMap
	server.bots = make(map[string]bot.Bot)
	mathHelper, err := mathhelper.NewMathHelper(new(raycaster.RayCasterImpl))
	server.mathHelper = mathHelper
	if err != nil {
		return nil, fmt.Errorf("error while instantiating the math-helper: %w", err)
	}
	server.timeFrame = 0
	server.players = make(map[string]animatedelement.AnimatedElement)
	server.eventQueue = make(chan event.Event, 100)
	server.clientConnections = make(map[string]connector.ClientConnection)
	botID := uuid.New().String()
	server.bots[botID] = NewBot(botID, worldMap, mathHelper, quit)
	server.bots[botID].RegisterListener(server.eventQueue)
	server.quit = quit
	server.clientUpdateRate = 10
	server.botsUpdateRate = worldUpdateRate
	go server.start()
	go server.startBots()
	return server, nil
}

//RegisterPlayer register a player and provide the environment
func (server *Impl) RegisterPlayer(clientConnection connector.ClientConnection) string {
	playerID := uuid.New().String()
	server.clientConnections[playerID] = clientConnection
	animatedElementState := state.AnimatedElementState{
		Position:  &internalmath.Point2D{X: 5, Y: 5},
		Angle:     0.0,
		Size:      0.5,
		Velocity:  0.1,
		StepAngle: 0.01,
		Style:     tcell.StyleDefault.Background(tcell.Color126),
	}
	server.players[playerID] = animatedelementImpl.NewAnimatedElementWithState(&animatedElementState, server.worldMap, server.mathHelper, server.quit)
	newPlayerEvent := event.Event{
		PlayerID:  playerID,
		State:     &animatedElementState,
		TimeFrame: server.timeFrame,
		Action:    "join",
	}
	worldMap := server.worldMap.Clone()
	server.eventQueue <- newPlayerEvent
	otherPlayers := make(map[string]state.AnimatedElementState)
	for id, player := range server.players {
		if id != playerID {
			otherPlayers[id] = player.GetState().Clone()
		}
	}
	for id, bot := range server.bots {
		otherPlayers[id] = bot.GetState().Clone()
	}
	timeFrame := server.timeFrame
	extraData := make(map[string]interface{})
	extraData["worldMap"] = worldMap
	extraData["otherPlayers"] = otherPlayers
	newPlayerInitializationEvent := event.Event{
		Action:    "init",
		PlayerID:  playerID,
		State:     &animatedElementState,
		TimeFrame: timeFrame,
		ExtraData: extraData,
	}
	clientConnection.SendEventsToClient([]event.Event{newPlayerInitializationEvent})
	server.players[playerID].Start()
	return playerID
}

//UnregisterClient removes a player
func (server *Impl) UnregisterClient(playerID string) {
	delete(server.players, playerID)
	event := event.Event{
		PlayerID:  playerID,
		TimeFrame: server.timeFrame,
		Action:    "quit",
	}
	server.eventQueue <- event
}

//ReceiveEventFromClient manage an event received from a client
// as it is supposed to override the previous ones
func (server *Impl) ReceiveEventFromClient(event event.Event) {
	server.players[event.PlayerID].SetState(event.State)
	server.eventQueue <- event
}

//The sync action is managed by sending the whole animated-element state when a change is done.
func (server *Impl) sendEventsToClients() {
	numberOfEvent := len(server.eventQueue)
	if numberOfEvent > 0 {
		eventsToSend := make([]event.Event, len(server.eventQueue))
		for i := 0; i < numberOfEvent; i++ {
			eventsToSend[i] = <-server.eventQueue
			eventsToSend[i].TimeFrame = server.timeFrame
		}
		for _, clientConnection := range server.clientConnections {
			clientConnection.SendEventsToClient(eventsToSend)
		}
	}
	server.timeFrame++
}

func (server *Impl) start() {
	clientUpdateTicker := time.NewTicker(time.Duration(1000/server.clientUpdateRate) * time.Millisecond)
	for {
		select {
		case <-server.quit:
			clientUpdateTicker.Stop()
			return
		case <-clientUpdateTicker.C:
			server.sendEventsToClients()
		}
	}
}

func (server *Impl) startBots() {
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
