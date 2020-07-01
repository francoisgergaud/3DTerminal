package impl

import (
	"fmt"
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/client/configuration"
	"francoisgergaud/3dGame/client/connector"
	"francoisgergaud/3dGame/client/consolemanager"
	"francoisgergaud/3dGame/client/render"
	renderImpl "francoisgergaud/3dGame/client/render/impl"
	renderMathHelperImpl "francoisgergaud/3dGame/client/render/mathhelper/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	animatedElementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement/projectile"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/common/math"
	"francoisgergaud/3dGame/common/math/helper"
	mathHelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/math/raycaster"
	"francoisgergaud/3dGame/common/runner"
	"log"
	originalMath "math"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/google/uuid"
)

var info = log.New(os.Stderr, "client ", 0)

//Impl implements the Engine interface.
type Impl struct {
	runner.Runner
	screen                                tcell.Screen
	playerID                              string
	worldMap                              world.WorldMap
	otherPlayers                          map[string]animatedelement.AnimatedElement
	projectiles                           map[string]projectile.Projectile
	player                                animatedelement.AnimatedElement
	otherPlayerLastUpdates                map[string]uint32
	renderer                              render.Renderer
	playerListener                        *playerListenerImpl
	worldElementUpdater                   *worldElementUpdaterImpl
	preInitializationEventFromServerQueue chan event.Event
	quit                                  <-chan interface{}
	frameRate                             int
	mathHelper                            mathHelper.MathHelper
	consoleEventManager                   consolemanager.ConsoleEventManager
	shutdown                              chan interface{}
	initialized, waitSpawnFromServer      bool
	connectionToServer                    connector.ServerConnector
	animatedElementFactory                func(id string, animatedElementState *state.AnimatedElementState, world world.WorldMap, mathHelper mathHelper.MathHelper) animatedelement.AnimatedElement
	projectileFactory                     func(id string, position *math.Point2D, angle float64, world world.WorldMap, otherPlayers map[string]animatedelement.AnimatedElement, mathHelper helper.MathHelper) projectile.Projectile
	identifierFactory                     func() uuid.UUID
}

//NewEngine provides a new engine.
func NewEngine(screen tcell.Screen, consoleEventManager consolemanager.ConsoleEventManager, engineConfig *configuration.Configuration, quit <-chan interface{}) (*Impl, error) {
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
		preInitializationEventFromServerQueue: make(chan event.Event, 100),
		quit:                                  quit,
		frameRate:                             engineConfig.FrameRate,
		mathHelper:                            mathHelper,
		consoleEventManager:                   consoleEventManager,
		shutdown:                              make(chan interface{}),
		initialized:                           false,
		waitSpawnFromServer:                   false,
		animatedElementFactory:                animatedElementImpl.NewAnimatedElementWithState,
		projectileFactory:                     projectile.NewProjectile,
		identifierFactory:                     uuid.New,
		playerListener: &playerListenerImpl{
			playerEventQueue: make(chan event.Event),
			quit:             quit,
		},
	}
	worldElementUpdater := &worldElementUpdaterImpl{
		updateRate: engineConfig.WorlUpdateRate,
		quit:       quit,
		engine:     &engine,
	}
	engine.worldElementUpdater = worldElementUpdater
	engine.Runner = &runner.AsyncRunner{}
	return &engine, nil
}

