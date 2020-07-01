package impl

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/projectile"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	internalMath "francoisgergaud/3dGame/common/math"
	"math"
	"testing"

	testRenderMathHelper "francoisgergaud/3dGame/internal/testutils/client/render/mathhelper"
	testAnimatedElement "francoisgergaud/3dGame/internal/testutils/common/environment/animatedelement"
	testprojectile "francoisgergaud/3dGame/internal/testutils/common/environment/projectile"
	testWorld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	testMathHelper "francoisgergaud/3dGame/internal/testutils/common/math/helper"
	testTcell "francoisgergaud/3dGame/internal/testutils/tcell"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWallRendererProducer struct {
	mock.Mock
}

func (mock *MockWallRendererProducer) getRenderer(screen tcell.Screen, player animatedelement.AnimatedElement, worldMap world.WorldMap, columnIndex int) elementRenderer {
	args := mock.Called(screen, player, worldMap, columnIndex)
	return args.Get(0).(elementRenderer)
}

type MockWorldElementRendererProducer struct {
	mock.Mock
}

func (mock *MockWorldElementRendererProducer) getRenderer(player animatedelement.AnimatedElement, fieldOfViewAngle float64, worldElement animatedelement.AnimatedElement) elementRenderer {
	args := mock.Called(player, fieldOfViewAngle, worldElement)
	return args.Get(0).(elementRenderer)
}

type MockElementRenderer struct {
	mock.Mock
}

func (mock *MockElementRenderer) getDistance() float64 {
	args := mock.Called()
	return args.Get(0).(float64)
}

func (mock *MockElementRenderer) render(screen tcell.Screen) {
	mock.Called(screen)
}

func TestCreateRenderer(t *testing.T) {
	screenWidth := 5
	screenHeight := 5
	fieldOfViewAngle := 0.7
	renderMathHelper := new(testRenderMathHelper.MockRendererMathHelper)
	mathHelper := new(testMathHelper.MockMathHelper)
	raySampler := new(MockRaySampler)
	visibility := 5.0
	renderer := CreateRenderer(screenWidth, screenHeight, raySampler, mathHelper, renderMathHelper, fieldOfViewAngle, visibility)
	assert.IsType(t, &RendererImpl{}, renderer)
}

func TestRender(t *testing.T) {
	screen := new(testTcell.MockScreen)
	screenWidth := 5
	screenHeight := 5
	fieldOfViewAngle := 0.7
	renderMathHelper := new(testRenderMathHelper.MockRendererMathHelper)
	wallRendererProducer := new(MockWallRendererProducer)
	worldElementRendererProducer := new(MockWorldElementRendererProducer)
	renderer := createRenderer(screenWidth, screenHeight, renderMathHelper, fieldOfViewAngle, wallRendererProducer, worldElementRendererProducer)
	worldMap := new(testWorld.MockWorldMap)
	player := new(testAnimatedElement.MockAnimatedElement)
	worldElement := new(testAnimatedElement.MockAnimatedElement)
	elementRenderer := new(MockElementRenderer)
	elementRenderer.On("getDistance").Return(0.1)
	elementRenderer.On("render", screen).Times(screenWidth)
	screen.On("Clear")
	for i := 0; i < screenWidth; i++ {
		wallRendererProducer.On("getRenderer", screen, player, worldMap, i).Return(elementRenderer)
	}
	screen.On("Show")
	worldElementRenderer := new(MockElementRenderer)
	worldElementRenderer.On("getDistance").Return(1.1)
	worldElementRenderer.On("render", screen)
	worldElementRendererProducer.On("getRenderer", player, 0.7, worldElement).Return(worldElementRenderer)
	worldElements := make(map[string]animatedelement.AnimatedElement)
	worldElements["worldElementID"] = worldElement
	projectiles := make(map[string]projectile.Projectile)
	projectile := new(testprojectile.MockProjectile)
	projectiles["projectileID"] = projectile
	worldElementRendererProducer.On("getRenderer", player, 0.7, projectile).Return(worldElementRenderer)

	renderer.Render("playerID", worldMap, player, worldElements, projectiles, screen)

	wallRendererProducer.AssertExpectations(t)
	worldElementRendererProducer.AssertExpectations(t)
	worldElementRenderer.AssertExpectations(t)
	elementRenderer.AssertExpectations(t)
	screen.AssertExpectations(t)
}

