package impl

import (
	"fmt"
	"francoisgergaud/3dGame/client"
	"francoisgergaud/3dGame/client/connector"
	"francoisgergaud/3dGame/client/player"
	playerImpl "francoisgergaud/3dGame/client/player/impl"
	"francoisgergaud/3dGame/client/render"
	renderImpl "francoisgergaud/3dGame/client/render/impl"
	renderMathHelperImpl "francoisgergaud/3dGame/client/render/mathhelper/impl"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	animatedElementImpl "francoisgergaud/3dGame/common/environment/animatedelement/impl"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	mathHelper "francoisgergaud/3dGame/common/math/helper"
	"francoisgergaud/3dGame/common/math/raycaster"
	"time"

	"github.com/gdamore/tcell"
)

//Impl implements the Engine interface.
type Impl struct {
	screen       tcell.Screen
	player       player.Player
	worldMap     world.WorldMap
	otherPlayers map[string]animatedelement.AnimatedElement
	renderer     render.Renderer
	quit         chan struct{}
	frameRate    int
	updateRate   int
	mathHelper   mathHelper.MathHelper
}

//NewEngine provides a new engine.
func NewEngine(screen tcell.Screen, engineConfig *client.Configuration, serverConnection connector.ServerConnection) (*Impl, error) {
	if engineConfig.WorldMap == nil {
		return nil, fmt.Errorf("world-map cannot be nil")
	}
	worldMap := engineConfig.WorldMap
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
	player := playerImpl.NewPlayer(engineConfig.PlayerID, engineConfig.PlayerConfiguration, worldMap, mathHelper, engineConfig.QuitChannel, serverConnection)
	renderMathHelper := renderMathHelperImpl.NewRendererMathHelper(mathHelper)
	renderer := renderImpl.CreateRenderer(engineConfig.ScreenWidth, engineConfig.ScreenHeight, raySampler, mathHelper, renderMathHelper, engineConfig.PlayerFieldOfViewAngle, engineConfig.Visibility)
	otherPlayers := make(map[string]animatedelement.AnimatedElement)
	for id, otherPlayerConfiguration := range engineConfig.OtherPlayerConfigurations {
		otherPlayers[id] = animatedElementImpl.NewAnimatedElementWithState(id, otherPlayerConfiguration, worldMap, mathHelper, engineConfig.QuitChannel)
	}
	engine := Impl{
		screen:       screen,
		player:       player,
		worldMap:     worldMap,
		otherPlayers: otherPlayers,
		renderer:     renderer,
		quit:         engineConfig.QuitChannel,
		frameRate:    engineConfig.FrameRate,
		updateRate:   engineConfig.WorlUpdateRate,
		mathHelper:   mathHelper,
	}
	return &engine, nil
}

//StartEngine initializes the required element and start the engine to render world's elements in pseudo-3D
func (engine *Impl) StartEngine() {
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
			return
		case <-frameUpdateTicker.C:
			engine.renderer.Render(engine.worldMap, engine.player, engine.otherPlayers, engine.screen)
		case updateWorldTickerTime := <-worldUpdateTicker.C:
			engine.player.GetUpdateChannel() <- updateWorldTickerTime
			for _, worldelement := range engine.otherPlayers {
				worldelement.GetUpdateChannel() <- updateWorldTickerTime
			}
		}
	}
}

//ReceiveEventsFromServer manages the event received from the server
func (engine *Impl) ReceiveEventsFromServer(events []event.Event) {
	for _, event := range events {
		if event.PlayerID != engine.player.GetID() {
			if event.Action == "join" {
				engine.otherPlayers[event.PlayerID] = animatedElementImpl.NewAnimatedElementWithState(event.PlayerID, *event.State, engine.worldMap, engine.mathHelper, engine.quit)
			} else if event.Action == "move" {
				engine.otherPlayers[event.PlayerID].SetState(event.State)
			}
		}
	}
}

//GetPlayer returns the engine's player.
func (engine *Impl) GetPlayer() player.Player {
	return engine.player
}
