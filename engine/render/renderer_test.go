package render

import (
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/environment"
	"francoisgergaud/3dGame/internal/testutils"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/mock"
)

//MockBackgroundRendererMathHelper mocks the calls to the BackgroundRendererMathHelper interface.
type MockBackgroundRendererMathHelper struct {
	mock.Mock
}

func (mock *MockBackgroundRendererMathHelper) calculateProjectionDistance(startPosition *common.Point2D, endPosition *common.Point2D, angle float64) float64 {
	args := mock.Called(startPosition, endPosition, angle)
	return args.Get(0).(float64)
}
func (mock *MockBackgroundRendererMathHelper) isWallAngle(point *common.Point2D) bool {
	args := mock.Called(point)
	return args.Bool(0)
}
func (mock *MockBackgroundRendererMathHelper) getRayTracingAngleForColumn(angle float64, columnIndex, screenWidth int, viewAngle float64) float64 {
	args := mock.Called(angle, columnIndex, screenWidth, viewAngle)
	return args.Get(0).(float64)
}
func (mock *MockBackgroundRendererMathHelper) castRay(origin *common.Point2D, worldMap environment.WorldMap, rayAngle, visibility float64) *common.Point2D {
	args := mock.Called(origin, worldMap, rayAngle, visibility)
	if args.Get(0) != nil {
		return args.Get(0).(*common.Point2D)
	}
	return nil
}
func (mock *MockBackgroundRendererMathHelper) GetFillRowRange(distance, screenHeight float64) (int, int) {
	args := mock.Called(distance, screenHeight)
	return args.Int(0), args.Int(1)
}

type MockBackgroundColumnRenderer struct {
	mock.Mock
}

func (mock *MockBackgroundColumnRenderer) render(screen tcell.Screen, player environment.Character, worldMap environment.WorldMap, columnIndex int) {
	mock.Called(screen, player, worldMap, columnIndex)
}

func TestBackgroundRenderer(t *testing.T) {
	screen := new(testutils.MockScreen)
	screenWidth := 5
	bgColRenderer := new(MockBackgroundColumnRenderer)
	bgRenderer := CreateBackgroundRenderer(screenWidth, bgColRenderer)
	worldMap := new(testutils.MockWorldMap)
	player := environment.NewPlayer(nil, 0.0, 0.0, worldMap)
	screen.On("Clear")
	for i := 0; i < screenWidth; i++ {
		bgColRenderer.On("render", screen, player, worldMap, i)
	}
	screen.On("Show")
	bgRenderer.Render(worldMap, player, screen)
}

func TestBackgroundColumnRenderer(t *testing.T) {
	screenWidth := 5
	screenHeight := 10
	fieldOfViewAngle := 0.5
	visibility := 5.0
	mathHelper := new(MockBackgroundRendererMathHelper)
	wallAngleStyle := tcell.StyleDefault
	raySampler := new(testutils.MockRaySampler)
	backgroundColumnRenderer := CreateBackgroundColumnRenderer(screenWidth, screenHeight, fieldOfViewAngle, visibility, mathHelper, wallAngleStyle, raySampler)
	screen := new(testutils.MockScreen)
	playerPosition := new(common.Point2D)
	playerAngle := 0.5
	worlMap := new(testutils.MockWorldMap)
	player := environment.NewPlayer(playerPosition, playerAngle, 1.0, worlMap)
	columnIndex := 1
	rayTracingAngle := 0.25
	projectedDistance := 1.5
	rayTracingDestinationPoint := new(common.Point2D)
	startRow := 2
	endRow := 8
	isWallAngle := false
	wallStyle := tcell.StyleDefault.Background(tcell.Color101)
	wallRune := '1'
	backgroundStyle := tcell.StyleDefault.Background(tcell.Color102)
	backgroundRune := '2'

	mathHelper.On("getRayTracingAngleForColumn", player.GetAngle(), columnIndex, screenWidth, fieldOfViewAngle).Return(rayTracingAngle)
	mathHelper.On("castRay", player.GetPosition(), worlMap, rayTracingAngle, visibility).Return(rayTracingDestinationPoint).Return(rayTracingDestinationPoint)
	mathHelper.On("calculateProjectionDistance", playerPosition, rayTracingDestinationPoint, player.GetAngle()-rayTracingAngle).Return(projectedDistance)
	mathHelper.On("GetFillRowRange", projectedDistance, float64(screenHeight)).Return(startRow, endRow)
	mathHelper.On("isWallAngle", rayTracingDestinationPoint).Return(isWallAngle)
	raySampler.On("GetWallStyleFromDistance", projectedDistance).Return(wallStyle)
	for rowIndex := 0; rowIndex <= startRow; rowIndex++ {
		raySampler.On("GetBackgroundRune", rowIndex).Return(backgroundRune)
		raySampler.On("GetBackgroundStyle", rowIndex).Return(backgroundStyle)
		screen.On("SetContent", columnIndex, rowIndex, backgroundRune, []int32(nil), backgroundStyle)
	}
	for rowIndex := startRow + 1; rowIndex < endRow; rowIndex++ {
		raySampler.On("GetWallRune", rowIndex).Return(wallRune)
		screen.On("SetContent", columnIndex, rowIndex, wallRune, []int32(nil), wallStyle)
	}
	for rowIndex := endRow; rowIndex <= screenHeight; rowIndex++ {
		raySampler.On("GetBackgroundRune", rowIndex).Return(backgroundRune)
		raySampler.On("GetBackgroundStyle", rowIndex).Return(backgroundStyle)
		screen.On("SetContent", columnIndex, rowIndex, backgroundRune, []int32(nil), backgroundStyle)
	}
	backgroundColumnRenderer.render(screen, player, worlMap, columnIndex)
	mathHelper.AssertExpectations(t)
}