func TestWallRendererProducer(t *testing.T) {
	screenWidth := 5
	screenHeight := 10
	fieldOfViewAngle := 0.5
	visibility := 5.0
	mathHelper := new(testMathHelper.MockMathHelper)
	rendererMathHelper := new(testRenderMathHelper.MockRendererMathHelper)
	wallAngleStyle := tcell.StyleDefault
	raySampler := new(MockRaySampler)
	wallRendererProducer := createWallRendererProducer(screenWidth, screenHeight, fieldOfViewAngle, visibility, mathHelper, rendererMathHelper, wallAngleStyle, raySampler)
	screen := new(testTcell.MockScreen)
	playerPosition := new(internalMath.Point2D)
	playerAngle := 0.5
	velocity := 1.0
	stepAngle := 0.01
	playerState := &state.AnimatedElementState{
		Position:  playerPosition,
		Angle:     playerAngle,
		Velocity:  velocity,
		StepAngle: stepAngle,
	}
	worldMap := new(testWorld.MockWorldMap)
	player := new(testAnimatedElement.MockAnimatedElement)
	player.On("State").Return(playerState)
	columnIndex := 1
	rayTracingAngle := 0.25
	projectedDistance := 1.5
	rayTracingDestinationPoint := new(internalMath.Point2D)
	startRow := 2
	endRow := 8
	isWallAngle := false
	wallStyle := tcell.StyleDefault.Background(tcell.Color101)
	wallHeight := 1.0
	rendererMathHelper.On("GetRayTracingAngleForColumn", player.State().Angle, columnIndex, screenWidth, fieldOfViewAngle).Return(rayTracingAngle)
	mathHelper.On("CastRay", player.State().Position, worldMap, rayTracingAngle, visibility).Return(rayTracingDestinationPoint).Return(rayTracingDestinationPoint)
	rendererMathHelper.On("CalculateProjectionDistance", playerPosition, rayTracingDestinationPoint, player.State().Angle-rayTracingAngle).Return(projectedDistance)
	rendererMathHelper.On("GetFillRowRange", projectedDistance, visibility, wallHeight, screenHeight).Return(startRow, endRow)
	rendererMathHelper.On("IsWallAngle", rayTracingDestinationPoint).Return(isWallAngle)
	raySampler.On("GetWallStyleFromDistance", projectedDistance).Return(wallStyle)
	wallRendererProducer.getRenderer(screen, player, worldMap, columnIndex)
	mathHelper.AssertExpectations(t)
	rendererMathHelper.AssertExpectations(t)
	raySampler.AssertExpectations(t)
}

func TestWallRendererProducerWithWallAngle(t *testing.T) {
	screenWidth := 5
	screenHeight := 10
	fieldOfViewAngle := 0.5
	visibility := 5.0
	mathHelper := new(testMathHelper.MockMathHelper)
	rendererMathHelper := new(testRenderMathHelper.MockRendererMathHelper)
	wallAngleStyle := tcell.StyleDefault.Background(tcell.Color104)
	raySampler := new(MockRaySampler)
	wallRendererProducer := createWallRendererProducer(screenWidth, screenHeight, fieldOfViewAngle, visibility, mathHelper, rendererMathHelper, wallAngleStyle, raySampler)
	screen := new(testTcell.MockScreen)
	playerPosition := new(internalMath.Point2D)
	playerAngle := 0.5
	velocity := 1.0
	stepAngle := 0.01
	playerState := state.AnimatedElementState{
		Position:  playerPosition,
		Angle:     playerAngle,
		Velocity:  velocity,
		StepAngle: stepAngle,
	}
	worldMap := new(testWorld.MockWorldMap)
	player := new(testAnimatedElement.MockAnimatedElement)
	player.On("State").Return(&playerState)
	columnIndex := 1
	rayTracingAngle := 0.25
	projectedDistance := 1.5
	rayTracingDestinationPoint := new(internalMath.Point2D)
	startRow := 2
	endRow := 8
	isWallAngle := true
	wallHeight := 1.0

	rendererMathHelper.On("GetRayTracingAngleForColumn", player.State().Angle, columnIndex, screenWidth, fieldOfViewAngle).Return(rayTracingAngle)
	mathHelper.On("CastRay", player.State().Position, worldMap, rayTracingAngle, visibility).Return(rayTracingDestinationPoint).Return(rayTracingDestinationPoint)
	rendererMathHelper.On("CalculateProjectionDistance", playerPosition, rayTracingDestinationPoint, player.State().Angle-rayTracingAngle).Return(projectedDistance)
	rendererMathHelper.On("GetFillRowRange", projectedDistance, visibility, wallHeight, screenHeight).Return(startRow, endRow)
	rendererMathHelper.On("IsWallAngle", rayTracingDestinationPoint).Return(isWallAngle)
	wallRendererProducer.getRenderer(screen, player, worldMap, columnIndex)
	mathHelper.AssertExpectations(t)
	rendererMathHelper.AssertExpectations(t)
}

