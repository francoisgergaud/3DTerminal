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
	mathHelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/math/raycaster"
	"time"

	"github.com/gdamore/tcell"
)

//Impl implements the Engine interface.
type Impl struct {
	screen                                tcell.Screen
	player                                player.Player
	playerID                              string
	worldMap                              world.WorldMap
	otherPlayers                          map[string]animatedelement.AnimatedElement
	otherPlayerLastUpdates                map[string]uint32
	renderer                              render.Renderer
	playerEventQueue                      chan event.Event
	preInitializationEventFromServerQueue chan event.Event
	quit                                  chan struct{}
	frameRate                             int
	updateRate                            int
	mathHelper                            mathHelper.MathHelper
	consoleEventManager                   consolemanager.ConsoleEventManager
	shutdown                              chan interface{}
	intialized                            bool
	connectionToServer                    connector.ServerConnector
}

//NewEngine provides a new engine.
func NewEngine(screen tcell.Screen, consoleEventManager consolemanager.ConsoleEventManager, engineConfig *client.Configuration) (client.Engine, error) {
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
	if err != nil {
		return nil, fmt.Errorf("error while instantiating the math-helper: %w", err)
	}
	renderMathHelper := renderMathHelperImpl.NewRendererMathHelper(mathHelper)
	renderer := renderImpl.CreateRenderer(engineConfig.ScreenWidth, engineConfig.ScreenHeight, raySampler, mathHelper, renderMathHelper, engineConfig.PlayerFieldOfViewAngle, engineConfig.Visibility)
	engine := Impl{
		screen:                                screen,
		renderer:                              renderer,
		playerEventQueue:                      make(chan event.Event),
		preInitializationEventFromServerQueue: make(chan event.Event, 100),
		quit:                                  engineConfig.QuitChannel,
		frameRate:                             engineConfig.FrameRate,
		updateRate:                            engineConfig.WorlUpdateRate,
		mathHelper:                            mathHelper,
		consoleEventManager:                   consoleEventManager,
		shutdown:                              make(chan interface{}),
		intialized:                            false,
	}
	consoleEventManager.Listen()
	return &engine, nil
}

//Initialize set the engine player and environment
func (engine *Impl) Initialize(playerID string, playerState state.AnimatedElementState, worldMap world.WorldMap, otherPlayersState map[string]state.AnimatedElementState, serverTimeFrame uint32) {
	engine.playerID = playerID
	engine.player = playerImpl.NewPlayer(playerState, worldMap, engine.mathHelper, engine.quit)
	engine.consoleEventManager.SetPlayer(engine.GetPlayer())
	engine.worldMap = worldMap
	engine.otherPlayers = make(map[string]animatedelement.AnimatedElement)
	engine.otherPlayerLastUpdates = make(map[string]uint32)
	for id, otherPlayerState := range otherPlayersState {
		engine.otherPlayers[id] = animatedElementImpl.NewAnimatedElementWithState(otherPlayerState, worldMap, engine.mathHelper, engine.quit)
		engine.otherPlayerLastUpdates[id] = serverTimeFrame
	}
}

//StartEngine initializes the required element and start the engine to render world's elements in pseudo-3D
func (engine *Impl) StartEngine() {
	go func() {
		engine.screen.Clear()
		engine.player.Start()
		frameUpdateTicker := time.NewTicker(time.Duration(1000/engine.frameRate) * time.Millisecond)
		worldUpdateTicker := time.NewTicker(time.Duration(1000/engine.updateRate) * time.Millisecond)
		for _, worldelement := range engine.otherPlayers {
			worldelement.Start()
		}
		for {
			select {
			case <-engine.quit:
				frameUpdateTicker.Stop()
				worldUpdateTicker.Stop()
				engine.screen.SetStyle(tcell.StyleDefault)
				engine.screen.Clear()
				engine.screen.Fini()
				close(engine.shutdown)
				return
			case <-frameUpdateTicker.C:
				engine.renderer.Render(engine.playerID, engine.worldMap, engine.player, engine.otherPlayers, engine.screen)
			case updateWorldTickerTime := <-worldUpdateTicker.C:
				engine.player.GetUpdateChannel() <- updateWorldTickerTime
				for _, worldelement := range engine.otherPlayers {
					worldelement.GetUpdateChannel() <- updateWorldTickerTime
				}
			}
		}
	}()
}

