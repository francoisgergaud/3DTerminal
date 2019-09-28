package render

import (
	"francoisgergaud/3dGame/environment"

	"github.com/gdamore/tcell"
)

//BackgroundRenderer provides the functionalities to render the environment's map.
type BackgroundRenderer interface {
	Render(worldMap environment.WorldMap, player environment.Character, screen tcell.Screen)
}

//BackgroundRendererImpl implements the BackgroundRenderer interface. For the purpose of code-readability,
//the code has been split between this struct an the BackgroundColumnRenderer struct (internal to this package).
type BackgroundRendererImpl struct {
	screenWidth   int
	bgColRenderer BackgroundColumnRenderer
}

//CreateBackgroundRenderer is a factory:
// - screenWidth: the width of the screen.
// - bgColRenderer: the background's color renderer.
func CreateBackgroundRenderer(screenWidth int, bgColRenderer BackgroundColumnRenderer) *BackgroundRendererImpl {
	return &BackgroundRendererImpl{
		screenWidth:   screenWidth,
		bgColRenderer: bgColRenderer,
	}
}

//Render a scene:
// 1 - clear the screen
// 2 - render each column
// 3 - update the screen
func (bgRenderer *BackgroundRendererImpl) Render(worldMap environment.WorldMap, player environment.Character, screen tcell.Screen) {
	screen.Clear()
	for columnIndex := 0; columnIndex < bgRenderer.screenWidth; columnIndex++ {
		bgRenderer.bgColRenderer.render(screen, player, worldMap, columnIndex)
	}
	screen.Show()
}

//BackgroundColumnRenderer provides functionalities to render a column.
type BackgroundColumnRenderer interface {
	render(screen tcell.Screen, player environment.Character, worldMap environment.WorldMap, columnIndex int)
}

//BackgroundColumnRendererImpl implements the BackgroundColumnRenderer interface.
type BackgroundColumnRendererImpl struct {
	// the screen's dimensions.
	screenWidth, screenHeight int
	// the camera view-angle and maximum-visibility
	fieldOfViewAngle, visibility float64
	// helper for math formula
	mathHelper BackgroundRendererMathHelper
	//style used to render a wall's angle
	wallAngleStyle tcell.Style
	//ray-sampler: contains the styles to be applied for wall and background
	raySampler RaySampler
}

//CreateBackgroundColumnRenderer is a factory: build a BackgroundColumnRenderer
func CreateBackgroundColumnRenderer(screenWidth, screenHeight int, fieldOfViewAngle, visibility float64, mathHelper BackgroundRendererMathHelper, wallAngleStyle tcell.Style, raySampler RaySampler) *BackgroundColumnRendererImpl {
	return &BackgroundColumnRendererImpl{
		screenWidth:      screenWidth,
		screenHeight:     screenHeight,
		fieldOfViewAngle: fieldOfViewAngle,
		visibility:       visibility,
		mathHelper:       mathHelper,
		wallAngleStyle:   wallAngleStyle,
		raySampler:       raySampler,
	}
}

//render renders a column:
// 1 - get the absolute angle of the ray to be casted (from the player's angle and the column-index)
// 2 - cast the ray and find the destination point.
// If there is a wall:
//   3 - get the projection-distance from the player to the destination of the ray-casted (to avoid the "fish-eye" effect.)
//   4 - Get the wall'style (this rendreralso manage the wall's angle to display them in another color)
//   5 - for each row of the column, set the style and rune to be rendered.
func (brColRenderer *BackgroundColumnRendererImpl) render(screen tcell.Screen, player environment.Character, worldMap environment.WorldMap, columnIndex int) {
	//calculate the ray's angle
	rayTracingAngle := brColRenderer.mathHelper.getRayTracingAngleForColumn(player.GetAngle(), columnIndex, brColRenderer.screenWidth, brColRenderer.fieldOfViewAngle)
	//cast the ray
	rayCastDestination := brColRenderer.mathHelper.castRay(player.GetPosition(), worldMap, rayTracingAngle, brColRenderer.visibility)
	if rayCastDestination != nil {
		distance := brColRenderer.mathHelper.calculateProjectionDistance(player.GetPosition(), rayCastDestination, player.GetAngle()-rayTracingAngle)
		//distance := player.Pos.Distance(rayCastDestination)
		var wallStyle tcell.Style
		wallRowStart, wallRowEnd := brColRenderer.mathHelper.GetFillRowRange(distance, float64(brColRenderer.screenHeight))
		isWallAngle := brColRenderer.mathHelper.isWallAngle(rayCastDestination)
		if isWallAngle {
			wallStyle = brColRenderer.wallAngleStyle
		} else {
			wallStyle = brColRenderer.raySampler.GetWallStyleFromDistance(distance)
		}
		for rowIndex := 0; rowIndex < int(brColRenderer.screenHeight); rowIndex++ {
			if rowIndex > wallRowStart && rowIndex < wallRowEnd {
				screen.SetContent(columnIndex, rowIndex, brColRenderer.raySampler.GetWallRune(rowIndex), nil, wallStyle)
			} else {
				screen.SetContent(columnIndex, rowIndex, brColRenderer.raySampler.GetBackgroundRune(rowIndex), nil, brColRenderer.raySampler.GetBackgroundStyle(rowIndex))
			}
		}
	} else {
		for rowIndex := 0; rowIndex < int(brColRenderer.screenHeight); rowIndex++ {
			screen.SetContent(columnIndex, rowIndex, brColRenderer.raySampler.GetBackgroundRune(rowIndex), nil, brColRenderer.raySampler.GetBackgroundStyle(rowIndex))
		}
	}
}
