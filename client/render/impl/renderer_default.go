package impl

import (
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/client/render"
	"francoisgergaud/3dGame/client/render/mathhelper"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/world"
	commonMathHelper "francoisgergaud/3dGame/common/math/helper"
	"math"
	"sort"

	"github.com/gdamore/tcell"
)

//RendererImpl implements the Renderer interface. For the purpose of code-readability,
//the code has been split between this struct and the RendererProducer (internal to this package).
type RendererImpl struct {
	// the screen's height and width
	screenWidth, screenHeight int
	// helper for math formula
	renderMathHelper mathhelper.RendererMathHelper
	//the wall-renderer-producer
	wallRendererProducer wallRendererProducer
	//field-of-view angle.
	fieldOfViewAngle float64
	//the world-element-renderer-producer
	worldElementRendererProducer worldElementRendererProducer
}

//CreateRenderer is a factory:
func CreateRenderer(screenWidth, screenHeight int, raySampler RaySampler, mathHelper commonMathHelper.MathHelper, renderMathHelper mathhelper.RendererMathHelper, fieldOfViewAngle, visibility float64) render.Renderer {
	wallRendererProducer := createWallRendererProducer(screenWidth, screenHeight, fieldOfViewAngle, visibility, mathHelper, renderMathHelper, tcell.StyleDefault.Background(tcell.ColorBlueViolet), raySampler)
	worldElementRendererProducer := createWorldElementRendererProducer(mathHelper, renderMathHelper, screenHeight, screenWidth)
	return createRenderer(screenWidth, screenHeight, renderMathHelper, fieldOfViewAngle, wallRendererProducer, worldElementRendererProducer)
}

func createRenderer(screenWidth, screenHeight int, renderMathHelper mathhelper.RendererMathHelper, fieldOfViewAngle float64, wallRendererProducer wallRendererProducer, worldElementRendererProducer worldElementRendererProducer) render.Renderer {
	return &RendererImpl{
		screenWidth:                  screenWidth,
		screenHeight:                 screenHeight,
		wallRendererProducer:         wallRendererProducer,
		worldElementRendererProducer: worldElementRendererProducer,
		renderMathHelper:             renderMathHelper,
		fieldOfViewAngle:             fieldOfViewAngle,
	}
}

//Render a scene:
// 1 - clear the screen
// 2 - get the wall/background and world-element renderers each column
// 3 - sort these renderers by depth
// 4 - render each renderer from the deepest to the nearest.
// 5 - update the screen
func (renderer *RendererImpl) Render(playerID string, worldMap world.WorldMap, player player.Player, worldElements map[string]animatedelement.AnimatedElement, screen tcell.Screen) {
	screen.Clear()
	renderers := make([]elementRenderer, 0)
	for columnIndex := 0; columnIndex < renderer.screenWidth; columnIndex++ {
		renderers = append(renderers, renderer.wallRendererProducer.getRenderer(screen, player, worldMap, columnIndex))
	}
	if worldElements != nil {
		for worldElementID, worldElement := range worldElements {
			if worldElementID != playerID {
				worldElementRenderer := renderer.worldElementRendererProducer.getRenderer(player, renderer.fieldOfViewAngle, worldElement)
				if worldElementRenderer != nil {
					renderers = append(renderers, worldElementRenderer)
				}
			}
		}
	}
	// sort the 'elementRenderers' array by their ditance (from grater to lower) and render them.
	sort.Slice(renderers, func(e1, e2 int) bool {
		return renderers[e1].getDistance() > renderers[e2].getDistance()
	})
	for _, elementRenderer := range renderers {
		elementRenderer.render(screen)
	}
	screen.Show()
}

//wallRendererProducer provides functionalities to produce a wall-and-background renderer.
type wallRendererProducer interface {
	getRenderer(screen tcell.Screen, player player.Player, worldMap world.WorldMap, columnIndex int) elementRenderer
}

//wallRendererProducerImpl implements the WallRendererProducer interface.
type wallRendererProducerImpl struct {
	// the screen's dimensions.
	screenWidth, screenHeight int
	// the camera view-angle and maximum-visibility
	fieldOfViewAngle, visibility float64
	// helper for math formula for rendering
	renderMathHelper mathhelper.RendererMathHelper
	// helper for math formula
	mathHelper commonMathHelper.MathHelper
	//style used to render a wall's angle
	wallAngleStyle tcell.Style
	//ray-sampler: contains the styles to be applied for wall and background
	raySampler RaySampler
}

