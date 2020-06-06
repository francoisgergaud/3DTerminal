package impl

import (
	"fmt"
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/client/connector"
	"francoisgergaud/3dGame/client/consolemanager"
	"francoisgergaud/3dGame/client/player"
	playerImpl "francoisgergaud/3dGame/client/player/impl"
	"francoisgergaud/3dGame/client/render"
	renderImpl "francoisgergaud/3dGame/client/render/impl"
	renderMathHelperImpl "francoisgergaud/3dGame/client/render/mathhelper/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	animatedElementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/math/helper"
	mathHelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/math/raycaster"
	"francoisgergaud/3dGame/common/runner"
	"time"

	"github.com/gdamore/tcell"
)

//Impl implements the Engine interface.
type Impl struct {
	runner.Runner
	screen                                tcell.Screen
	playerID                              string
	worldMap                              world.WorldMap
	otherPlayers                          map[string]animatedelement.AnimatedElement
	player                                player.Player
	otherPlayerLastUpdates                map[string]uint32
	renderer                              render.Renderer
	playerListener                        *playerListenerImpl
	worldElementUpdater                   *worldElementUpdaterImpl
	preInitializationEventFromServerQueue chan event.Event
	quit                                  chan struct{}
	frameRate                             int
	mathHelper                            mathHelper.MathHelper
	consoleEventManager                   consolemanager.ConsoleEventManager
	shutdown                              chan interface{}
	initialized                           bool
	connectionToServer                    connector.ServerConnector
	animatedElementFactory                func(animatedElementState *state.AnimatedElementState, world world.WorldMap, mathHelper mathHelper.MathHelper) animatedelement.AnimatedElement
	playerFactory                         func(playerState *state.AnimatedElementState, world world.WorldMap, mathHelper helper.MathHelper) player.Player
}

//NewEngine provides a new engine.
func NewEngine(screen tcell.Screen, consoleEventManager consolemanager.ConsoleEventManager, engineConfig *client.Configuration) (*Impl, error) {
	raySampler, err := renderImpl.CreateRaySamplerForAnsiColorTerminal(
		engineConfig.GradientRSFirst,
		engineConfig.GradientRSMultiplicator,
		engineConfig.GradientRSLimit,
		engineConfig.GradientRSWallStartColor,
		engineConfig.GradientRSWallEndColor,
		engineConfig.ScreenHeight,
		engineConfig.GradientRSBackgroundRange,
		engineConfig.GradientRSBackgroundColors)
	if err != nil {
		return nil, fmt.Errorf("error while instantiating the ray-sampler: %w", err)
	}
	mathHelper, err := mathHelper.NewMathHelper(new(raycaster.RayCasterImpl))
	//mathHelper, err := mathHelper.NewMathHelper(new(raycaster.RayCasterImpl))
	if err != nil {
		return nil, fmt.Errorf("error while instantiating the math-helper: %w", err)
	}
	renderMathHelper := renderMathHelperImpl.NewRendererMathHelper(mathHelper)
	renderer := renderImpl.CreateRenderer(engineConfig.ScreenWidth, engineConfig.ScreenHeight, raySampler, mathHelper, renderMathHelper, engineConfig.PlayerFieldOfViewAngle, engineConfig.Visibility)
	engine := Impl{
		screen:                                screen,
		renderer:                              renderer,
		preInitializationEventFromServerQueue: make(chan event.Event, 100),
		quit:                                  engineConfig.QuitChannel,
		frameRate:                             engineConfig.FrameRate,
		mathHelper:                            mathHelper,
		consoleEventManager:                   consoleEventManager,
		shutdown:                              make(chan interface{}),
		initialized:                           false,
		animatedElementFactory:                animatedElementImpl.NewAnimatedElementWithState,
		playerFactory:                         playerImpl.NewPlayer,
		playerListener: &playerListenerImpl{
			playerEventQueue: make(chan event.Event),
			quit:             engineConfig.QuitChannel,
		},
	}
	worldElementUpdater := &worldElementUpdaterImpl{
		updateRate: engineConfig.WorlUpdateRate,
		quit:       engineConfig.QuitChannel,
		engine:     &engine,
	}
	engine.worldElementUpdater = worldElementUpdater
	engine.Runner = &runner.AsyncRunner{}
	consoleEventManager.Listen()
	return &engine, nil
}

//Initialize set the engine player and environment
func (engine *Impl) initialize(playerID string, playerState *state.AnimatedElementState, worldMap world.WorldMap, otherPlayersState map[string]*state.AnimatedElementState, serverTimeFrame uint32) {
	engine.playerID = playerID
	engine.player = engine.playerFactory(playerState, worldMap, engine.mathHelper)
	engine.consoleEventManager.SetPlayer(engine.player)
	engine.worldMap = worldMap
	engine.otherPlayers = make(map[string]animatedelement.AnimatedElement)
	engine.otherPlayerLastUpdates = make(map[string]uint32)
	for id, otherPlayerState := range otherPlayersState {
		engine.otherPlayers[id] = engine.animatedElementFactory(otherPlayerState, worldMap, engine.mathHelper)
		engine.otherPlayerLastUpdates[id] = serverTimeFrame
	}
}

//ReceiveEventsFromServer manages the event received from the server
func (engine *Impl) ReceiveEventsFromServer(events []event.Event) {
	if engine.initialized {
		engine.processPostInitializationEvents(events)
	} else {
		engine.processPreInitializationEvents(events)
	}
}

