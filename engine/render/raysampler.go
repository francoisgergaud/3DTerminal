package render

import (
	"math"

	"github.com/gdamore/tcell"
)

//RaySampler is the Ray-sampler interface
type RaySampler interface {
	GetBackgroundRune(rowIndex int) rune
	GetWallRune(rowIndex int) rune
	GetBackgroundStyle(rowIndex int) tcell.Style
	GetWallStyleFromDistance(distance float64) tcell.Style
}

//GradientRaySampler sample a ray given it depth and some other properties
//TODO: remove outOfRangeStyle property
type GradientRaySampler struct {
	//style to be used to render a wall (i.e: the color used in a range of distance).
	wallStyles []tcell.Style
	//style used to render the floor
	backgroundStyle []tcell.Style
	//the depth ranges: each distance to be sampled will be in a range. The depthRanges set the upper-limit of a range at its index.
	depthRanges []float64
	//the ordered background-range to which apply the background-colors. The number is the upper-limit of a range at its index.
	backgroundRange []float32
	//the ordered background-colors to be applied by background-range.
	backgroundRangeColor []tcell.Style
}

//CreateRaySamplerForAnsiColorTerminal create the Gradients
func CreateRaySamplerForAnsiColorTerminal(first float64, multiplicator float64, maxLimit float64, wallStartColor int, wallEndColor int, screenHeight int, backgroundRange []float32, backgroundColors []int) (g *GradientRaySampler) {
	result := &GradientRaySampler{
		wallStyles:      make([]tcell.Style, 1),
		backgroundRange: backgroundRange,
	}
	currentLimit := first
	for currentLimit < maxLimit {
		result.depthRanges = append(result.depthRanges, currentLimit)
		currentLimit *= multiplicator
	}
	result.depthRanges = append(result.depthRanges, math.Inf(1))

	result.setColorArrayFromDepthRange(&result.wallStyles, wallStartColor, wallEndColor)

	result.backgroundRangeColor = make([]tcell.Style, len(backgroundColors))
	for index, element := range backgroundColors {
		result.backgroundRangeColor[index] = tcell.StyleDefault.Background(tcell.Color(element))
	}

	result.setBackground(screenHeight)

	return result
}

func (raySampler *GradientRaySampler) setColorArrayFromDepthRange(array *[]tcell.Style, startColor, endColor int) {
	if startColor > endColor {
		colorStep := (startColor - endColor) / len(raySampler.depthRanges)
		for i := 0; i < len(raySampler.depthRanges); i++ {
			*array = append(*array, tcell.StyleDefault.Background(tcell.Color(startColor-colorStep*i)))
		}
	} else {
		colorStep := (endColor - startColor) / len(raySampler.depthRanges)
		for i := 0; i < len(raySampler.depthRanges); i++ {
			*array = append(*array, tcell.StyleDefault.Background(tcell.Color(startColor+colorStep*i)))
		}
	}
}

//setBackground reset the background-style and runes based from the input screen height.
func (raySampler *GradientRaySampler) setBackground(screenHeight int) {
	raySampler.backgroundStyle = make([]tcell.Style, screenHeight)
	currentBackgroundRange := 0
	for i := 0; i < screenHeight; i++ {
		for i > int(float32(screenHeight)*raySampler.backgroundRange[currentBackgroundRange]) {
			currentBackgroundRange++
		}
		raySampler.backgroundStyle[i] = raySampler.backgroundRangeColor[currentBackgroundRange]
	}
}

//GetBackgroundRune returns the rune used for the background at a specific row number.
func (raySampler *GradientRaySampler) GetBackgroundRune(rowIndex int) rune {
	return ' '
}

//GetWallRune returns the rune used for the wall at a specific row number.
func (raySampler *GradientRaySampler) GetWallRune(rowIndex int) rune {
	return ' '
}

//GetBackgroundStyle returns the style used for the background at a specific row number.
func (raySampler *GradientRaySampler) GetBackgroundStyle(rowIndex int) tcell.Style {
	return raySampler.backgroundStyle[rowIndex]
}

//GetWallStyleFromDistance returns the wall's style for a given distance.
func (raySampler *GradientRaySampler) GetWallStyleFromDistance(distance float64) tcell.Style {
	rangeNumber := 0
	for rangeNumber < len(raySampler.depthRanges) {
		if distance < raySampler.depthRanges[rangeNumber] {
			break
		}
		rangeNumber++
	}
	return raySampler.wallStyles[rangeNumber]
}