//createWallRendererProducer is a factory: build a WallRendererProducer
func createWallRendererProducer(screenWidth, screenHeight int, fieldOfViewAngle, visibility float64, mathHelper commonMathHelper.MathHelper, renderMathHelper mathhelper.RendererMathHelper, wallAngleStyle tcell.Style, raySampler RaySampler) *wallRendererProducerImpl {
	return &wallRendererProducerImpl{
		screenWidth:      screenWidth,
		screenHeight:     screenHeight,
		fieldOfViewAngle: fieldOfViewAngle,
		visibility:       visibility,
		renderMathHelper: renderMathHelper,
		mathHelper:       mathHelper,
		wallAngleStyle:   wallAngleStyle,
		raySampler:       raySampler,
	}
}

//getRenderer get teh rendering-data for a wall/background:
// 1 - get the absolute angle of the ray to be casted (from the player's angle and the column-index)
// 2 - cast the ray and find the destination point.
// If there is a wall:
//   3 - get the projection-distance from the player to the destination of the ray-casted (to avoid the "fish-eye" effect.)
//   4 - Get the wall'style (this rendreralso manage the wall's angle to display them in another color)
//   5 - for each row of the column, set the style and rune to be rendered.
func (wallRendererProducer *wallRendererProducerImpl) getRenderer(screen tcell.Screen, player player.Player, worldMap world.WorldMap, columnIndex int) elementRenderer {
	playerState := player.GetState()
	//calculate the ray's angle
	rayTracingAngle := wallRendererProducer.renderMathHelper.GetRayTracingAngleForColumn(playerState.Angle, columnIndex, wallRendererProducer.screenWidth, wallRendererProducer.fieldOfViewAngle)
	//cast the ray
	rayCastDestination := wallRendererProducer.mathHelper.CastRay(playerState.Position, worldMap, rayTracingAngle, wallRendererProducer.visibility)
	if rayCastDestination != nil {
		distance := wallRendererProducer.renderMathHelper.CalculateProjectionDistance(playerState.Position, rayCastDestination, playerState.Angle-rayTracingAngle)
		var wallStyle tcell.Style
		wallRowStart, wallRowEnd := wallRendererProducer.renderMathHelper.GetFillRowRange(distance, float64(wallRendererProducer.screenHeight))
		isWallAngle := wallRendererProducer.renderMathHelper.IsWallAngle(rayCastDestination)
		if isWallAngle {
			wallStyle = wallRendererProducer.wallAngleStyle
		} else {
			wallStyle = wallRendererProducer.raySampler.GetWallStyleFromDistance(distance)
		}
		return &wallRenderer{
			distance:     distance,
			columnIndex:  columnIndex,
			wallRowStart: wallRowStart,
			wallRowEnd:   wallRowEnd,
			wallStyle:    wallStyle,
			raySampler:   wallRendererProducer.raySampler,
			screenHeight: wallRendererProducer.screenHeight,
		}
	}
	return &backgroundRenderer{
		columnIndex:  columnIndex,
		raySampler:   wallRendererProducer.raySampler,
		screenHeight: wallRendererProducer.screenHeight,
	}
}

//worldElementRendererProducer provides functionalities to produce a wall-and-background renderer.
type worldElementRendererProducer interface {
	getRenderer(player player.Player, fieldOfViewAngle float64, worldElement animatedelement.AnimatedElement) elementRenderer
}

//worldElementRendererProducerImpl implements the WorldElementRendererProducer.
type worldElementRendererProducerImpl struct {
	// helper for math formula for rendering
	renderMathHelper mathhelper.RendererMathHelper
	// helper for math formula
	mathHelper commonMathHelper.MathHelper
	//the screen-height.
	screenHeight int
	//the screen-width.
	screenWidth int
}

//createWorldElementRendererProducer creates a WorldElementRendererProducer.
func createWorldElementRendererProducer(mathHelper commonMathHelper.MathHelper, rendererMathHelper mathhelper.RendererMathHelper, screenHeight, screenWidth int) worldElementRendererProducer {
	return &worldElementRendererProducerImpl{
		mathHelper:       mathHelper,
		renderMathHelper: rendererMathHelper,
		screenHeight:     screenHeight,
		screenWidth:      screenWidth,
	}
}

