package render

import (
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/environment/character"
	"francoisgergaud/3dGame/environment/world"
	"francoisgergaud/3dGame/environment/worldelement"
	"francoisgergaud/3dGame/internal/testutils"
	testcharacter "francoisgergaud/3dGame/internal/testutils/environment/character"
	testworldelement "francoisgergaud/3dGame/internal/testutils/environment/worldelement"
	"francoisgergaud/3dGame/internal/testutils/environment/worldmap"
	"math"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//MockBackgroundRendererMathHelper mocks the calls to the BackgroundRendererMathHelper interface.
type MockRendererMathHelper struct {
	mock.Mock
}

func (mock *MockRendererMathHelper) calculateProjectionDistance(startPosition *common.Point2D, endPosition *common.Point2D, angle float64) float64 {
	args := mock.Called(startPosition, endPosition, angle)
	return args.Get(0).(float64)
}
func (mock *MockRendererMathHelper) isWallAngle(point *common.Point2D) bool {
	args := mock.Called(point)
	return args.Bool(0)
}
func (mock *MockRendererMathHelper) getRayTracingAngleForColumn(angle float64, columnIndex, screenWidth int, viewAngle float64) float64 {
	args := mock.Called(angle, columnIndex, screenWidth, viewAngle)
	return args.Get(0).(float64)
}

func (mock *MockRendererMathHelper) getFillRowRange(distance, screenHeight float64) (int, int) {
	args := mock.Called(distance, screenHeight)
	return args.Int(0), args.Int(1)
}

type MockWallRendererProducer struct {
	mock.Mock
}

func (mock *MockWallRendererProducer) getRenderer(screen tcell.Screen, player character.Character, worldMap world.WorldMap, columnIndex int) elementRenderer {
	args := mock.Called(screen, player, worldMap, columnIndex)
	return args.Get(0).(elementRenderer)
}

type MockWorldElementRendererProducer struct {
	mock.Mock
}

func (mock *MockWorldElementRendererProducer) getRenderer(player character.Character, fieldOfViewAngle float64, worldElement worldelement.WorldElement) elementRenderer {
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

func TestRendererImpl(t *testing.T) {
	screen := new(testutils.MockScreen)
	screenWidth := 5
	screenHeight := 5
	fieldOfViewAngle := 0.7
	renderMathHelper := new(MockRendererMathHelper)
	wallRendererProducer := new(MockWallRendererProducer)
	worldElementRendererProducer := new(MockWorldElementRendererProducer)
	renderer := CreateRenderer(screenWidth, screenHeight, renderMathHelper, fieldOfViewAngle, wallRendererProducer, worldElementRendererProducer)
	worldMap := new(worldmap.MockWorldMap)
	player := character.NewPlayableCharacter(nil, 0.0, 1.0, 0.01, worldMap)
	worldElement := new(testworldelement.MockWorldElement)
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
	renderer.Render(worldMap, player, []worldelement.WorldElement{worldElement}, screen)
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
	mathHelper := new(testutils.MockMathHelper)
	rendererMathHelper := new(MockRendererMathHelper)
	wallAngleStyle := tcell.StyleDefault
	raySampler := new(testutils.MockRaySampler)
	wallRendererProducer := CreateWallRendererProducer(screenWidth, screenHeight, fieldOfViewAngle, visibility, mathHelper, rendererMathHelper, wallAngleStyle, raySampler)
	screen := new(testutils.MockScreen)
	playerPosition := new(common.Point2D)
	playerAngle := 0.5
	velocity := 1.0
	stepAngle := 0.01
	worlMap := new(worldmap.MockWorldMap)
	player := character.NewPlayableCharacter(playerPosition, playerAngle, velocity, stepAngle, worlMap)
	columnIndex := 1
	rayTracingAngle := 0.25
	projectedDistance := 1.5
	rayTracingDestinationPoint := new(common.Point2D)
	startRow := 2
	endRow := 8
	isWallAngle := false
	wallStyle := tcell.StyleDefault.Background(tcell.Color101)

	rendererMathHelper.On("getRayTracingAngleForColumn", player.GetAngle(), columnIndex, screenWidth, fieldOfViewAngle).Return(rayTracingAngle)
	mathHelper.On("CastRay", player.GetPosition(), worlMap, rayTracingAngle, visibility).Return(rayTracingDestinationPoint).Return(rayTracingDestinationPoint)
	rendererMathHelper.On("calculateProjectionDistance", playerPosition, rayTracingDestinationPoint, player.GetAngle()-rayTracingAngle).Return(projectedDistance)
	rendererMathHelper.On("getFillRowRange", projectedDistance, float64(screenHeight)).Return(startRow, endRow)
	rendererMathHelper.On("isWallAngle", rayTracingDestinationPoint).Return(isWallAngle)
	raySampler.On("GetWallStyleFromDistance", projectedDistance).Return(wallStyle)
	wallRendererProducer.getRenderer(screen, player, worlMap, columnIndex)
	mathHelper.AssertExpectations(t)
	rendererMathHelper.AssertExpectations(t)
	raySampler.AssertExpectations(t)
}

func TestWallRendererProducerWithWallAngle(t *testing.T) {
	screenWidth := 5
	screenHeight := 10
	fieldOfViewAngle := 0.5
	visibility := 5.0
	mathHelper := new(testutils.MockMathHelper)
	rendererMathHelper := new(MockRendererMathHelper)
	wallAngleStyle := tcell.StyleDefault.Background(tcell.Color104)
	raySampler := new(testutils.MockRaySampler)
	wallRendererProducer := CreateWallRendererProducer(screenWidth, screenHeight, fieldOfViewAngle, visibility, mathHelper, rendererMathHelper, wallAngleStyle, raySampler)
	screen := new(testutils.MockScreen)
	playerPosition := new(common.Point2D)
	playerAngle := 0.5
	worlMap := new(worldmap.MockWorldMap)
	velocity := 1.0
	stepAngle := 0.01
	player := character.NewPlayableCharacter(playerPosition, playerAngle, velocity, stepAngle, worlMap)
	columnIndex := 1
	rayTracingAngle := 0.25
	projectedDistance := 1.5
	rayTracingDestinationPoint := new(common.Point2D)
	startRow := 2
	endRow := 8
	isWallAngle := true

	rendererMathHelper.On("getRayTracingAngleForColumn", player.GetAngle(), columnIndex, screenWidth, fieldOfViewAngle).Return(rayTracingAngle)
	mathHelper.On("CastRay", player.GetPosition(), worlMap, rayTracingAngle, visibility).Return(rayTracingDestinationPoint).Return(rayTracingDestinationPoint)
	rendererMathHelper.On("calculateProjectionDistance", playerPosition, rayTracingDestinationPoint, player.GetAngle()-rayTracingAngle).Return(projectedDistance)
	rendererMathHelper.On("getFillRowRange", projectedDistance, float64(screenHeight)).Return(startRow, endRow)
	rendererMathHelper.On("isWallAngle", rayTracingDestinationPoint).Return(isWallAngle)
	wallRendererProducer.getRenderer(screen, player, worlMap, columnIndex)
	mathHelper.AssertExpectations(t)
	rendererMathHelper.AssertExpectations(t)
}

func TestWallRendererProducerWithNilRayTracing(t *testing.T) {
	screenWidth := 5
	screenHeight := 10
	fieldOfViewAngle := 0.5
	visibility := 5.0
	mathHelper := new(testutils.MockMathHelper)
	rendererMathHelper := new(MockRendererMathHelper)
	wallAngleStyle := tcell.StyleDefault.Background(tcell.Color104)
	raySampler := new(testutils.MockRaySampler)
	backgroundColumnRenderer := CreateWallRendererProducer(screenWidth, screenHeight, fieldOfViewAngle, visibility, mathHelper, rendererMathHelper, wallAngleStyle, raySampler)
	screen := new(testutils.MockScreen)
	playerPosition := new(common.Point2D)
	playerAngle := 0.5
	velocity := 1.0
	stepAngle := 0.01
	worlMap := new(worldmap.MockWorldMap)
	player := character.NewPlayableCharacter(playerPosition, playerAngle, velocity, stepAngle, worlMap)
	columnIndex := 1
	rayTracingAngle := 0.25

	rendererMathHelper.On("getRayTracingAngleForColumn", player.GetAngle(), columnIndex, screenWidth, fieldOfViewAngle).Return(rayTracingAngle)
	mathHelper.On("CastRay", player.GetPosition(), worlMap, rayTracingAngle, visibility).Return(nil)
	backgroundColumnRenderer.getRenderer(screen, player, worlMap, columnIndex)
	mathHelper.AssertExpectations(t)
}

func TestWallRenderer(t *testing.T) {
	wallRowStart := 3
	wallRowEnd := 7
	screenHeight := 10
	wallStyle := tcell.StyleDefault.Background(tcell.Color108)
	wallRune := '3'
	raySampler := new(testutils.MockRaySampler)
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
	screen := new(testutils.MockScreen)
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
	playerPosition := &common.Point2D{X: 0, Y: 5}
	playerAngle := 0.0
	worlElementPosition := &common.Point2D{X: 5, Y: 5}
	worldElementStyle := tcell.StyleDefault.Background(tcell.Color107)
	distance := playerPosition.Distance(worlElementPosition)
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
	mathHelper := new(testutils.MockMathHelper)
	rendererMathHelper := new(MockRendererMathHelper)
	mathHelper.On("GetWorldElementProjection", playerPosition, playerAngle, fieldOfView, worlElementPosition, worldElementSize).Return(isVisible, startScreenWidthRatio, startOffset, endScreenWidthRatio, endOffset)
	rendererMathHelper.On("getFillRowRange", distance, float64(screenHeight)).Return(worldElementRowStart, worldElementRowEnd)
	worldElementRendererProducer := CreateWorldElementRendererProducer(mathHelper, rendererMathHelper, screenHeight, screenWidth)
	worldElement := new(testworldelement.MockWorldElement)
	worldElement.On("GetPosition").Return(worlElementPosition)
	worldElement.On("GetSize").Return(worldElementSize)
	worldElement.On("GetStyle").Return(worldElementStyle)
	player := new(testcharacter.MockCharacter)
	player.On("GetPosition").Return(playerPosition)
	player.On("GetAngle").Return(playerAngle)
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
	screen := new(testutils.MockScreen)
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
	raySampler := new(testutils.MockRaySampler)
	columnIndex := 4
	backgroundRenderer := backgroundRenderer{
		screenHeight: screenHeight,
		raySampler:   raySampler,
		columnIndex:  columnIndex,
	}
	screen := new(testutils.MockScreen)
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
