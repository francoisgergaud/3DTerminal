package render

import (
	"testing"

	"github.com/gdamore/tcell"

	"github.com/stretchr/testify/assert"
)

func TestRaySamplerGetBackgroundRune(t *testing.T) {
	first := 1.0
	multiplicator := 2.0
	maxLimit := 5.0
	wallStartColor := 0
	wallEndColor := 5
	screenHeight := 10
	backgroundRange := []float32{1.0, 3.0, 5.0}
	backgroundColors := []int{10, 11, 12, 13}
	gradientRaySampler := CreateRaySamplerForAnsiColorTerminal(first, multiplicator, maxLimit, wallStartColor, wallEndColor, screenHeight, backgroundRange, backgroundColors)
	rowIndex := 1
	result := gradientRaySampler.GetBackgroundRune(rowIndex)
	assert.Equal(t, result, ' ')
}

func TestRaySamplerGetWallRune(t *testing.T) {
	first := 1.0
	multiplicator := 2.0
	maxLimit := 5.0
	wallStartColor := 0
	wallEndColor := 5
	screenHeight := 10
	backgroundRange := []float32{1.0, 3.0, 5.0}
	backgroundColors := []int{10, 11, 12, 13}
	gradientRaySampler := CreateRaySamplerForAnsiColorTerminal(first, multiplicator, maxLimit, wallStartColor, wallEndColor, screenHeight, backgroundRange, backgroundColors)
	rowIndex := 1
	result := gradientRaySampler.GetWallRune(rowIndex)
	assert.Equal(t, result, ' ')
}

func TestRaySamplerGetBackgroundStyle(t *testing.T) {
	first := 1.0
	multiplicator := 2.0
	maxLimit := 5.0
	wallStartColor := 0
	wallEndColor := 5
	screenHeight := 10
	backgroundRange := []float32{0.2, 0.5, 1.0}
	backgroundColors := []int{10, 11, 12}
	gradientRaySampler := CreateRaySamplerForAnsiColorTerminal(first, multiplicator, maxLimit, wallStartColor, wallEndColor, screenHeight, backgroundRange, backgroundColors)
	assert.Equal(t, tcell.StyleDefault.Background(10), gradientRaySampler.GetBackgroundStyle(0))
	assert.Equal(t, tcell.StyleDefault.Background(10), gradientRaySampler.GetBackgroundStyle(2))
	assert.Equal(t, tcell.StyleDefault.Background(11), gradientRaySampler.GetBackgroundStyle(3))
	assert.Equal(t, tcell.StyleDefault.Background(11), gradientRaySampler.GetBackgroundStyle(5))
	assert.Equal(t, tcell.StyleDefault.Background(12), gradientRaySampler.GetBackgroundStyle(6))
}