func TestWallRendererProducerWithNilRayTracing(t *testing.T) {
	screenWidth := 5
	screenHeight := 10
	fieldOfViewAngle := 0.5
	visibility := 5.0
	mathHelper := new(testMathHelper.MockMathHelper)
	rendererMathHelper := new(testRenderMathHelper.MockRendererMathHelper)
	wallAngleStyle := tcell.StyleDefault.Background(tcell.Color104)
	raySampler := new(MockRaySampler)
	backgroundColumnRenderer := createWallRendererProducer(screenWidth, screenHeight, fieldOfViewAngle, visibility, mathHelper, rendererMathHelper, wallAngleStyle, raySampler)
	screen := new(testTcell.MockScreen)
	playerPosition := new(internalMath.Point2D)
	playerAngle := 0.5
	velocity := 1.0
	stepAngle := 0.01
	playerState := state.AnimatedElementState{
		Position:  playerPosition,
		Angle:     playerAngle,
		Velocity:  velocity,
		StepAngle: stepAngle,
	}
	worldMap := new(testWorld.MockWorldMap)
	player := new(testAnimatedElement.MockAnimatedElement)
	player.On("State").Return(&playerState)
	columnIndex := 1
	rayTracingAngle := 0.25

	rendererMathHelper.On("GetRayTracingAngleForColumn", player.State().Angle, columnIndex, screenWidth, fieldOfViewAngle).Return(rayTracingAngle)
	mathHelper.On("CastRay", player.State().Position, worldMap, rayTracingAngle, visibility).Return(nil)
	backgroundColumnRenderer.getRenderer(screen, player, worldMap, columnIndex)
	mathHelper.AssertExpectations(t)
}

func TestWallRenderer(t *testing.T) {
	wallRowStart := 3
	wallRowEnd := 7
	screenHeight := 10
	wallStyle := tcell.StyleDefault.Background(tcell.Color108)
	wallRune := '3'
	raySampler := new(MockRaySampler)
	columnIndex := 9
	distance := 5.9
	wallRenderer := wallRenderer{
		distance:     distance,
		columnIndex:  columnIndex,
		screenHeight: screenHeight,
		wallRowStart: wallRowStart,
		wallRowEnd:   wallRowEnd,
		wallStyle:    wallStyle,
		raySampler:   raySampler,
	}
	screen := new(testTcell.MockScreen)
	backgroundStyle := tcell.StyleDefault.Background(tcell.Color102)
	backgroundRune := '2'
	for rowIndex := 0; rowIndex < screenHeight; rowIndex++ {
		if rowIndex <= wallRowStart || rowIndex >= wallRowEnd {
			raySampler.On("GetBackgroundRune", rowIndex).Return(backgroundRune)
			raySampler.On("GetBackgroundStyle", rowIndex).Return(backgroundStyle)
			screen.On("SetContent", columnIndex, rowIndex, backgroundRune, []int32(nil), backgroundStyle)
		} else {
			raySampler.On("GetWallRune", rowIndex).Return(wallRune)
			screen.On("SetContent", columnIndex, rowIndex, wallRune, []int32(nil), wallStyle)
		}
	}
	wallRenderer.render(screen)
	screen.AssertExpectations(t)
	raySampler.AssertExpectations(t)
	assert.Equal(t, distance, wallRenderer.getDistance())
}