func (WorldElementRendererProducer *worldElementRendererProducerImpl) getRenderer(player player.Player, fieldOfViewAngle float64, worldElement animatedelement.AnimatedElement) elementRenderer {
	playerState := player.GetState()
	worldElementState := worldElement.GetState()
	isVisible, startScreenWidthRatio, startOffset, endScreenWidthRatio, endOffset := WorldElementRendererProducer.mathHelper.GetWorldElementProjection(playerState.Position, playerState.Angle, fieldOfViewAngle, worldElementState.Position, worldElementState.Size)
	if isVisible {
		distance := playerState.Position.Distance(worldElement.GetState().Position)
		worldElementRowStart, worldElementRowEnd := WorldElementRendererProducer.renderMathHelper.GetFillRowRange(distance, float64(WorldElementRendererProducer.screenHeight))
		return &worldElementRenderer{
			distance:                distance,
			screenHeight:            WorldElementRendererProducer.screenHeight,
			screenWidth:             float64(WorldElementRendererProducer.screenWidth),
			worldElementRowStart:    worldElementRowStart,
			worldElementRowEnd:      worldElementRowEnd,
			startScreenWidthRatio:   startScreenWidthRatio,
			startWorldElementOffset: startOffset,
			endScreenWidthRatio:     endScreenWidthRatio,
			endtWorldElementOffset:  endOffset,
			worldElementStyle:       worldElementState.Style,
		}
	}
	return nil
}

type elementRenderer interface {
	getDistance() float64
	render(screen tcell.Screen)
}

type wallRenderer struct {
	distance     float64
	columnIndex  int
	wallRowStart int
	wallRowEnd   int
	wallStyle    tcell.Style
	raySampler   RaySampler
	screenHeight int
}

func (wallRenderer *wallRenderer) render(screen tcell.Screen) {
	for rowIndex := 0; rowIndex < int(wallRenderer.screenHeight); rowIndex++ {
		if rowIndex > wallRenderer.wallRowStart && rowIndex < wallRenderer.wallRowEnd {
			screen.SetContent(wallRenderer.columnIndex, rowIndex, wallRenderer.raySampler.GetWallRune(rowIndex), nil, wallRenderer.wallStyle)
		} else {
			screen.SetContent(wallRenderer.columnIndex, rowIndex, wallRenderer.raySampler.GetBackgroundRune(rowIndex), nil, wallRenderer.raySampler.GetBackgroundStyle(rowIndex))
		}
	}
}

func (wallRenderer *wallRenderer) getDistance() float64 {
	return wallRenderer.distance
}

type backgroundRenderer struct {
	columnIndex  int
	raySampler   RaySampler
	screenHeight int
}

func (backgroundRenderer *backgroundRenderer) render(screen tcell.Screen) {
	for rowIndex := 0; rowIndex < int(backgroundRenderer.screenHeight); rowIndex++ {
		screen.SetContent(backgroundRenderer.columnIndex, rowIndex, backgroundRenderer.raySampler.GetBackgroundRune(rowIndex), nil, backgroundRenderer.raySampler.GetBackgroundStyle(rowIndex))
	}
}

func (backgroundRenderer *backgroundRenderer) getDistance() float64 {
	return math.Inf(1)
}

type worldElementRenderer struct {
	distance                float64
	screenHeight            int
	screenWidth             float64
	worldElementRowStart    int
	worldElementRowEnd      int
	startScreenWidthRatio   float64
	startWorldElementOffset float64
	endScreenWidthRatio     float64
	endtWorldElementOffset  float64
	worldElementStyle       tcell.Style
}

func (worldElementRenderer *worldElementRenderer) render(screen tcell.Screen) {
	columnStart := int(math.Round(worldElementRenderer.screenWidth * worldElementRenderer.startScreenWidthRatio))
	columnEnd := int(math.Round(worldElementRenderer.screenWidth * worldElementRenderer.endScreenWidthRatio))
	for columnIndex := columnStart; columnIndex <= columnEnd; columnIndex++ {
		for rowIndex := worldElementRenderer.worldElementRowStart; rowIndex <= worldElementRenderer.worldElementRowEnd; rowIndex++ {
			screen.SetContent(columnIndex, rowIndex, ' ', nil, worldElementRenderer.worldElementStyle)
		}
	}
}

func (worldElementRenderer *worldElementRenderer) getDistance() float64 {
	return worldElementRenderer.distance
}
