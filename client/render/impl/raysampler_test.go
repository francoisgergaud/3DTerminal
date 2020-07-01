package impl

import (
	"testing"

	"github.com/gdamore/tcell"

	"github.com/stretchr/testify/assert"
)

func TestCreateRaySamplerForAnsiColorTerminal(t *testing.T) {
	backgroundRanges := []float32{0.1, 0.3, 0.5}
	backgroundColors := []int{10, 11, 12, 13}
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(1.0, 2.0, 5.0, 0, 5, 10, backgroundRanges, backgroundColors)
	assert.Nil(t, err)
	rowIndex := 1
	assert.Equal(t, gradientRaySampler.GetBackgroundRune(rowIndex), ' ')
	assert.Equal(t, gradientRaySampler.GetWallRune(rowIndex), ' ')
	assert.Equal(t, tcell.StyleDefault.Background(10), gradientRaySampler.GetBackgroundStyle(0))
	assert.Equal(t, tcell.StyleDefault.Background(11), gradientRaySampler.GetBackgroundStyle(2))
	assert.Equal(t, tcell.StyleDefault.Background(11), gradientRaySampler.GetBackgroundStyle(3))
	assert.Equal(t, tcell.StyleDefault.Background(12), gradientRaySampler.GetBackgroundStyle(5))
	assert.Equal(t, tcell.StyleDefault.Background(12), gradientRaySampler.GetBackgroundStyle(6))
	assert.Equal(t, tcell.StyleDefault.Background(0), gradientRaySampler.GetWallStyleFromDistance(0))
	assert.Equal(t, tcell.StyleDefault.Background(1), gradientRaySampler.GetWallStyleFromDistance(1))
	assert.Equal(t, tcell.StyleDefault.Background(2), gradientRaySampler.GetWallStyleFromDistance(2))
	assert.Equal(t, tcell.StyleDefault.Background(2), gradientRaySampler.GetWallStyleFromDistance(3))
	assert.Equal(t, tcell.StyleDefault.Background(0), gradientRaySampler.GetWallStyleFromDistance(4))
	assert.Equal(t, tcell.StyleDefault.Background(0), gradientRaySampler.GetWallStyleFromDistance(5))
	assert.Equal(t, tcell.StyleDefault.Background(0), gradientRaySampler.GetWallStyleFromDistance(6))
}

func TestCreateRaySamplerForAnsiColorTerminalWithInvertedColors(t *testing.T) {
	backgroundRanges := []float32{0.1, 0.3, 0.5}
	backgroundColors := []int{10, 11, 12, 13}
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(1.0, 2.0, 5.0, 5, 0, 10, backgroundRanges, backgroundColors)
	assert.Nil(t, err)
	rowIndex := 1
	assert.Equal(t, gradientRaySampler.GetBackgroundRune(rowIndex), ' ')
	assert.Equal(t, gradientRaySampler.GetWallRune(rowIndex), ' ')
	assert.Equal(t, tcell.StyleDefault.Background(10), gradientRaySampler.GetBackgroundStyle(0))
	assert.Equal(t, tcell.StyleDefault.Background(11), gradientRaySampler.GetBackgroundStyle(2))
	assert.Equal(t, tcell.StyleDefault.Background(11), gradientRaySampler.GetBackgroundStyle(3))
	assert.Equal(t, tcell.StyleDefault.Background(12), gradientRaySampler.GetBackgroundStyle(5))
	assert.Equal(t, tcell.StyleDefault.Background(12), gradientRaySampler.GetBackgroundStyle(6))
	assert.Equal(t, tcell.StyleDefault.Background(0), gradientRaySampler.GetWallStyleFromDistance(0))
	assert.Equal(t, tcell.StyleDefault.Background(1), gradientRaySampler.GetWallStyleFromDistance(1))
	assert.Equal(t, tcell.StyleDefault.Background(2), gradientRaySampler.GetWallStyleFromDistance(2))
	assert.Equal(t, tcell.StyleDefault.Background(2), gradientRaySampler.GetWallStyleFromDistance(3))
	assert.Equal(t, tcell.StyleDefault.Background(0), gradientRaySampler.GetWallStyleFromDistance(4))
	assert.Equal(t, tcell.StyleDefault.Background(0), gradientRaySampler.GetWallStyleFromDistance(5))
	assert.Equal(t, tcell.StyleDefault.Background(0), gradientRaySampler.GetWallStyleFromDistance(6))
}

func TestCreateRaySamplerForAnsiColorTerminalWithInvalidFirst(t *testing.T) {
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(-1.0, 2.0, 5.0, 0, 5, 10, []float32{0.1, 0.3, 0.5}, []int{10, 11, 12, 13})
	assert.Nil(t, gradientRaySampler)
	assert.Error(t, err)
}

func TestCreateRaySamplerForAnsiColorTerminalWithInvalidMultiplicator(t *testing.T) {
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(1.0, -2.0, 5.0, 0, 5, 10, []float32{0.1, 0.3, 0.5}, []int{10, 11, 12, 13})
	assert.Nil(t, gradientRaySampler)
	assert.Error(t, err)
}

func TestCreateRaySamplerForAnsiColorTerminalWithInvalidMaxLimit(t *testing.T) {
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(1.0, 2.0, -5.0, 0, 5, 10, []float32{0.1, 0.3, 0.5}, []int{10, 11, 12, 13})
	assert.Nil(t, gradientRaySampler)
	assert.Error(t, err)
}

func TestCreateRaySamplerForAnsiColorTerminalWithInvalidStartColor(t *testing.T) {
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(1.0, 2.0, 5.0, -1, 5, 10, []float32{0.1, 0.3, 0.5}, []int{10, 11, 12, 13})
	assert.Nil(t, gradientRaySampler)
	assert.Error(t, err)
}

func TestCreateRaySamplerForAnsiColorTerminalWithInvalidEndColor(t *testing.T) {
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(1.0, 2.0, 5.0, 1, -5, 10, []float32{0.1, 0.3, 0.5}, []int{10, 11, 12, 13})
	assert.Nil(t, gradientRaySampler)
	assert.Error(t, err)
}

func TestCreateRaySamplerForAnsiColorTerminalWithInvalidScreenHeight(t *testing.T) {
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(1.0, 2.0, 5.0, 1, 5, -10, []float32{0.1, 0.3, 0.5}, []int{10, 11, 12, 13})
	assert.Nil(t, gradientRaySampler)
	assert.Error(t, err)
}

func TestCreateRaySamplerForAnsiColorTerminalWithEmptyBackgroundRange(t *testing.T) {
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(1.0, 2.0, 5.0, 1, 5, 10, []float32{}, []int{10, 11, 12, 13})
	assert.Nil(t, gradientRaySampler)
	assert.Error(t, err)
}

func TestCreateRaySamplerForAnsiColorTerminalWithInvalidBackgroundRange(t *testing.T) {
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(1.0, 2.0, 5.0, 1, 5, 10, []float32{0.3, 0.1, 0.5}, []int{10, 11, 12, 13})
	assert.Nil(t, gradientRaySampler)
	assert.Error(t, err)
}

func TestCreateRaySamplerForAnsiColorTerminalWithInvalidBackgroundColors(t *testing.T) {
	gradientRaySampler, err := CreateRaySamplerForAnsiColorTerminal(1.0, 2.0, 5.0, 1, 5, 10, []float32{0.1, 0.3, 0.5}, []int{10})
	assert.Nil(t, gradientRaySampler)
	assert.Error(t, err)
}