func TestWorldElementRendererProducerImpl(t *testing.T) {
	playerPosition := &internalMath.Point2D{X: 0, Y: 5}
	playerAngle := 0.0
	worldElementPosition := &internalMath.Point2D{X: 5, Y: 5}
	worldElementStyle := tcell.StyleDefault.Background(tcell.Color107)
	distance := playerPosition.Distance(worldElementPosition)
	worldElementSize := 0.5
	fieldOfView := 0.5
	screenHeight := 10
	screenWidth := 10
	isVisible := true
	startScreenWidthRatio := 0.4365
	startOffset := 0.0
	endScreenWidthRatio := 0.5634
	endOffset := 1.0
	worldElementRowStart := 3
	worldElementRowEnd := 7
	maxVisibility := 10.0
	wallHeight := 1.0
	mathHelper := new(testMathHelper.MockMathHelper)
	rendererMathHelper := new(testRenderMathHelper.MockRendererMathHelper)
	mathHelper.On("GetWorldElementProjection", playerPosition, playerAngle, fieldOfView, worldElementPosition, worldElementSize).Return(isVisible, startScreenWidthRatio, startOffset, endScreenWidthRatio, endOffset)
	rendererMathHelper.On("GetFillRowRange", distance, maxVisibility, wallHeight, screenHeight).Return(worldElementRowStart, worldElementRowEnd)
	worldElementRendererProducer := createWorldElementRendererProducer(mathHelper, rendererMathHelper, screenHeight, screenWidth, maxVisibility)
	worldElement := new(testAnimatedElement.MockAnimatedElement)
	worldElementState := state.AnimatedElementState{
		Position: worldElementPosition,
		Style:    worldElementStyle,
		Size:     worldElementSize,
	}
	worldElement.On("State").Return(&worldElementState)
	player := new(testAnimatedElement.MockAnimatedElement)
	playerState := state.AnimatedElementState{
		Position: playerPosition,
		Angle:    playerAngle,
	}
	player.On("State").Return(&playerState)
	worldElementRenderer := worldElementRendererProducer.getRenderer(player, fieldOfView, worldElement).(*worldElementRenderer)
	assert.Equal(t, worldElementRenderer.distance, distance)
	assert.Equal(t, worldElementRenderer.screenHeight, screenHeight)
	assert.Equal(t, worldElementRenderer.screenWidth, float64(screenWidth))
	assert.Equal(t, worldElementRenderer.worldElementRowStart, worldElementRowStart)
	assert.Equal(t, worldElementRenderer.worldElementRowEnd, worldElementRowEnd)
	assert.Equal(t, worldElementRenderer.startScreenWidthRatio, startScreenWidthRatio)
	assert.Equal(t, worldElementRenderer.startWorldElementOffset, startOffset)
	assert.Equal(t, worldElementRenderer.endScreenWidthRatio, endScreenWidthRatio)
	assert.Equal(t, worldElementRenderer.endtWorldElementOffset, endOffset)
	assert.Equal(t, worldElementRenderer.worldElementStyle, worldElementStyle)
	rendererMathHelper.AssertExpectations(t)
	worldElement.AssertExpectations(t)
	player.AssertExpectations(t)
}

func TestWorldElementRenderer(t *testing.T) {
	startScreenWidthRatio := 0.4365
	startOffset := 0.0
	endScreenWidthRatio := 0.5634
	endOffset := 1.0
	screenHeight := 10
	screenWidth := 10.0
	worldElementRowStart := 3
	worldElementRowEnd := 7
	distance := 8.3
	worldElementStyle := tcell.StyleDefault.Background(tcell.Color108)
	worldElementRenderer := worldElementRenderer{
		distance:                distance,
		screenHeight:            screenHeight,
		screenWidth:             screenWidth,
		worldElementRowStart:    worldElementRowStart,
		worldElementRowEnd:      worldElementRowEnd,
		startScreenWidthRatio:   startScreenWidthRatio,
		startWorldElementOffset: startOffset,
		endScreenWidthRatio:     endScreenWidthRatio,
		endtWorldElementOffset:  endOffset,
		worldElementStyle:       worldElementStyle,
	}
	screen := new(testTcell.MockScreen)
	for column := int(math.Round(startScreenWidthRatio * screenWidth)); column <= int(math.Round(endScreenWidthRatio*screenWidth)); column++ {
		for row := worldElementRowStart; row <= worldElementRowEnd; row++ {
			screen.On("SetContent", column, row, ' ', []int32(nil), worldElementStyle)
		}
	}
	worldElementRenderer.render(screen)
	screen.AssertExpectations(t)
	assert.Equal(t, distance, worldElementRenderer.getDistance())
}

func TestBackgroundRenderer(t *testing.T) {
	screenHeight := 10
	raySampler := new(MockRaySampler)
	columnIndex := 4
	backgroundRenderer := backgroundRenderer{
		screenHeight: screenHeight,
		raySampler:   raySampler,
		columnIndex:  columnIndex,
	}
	screen := new(testTcell.MockScreen)
	backgroundStyle := tcell.StyleDefault.Background(tcell.Color102)
	backgroundRune := '2'
	for rowIndex := 0; rowIndex < screenHeight; rowIndex++ {
		raySampler.On("GetBackgroundRune", rowIndex).Return(backgroundRune)
		raySampler.On("GetBackgroundStyle", rowIndex).Return(backgroundStyle)
		screen.On("SetContent", columnIndex, rowIndex, backgroundRune, []int32(nil), backgroundStyle)
	}
	backgroundRenderer.render(screen)
	screen.AssertExpectations(t)
	raySampler.AssertExpectations(t)
	assert.Equal(t, math.Inf(1), backgroundRenderer.getDistance())
}