func (engine *Impl) processPostInitializationEvents(events []event.Event) {
	for _, event := range events {
		if event.PlayerID != engine.playerID {
			if event.Action == "join" {
				engine.otherPlayers[event.PlayerID] = animatedElementImpl.NewAnimatedElementWithState(event.State, engine.worldMap, engine.mathHelper)
				engine.otherPlayerLastUpdates[event.PlayerID] = event.TimeFrame
			} else if event.Action == "move" {
				if event.TimeFrame > engine.otherPlayerLastUpdates[event.PlayerID] {
					engine.otherPlayers[event.PlayerID].SetState(event.State)
					engine.otherPlayerLastUpdates[event.PlayerID] = event.TimeFrame
				}
			} else if event.Action == "quit" {
				//other-player removed
				delete(engine.otherPlayerLastUpdates, event.PlayerID)
				delete(engine.otherPlayers, event.PlayerID)
			}
		}
	}
}

func (engine *Impl) processPreInitializationEvents(events []event.Event) {
	var initializationEvent *event.Event
	for _, eventFromServer := range events {
		if eventFromServer.Action == "init" {
			initializationEvent = &eventFromServer
		} else {
			//TODO: manage properly the pre-initialization-events queue (don't block if the queue is full)
			engine.preInitializationEventFromServerQueue <- eventFromServer
		}
	}
	if initializationEvent != nil {
		//initialize and start the client
		engine.playerID = initializationEvent.PlayerID
		worldMap, _ := initializationEvent.ExtraData["worldMap"].(world.WorldMap)
		//TODO: manage cast error (use a warn in log file)
		otherPlayerStates, _ := initializationEvent.ExtraData["otherPlayers"].(map[string]state.AnimatedElementState)
		otherPlayerStatesClone := make(map[string]*state.AnimatedElementState)
		for otherPlayerID, otherPlayerState := range otherPlayerStates {
			otherPlayerStateClone := otherPlayerState.Clone()
			otherPlayerStatesClone[otherPlayerID] = &otherPlayerStateClone
		}
		playerStateClone := initializationEvent.State.Clone()
		engine.initialize(initializationEvent.PlayerID, &playerStateClone, worldMap.Clone(), otherPlayerStatesClone, initializationEvent.TimeFrame)
		engine.Runner.Start(engine)
		engine.Runner.Start(engine.worldElementUpdater)
		//propagate events from player
		engine.Player().RegisterListener(engine.playerListener)
		//process all previous events
		numberOfPreInitializationEvents := len(engine.preInitializationEventFromServerQueue)
		if numberOfPreInitializationEvents > 0 {
			preInitializationEvents := make([]event.Event, numberOfPreInitializationEvents)
			for i := 0; i < numberOfPreInitializationEvents; i++ {
				preInitializationEvents[i] = <-engine.preInitializationEventFromServerQueue
			}
			engine.processPostInitializationEvents(preInitializationEvents)
		}
		engine.Runner.Start(engine.playerListener)
		//change the state
		engine.initialized = true
	}
}

//Player returns the engine's player.
func (engine *Impl) Player() player.Player {
	return engine.player
}

//OtherPlayers returns the engine's other players.
func (engine *Impl) OtherPlayers() map[string]animatedelement.AnimatedElement {
	return engine.otherPlayers
}

//Shutdown return the channel to be closed when shutdown is gracefully operated
func (engine *Impl) Shutdown() <-chan interface{} {
	return engine.shutdown
}

//ConnectToServer set the connection to server once initialized
func (engine *Impl) ConnectToServer(connectionToServer connector.ServerConnector) {
	engine.connectionToServer = connectionToServer
	engine.playerListener.connectionToServer = connectionToServer
}

//Run initializes the required element and start the engine to render world's elements in pseudo-3D
func (engine *Impl) Run() {
	engine.screen.Clear()
	//TODO: manage division by 0 in a cleaner way
	frameUpdateTicker := time.NewTicker(time.Duration(1000/engine.frameRate) * time.Millisecond)
	for {
		select {
		case <-engine.quit:
			frameUpdateTicker.Stop()
			engine.screen.SetStyle(tcell.StyleDefault)
			engine.screen.Clear()
			engine.screen.Fini()
			if engine.connectionToServer != nil {
				engine.connectionToServer.Disconnect()
			}
			close(engine.shutdown)
			return
		case <-frameUpdateTicker.C:
			engine.renderer.Render(engine.playerID, engine.worldMap, engine.player, engine.otherPlayers, engine.screen)
		}
	}
}

//playerListenerImpl results from an internal decompostion of the client
type playerListenerImpl struct {
	playerEventQueue   chan event.Event
	quit               chan struct{}
	connectionToServer connector.ServerConnector
}

func (playerListener *playerListenerImpl) Run() {
	for {
		select {
		case /*eventsFromPlayer[0] =*/ eventFromPlayer := <-playerListener.playerEventQueue:
			playerListener.connectionToServer.NotifyServer([]event.Event{eventFromPlayer})
		case <-playerListener.quit:
			return
		}
	}
}

func (playerListener *playerListenerImpl) ReceiveEvent(event event.Event) {
	playerListener.playerEventQueue <- event
}

//worldElementUpdaterImpl results from an internal decompostion of the client to manage the client-side worl-update
type worldElementUpdaterImpl struct {
	updateRate int
	engine     client.Engine
	quit       chan struct{}
}

//loop of an internal clock events to update the player an world-elements based of their state (direction, position, velocity etc...)
func (worldElementUpdater *worldElementUpdaterImpl) Run() {
	worldUpdateTicker := time.NewTicker(time.Duration(1000/worldElementUpdater.updateRate) * time.Millisecond)
	for {
		select {
		case <-worldElementUpdater.quit:
			worldUpdateTicker.Stop()
			return
		case <-worldUpdateTicker.C:
			worldElementUpdater.engine.Player().Move()
			for _, worldelement := range worldElementUpdater.engine.OtherPlayers() {
				worldelement.Move()
			}
		}
	}
}