func TestBackgroundColumnRendererWallAngle(t *testing.T) {
	screenWidth := 5
	screenHeight := 10
	fieldOfViewAngle := 0.5
	visibility := 5.0
	mathHelper := new(MockBackgroundRendererMathHelper)
	wallAngleStyle := tcell.StyleDefault.Background(tcell.Color104)
	raySampler := new(testutils.MockRaySampler)
	backgroundColumnRenderer := CreateBackgroundColumnRenderer(screenWidth, screenHeight, fieldOfViewAngle, visibility, mathHelper, wallAngleStyle, raySampler)
	screen := new(testutils.MockScreen)
	playerPosition := new(common.Point2D)
	playerAngle := 0.5
	worlMap := new(testutils.MockWorldMap)
	player := environment.NewPlayer(playerPosition, playerAngle, 1.0, worlMap)
	columnIndex := 1
	rayTracingAngle := 0.25
	projectedDistance := 1.5
	rayTracingDestinationPoint := new(common.Point2D)
	startRow := 2
	endRow := 8
	isWallAngle := true
	wallRune := '1'
	backgroundStyle := tcell.StyleDefault.Background(tcell.Color102)
	backgroundRune := '2'

	mathHelper.On("getRayTracingAngleForColumn", player.GetAngle(), columnIndex, screenWidth, fieldOfViewAngle).Return(rayTracingAngle)
	mathHelper.On("castRay", player.GetPosition(), worlMap, rayTracingAngle, visibility).Return(rayTracingDestinationPoint).Return(rayTracingDestinationPoint)
	mathHelper.On("calculateProjectionDistance", playerPosition, rayTracingDestinationPoint, player.GetAngle()-rayTracingAngle).Return(projectedDistance)
	mathHelper.On("GetFillRowRange", projectedDistance, float64(screenHeight)).Return(startRow, endRow)
	mathHelper.On("isWallAngle", rayTracingDestinationPoint).Return(isWallAngle)
	for rowIndex := 0; rowIndex <= startRow; rowIndex++ {
		raySampler.On("GetBackgroundRune", rowIndex).Return(backgroundRune)
		raySampler.On("GetBackgroundStyle", rowIndex).Return(backgroundStyle)
		screen.On("SetContent", columnIndex, rowIndex, backgroundRune, []int32(nil), backgroundStyle)
	}
	for rowIndex := startRow + 1; rowIndex < endRow; rowIndex++ {
		raySampler.On("GetWallRune", rowIndex).Return(wallRune)
		screen.On("SetContent", columnIndex, rowIndex, wallRune, []int32(nil), wallAngleStyle)
	}
	for rowIndex := endRow; rowIndex <= screenHeight; rowIndex++ {
		raySampler.On("GetBackgroundRune", rowIndex).Return(backgroundRune)
		raySampler.On("GetBackgroundStyle", rowIndex).Return(backgroundStyle)
		screen.On("SetContent", columnIndex, rowIndex, backgroundRune, []int32(nil), backgroundStyle)
	}
	backgroundColumnRenderer.render(screen, player, worlMap, columnIndex)
	mathHelper.AssertExpectations(t)
}

func TestBackgroundColumnRendererNilRayTracing(t *testing.T) {
	screenWidth := 5
	screenHeight := 10
	fieldOfViewAngle := 0.5
	visibility := 5.0
	mathHelper := new(MockBackgroundRendererMathHelper)
	wallAngleStyle := tcell.StyleDefault.Background(tcell.Color104)
	raySampler := new(testutils.MockRaySampler)
	backgroundColumnRenderer := CreateBackgroundColumnRenderer(screenWidth, screenHeight, fieldOfViewAngle, visibility, mathHelper, wallAngleStyle, raySampler)
	screen := new(testutils.MockScreen)
	playerPosition := new(common.Point2D)
	playerAngle := 0.5
	worlMap := new(testutils.MockWorldMap)
	player := environment.NewPlayer(playerPosition, playerAngle, 1.0, worlMap)
	columnIndex := 1
	rayTracingAngle := 0.25
	backgroundStyle := tcell.StyleDefault.Background(tcell.Color102)
	backgroundRune := '2'

	mathHelper.On("getRayTracingAngleForColumn", player.GetAngle(), columnIndex, screenWidth, fieldOfViewAngle).Return(rayTracingAngle)
	mathHelper.On("castRay", player.GetPosition(), worlMap, rayTracingAngle, visibility).Return(nil)
	for rowIndex := 0; rowIndex <= screenHeight; rowIndex++ {
		raySampler.On("GetBackgroundRune", rowIndex).Return(backgroundRune)
		raySampler.On("GetBackgroundStyle", rowIndex).Return(backgroundStyle)
		screen.On("SetContent", columnIndex, rowIndex, backgroundRune, []int32(nil), backgroundStyle)
	}
	backgroundColumnRenderer.render(screen, player, worlMap, columnIndex)
	mathHelper.AssertExpectations(t)
}