//ReceiveEventsFromServer manages the event received from the server
func (engine *Impl) ReceiveEventsFromServer(events []event.Event) {
	if engine.intialized {
		engine.processOtherPlayerEvents(events)
	} else {
		for _, eventFromServer := range events {
			if eventFromServer.Action == "init" {
				//initialize and start the client
				engine.playerID = eventFromServer.PlayerID
				worldMap, _ := eventFromServer.ExtraData["worldMap"].(world.WorldMap)
				otherPlayerStates, _ := eventFromServer.ExtraData["otherPlayers"].(map[string]state.AnimatedElementState)
				otherPlayerStatesClone := make(map[string]state.AnimatedElementState)
				for otherPlayerID, otherPlayerState := range otherPlayerStates {
					otherPlayerStatesClone[otherPlayerID] = otherPlayerState.Clone()
				}
				engine.Initialize(eventFromServer.PlayerID, eventFromServer.State.Clone(), worldMap.Clone(), otherPlayerStatesClone, eventFromServer.TimeFrame)
				engine.StartEngine()
				//propagate events from player
				engine.GetPlayer().RegisterListener(engine.playerEventQueue)
				//process all previous events
				numberOfPreInitializationEvents := len(engine.preInitializationEventFromServerQueue)
				if numberOfPreInitializationEvents > 0 {
					preInitializationEvents := make([]event.Event, numberOfPreInitializationEvents)
					for i := 0; i < numberOfPreInitializationEvents; i++ {
						preInitializationEvents[i] = <-engine.preInitializationEventFromServerQueue
					}
					engine.processOtherPlayerEvents(preInitializationEvents)
				}
				go engine.listenPlayer()
				//change the state
				engine.intialized = true
			} else {
				engine.preInitializationEventFromServerQueue <- eventFromServer
			}
		}
	}
}

func (engine *Impl) processOtherPlayerEvents(events []event.Event) {
	for _, event := range events {
		if event.PlayerID != engine.playerID {
			if event.Action == "join" {
				engine.otherPlayers[event.PlayerID] = animatedElementImpl.NewAnimatedElementWithState(*event.State, engine.worldMap, engine.mathHelper, engine.quit)
				engine.otherPlayerLastUpdates[event.PlayerID] = event.TimeFrame
				engine.otherPlayers[event.PlayerID].Start()
			} else if event.Action == "move" {
				if event.TimeFrame > engine.otherPlayerLastUpdates[event.PlayerID] {
					engine.otherPlayers[event.PlayerID].SetState(event.State)
					engine.otherPlayerLastUpdates[event.PlayerID] = event.TimeFrame
				}
			}
		}
	}
}

//GetPlayer returns the engine's player.
func (engine *Impl) GetPlayer() player.Player {
	return engine.player
}

//GetShutdown return the channel to be closed when shutdown is gracefully operated
func (engine *Impl) GetShutdown() <-chan interface{} {
	return engine.shutdown
}

func (engine *Impl) SetConnectionToServer(connectionToServer connector.ServerConnector) {
	engine.connectionToServer = connectionToServer
}

func (engine *Impl) listenPlayer() {
	//eventsFromPlayer := make([]event.Event, 1)
	for {
		select {
		case /*eventsFromPlayer[0] =*/ eventFromPlayer := <-engine.playerEventQueue:
			engine.connectionToServer.NotifyServer([]event.Event{eventFromPlayer})
		case <-engine.quit:
			engine.connectionToServer.Disconnect()
		}
	}
}