//Initialize set the engine player and environment
func (engine *Impl) initialize(playerID string, playerState *state.AnimatedElementState, worldMap world.WorldMap, otherPlayerStates map[string]*state.AnimatedElementState, projectileStates map[string]*state.AnimatedElementState, serverTimeFrame uint32) {
	engine.playerID = playerID
	engine.player = engine.animatedElementFactory(playerID, playerState, worldMap, engine.mathHelper)
	engine.consoleEventManager.SetPlayer(engine)
	engine.Runner.Start(engine.consoleEventManager)
	engine.worldMap = worldMap
	engine.otherPlayers = make(map[string]animatedelement.AnimatedElement)
	engine.projectiles = make(map[string]projectile.Projectile)
	engine.otherPlayerLastUpdates = make(map[string]uint32)
	for id, otherPlayerState := range otherPlayerStates {
		engine.otherPlayers[id] = engine.animatedElementFactory(id, otherPlayerState, worldMap, engine.mathHelper)
		engine.otherPlayerLastUpdates[id] = serverTimeFrame
	}
	for id, projectileState := range projectileStates {
		engine.projectiles[id] = engine.projectileFactory(id, projectileState.Position, projectileState.Angle, worldMap, engine.otherPlayers, engine.mathHelper)
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
			if event.Action == "join" || event.Action == "spawn" {
				engine.otherPlayers[event.PlayerID] = animatedElementImpl.NewAnimatedElementWithState(event.PlayerID, event.State, engine.worldMap, engine.mathHelper)
				engine.otherPlayerLastUpdates[event.PlayerID] = event.TimeFrame
			} else if event.Action == "move" {
				if event.TimeFrame > engine.otherPlayerLastUpdates[event.PlayerID] {
					engine.otherPlayers[event.PlayerID].SetState(event.State)
					engine.otherPlayerLastUpdates[event.PlayerID] = event.TimeFrame
				}
			} else if event.Action == "quit" || event.Action == "kill" {
				//other-player removed
				delete(engine.otherPlayerLastUpdates, event.PlayerID)
				delete(engine.otherPlayers, event.PlayerID)
			} else if event.Action == "fire" {
				//On fire-event, the playerID field is the player firing
				projectileID := event.ExtraData["projectileID"].(string)
				engine.projectiles[projectileID] = engine.projectileFactory(projectileID, event.State.Position, event.State.Angle, engine.worldMap, engine.otherPlayers, engine.mathHelper)
			} else if event.Action == "projectileImpact" {
				//On projectileImpact-event, the playerID field is the projectile's identifier
				delete(engine.projectiles, event.PlayerID)
			}
		} else {
			if event.Action == "kill" {
				engine.waitSpawnFromServer = true
				fmt.Printf("killed. Wait for respawn...")
			} else if event.Action == "spawn" {
				engine.player.SetState(event.State)
				engine.waitSpawnFromServer = false
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
		otherPlayerStates, _ := initializationEvent.ExtraData["otherPlayers"].(map[string]*state.AnimatedElementState)
		projectileStates, _ := initializationEvent.ExtraData["projectiles"].(map[string]*state.AnimatedElementState)
		playerState := initializationEvent.State
		engine.initialize(initializationEvent.PlayerID, playerState, worldMap, otherPlayerStates, projectileStates, initializationEvent.TimeFrame)
		engine.Runner.Start(engine)
		engine.Runner.Start(engine.worldElementUpdater)
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

//Player returns the player
func (engine *Impl) Player() animatedelement.AnimatedElement {
	return engine.player
}

//OtherPlayers returns the engine's other players.
func (engine *Impl) OtherPlayers() map[string]animatedelement.AnimatedElement {
	return engine.otherPlayers
}

//Projectiles returns the engine's projectiles.
func (engine *Impl) Projectiles() map[string]projectile.Projectile {
	return engine.projectiles
}

//Shutdown waits for the gracefull shutdown to complete
func (engine *Impl) Shutdown() {
	<-engine.shutdown
}

//ConnectToServer set the connection to server once initialized
func (engine *Impl) ConnectToServer(connectionToServer connector.ServerConnector) {
	engine.connectionToServer = connectionToServer
	engine.playerListener.connectionToServer = connectionToServer
}

//Run initializes the required element and start the engine to render world's elements in pseudo-3D
func (engine *Impl) Run() error {
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
			return nil
		case <-frameUpdateTicker.C:
			engine.renderer.Render(engine.playerID, engine.worldMap, engine.player, engine.otherPlayers, engine.projectiles, engine.screen)
		}
	}
}

// Action the player according to the input key
func (engine *Impl) Action(eventKey *tcell.EventKey) {
	playerState := engine.player.State()
	var eventToSend event.Event
	switch eventKey.Key() {
	case tcell.KeyUp:
		if playerState.MoveDirection == state.Backward {
			playerState.MoveDirection = state.None
		} else {
			playerState.MoveDirection = state.Forward
		}
		eventToSend = event.Event{Action: "move", State: playerState, TimeFrame: 0}
	case tcell.KeyDown:
		if playerState.MoveDirection == state.Forward {
			playerState.MoveDirection = state.None
		} else {
			playerState.MoveDirection = state.Backward
		}
		eventToSend = event.Event{Action: "move", State: playerState, TimeFrame: 0}
	case tcell.KeyLeft:
		if playerState.RotateDirection == state.Right {
			playerState.RotateDirection = state.None
		} else {
			playerState.RotateDirection = state.Left
		}
		eventToSend = event.Event{Action: "move", State: playerState, TimeFrame: 0}
	case tcell.KeyRight:
		if playerState.RotateDirection == state.Left {
			playerState.RotateDirection = state.None
		} else {
			playerState.RotateDirection = state.Right
		}
		eventToSend = event.Event{Action: "move", State: playerState, TimeFrame: 0}
	case tcell.KeyEnter:
		projectileID := engine.playerID + "." + engine.identifierFactory().String()
		projectileStartFactor := 1.5
		projectilePosition := &math.Point2D{
			X: playerState.Position.X + (playerState.Size*projectileStartFactor)*originalMath.Cos(playerState.Angle*originalMath.Pi),
			Y: playerState.Position.Y + (playerState.Size*projectileStartFactor)*originalMath.Sin(playerState.Angle*originalMath.Pi),
		}
		engine.projectiles[projectileID] = engine.projectileFactory(projectileID, projectilePosition, playerState.Angle, engine.worldMap, engine.otherPlayers, engine.mathHelper)
		eventToSend = event.Event{
			State:  engine.projectiles[projectileID].State(),
			Action: "fire",
			ExtraData: map[string]interface{}{
				"projectileID": projectileID,
			},
		}
	}
	if !engine.waitSpawnFromServer {
		engine.playerListener.playerEventQueue <- eventToSend
	}
}

//playerListenerImpl results from an internal decompostion of the client
type playerListenerImpl struct {
	playerEventQueue   chan event.Event
	quit               <-chan interface{}
	connectionToServer connector.ServerConnector
}

func (playerListener *playerListenerImpl) Run() error {
	for {
		select {
		case eventFromPlayer := <-playerListener.playerEventQueue:
			playerListener.connectionToServer.NotifyServer([]event.Event{eventFromPlayer})
		case <-playerListener.quit:
			return nil
		}
	}
}

//worldElementUpdaterImpl results from an internal decompostion of the client to manage the client-side worl-update
type worldElementUpdaterImpl struct {
	updateRate int
	engine     client.Engine
	quit       <-chan interface{}
}

//loop of an internal clock events to update the player an world-elements based of their state (direction, position, velocity etc...)
func (worldElementUpdater *worldElementUpdaterImpl) Run() error {
	worldUpdateTicker := time.NewTicker(time.Duration(1000/worldElementUpdater.updateRate) * time.Millisecond)
	for {
		select {
		case <-worldElementUpdater.quit:
			worldUpdateTicker.Stop()
			return nil
		case <-worldUpdateTicker.C:
			worldElementUpdater.engine.Player().Move()
			for _, worldelement := range worldElementUpdater.engine.OtherPlayers() {
				worldelement.Move()
			}
			for _, projectile := range worldElementUpdater.engine.Projectiles() {
				projectile.Move()
			}
		}
	}
}
