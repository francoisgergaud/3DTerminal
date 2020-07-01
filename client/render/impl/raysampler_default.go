package impl

import (
	"fmt"

	"github.com/gdamore/tcell"
)

//GradientRaySampler sample a ray given it depth and some other properties
//TODO: remove outOfRangeStyle property
type GradientRaySampler struct {
	//styles to be used to render a wall (i.e: the color used in a range of distance).
	wallStyles []tcell.Style
	//styles used to render the floor
	backgroundStyles []tcell.Style
	//the depth-ranges: each distance to be sampled will be in a range. The depthRanges set the upper-limit of a range at its index.
	depthRanges []float64
	//the ordered-background-ranges to which apply the background-colors. The number is the upper-limit of a range at its index.
	backgroundRanges []float32
	//the ordered background-colors to be applied by background-ranges.
	backgroundRangesColors []tcell.Style
}

//CreateRaySamplerForAnsiColorTerminal create the Gradients
func CreateRaySamplerForAnsiColorTerminal(first float64, multiplicator float64, maxLimit float64, wallStartColor int, wallEndColor int, screenHeight int, backgroundRange []float32, backgroundColors []int) (g RaySampler, err error) {
	if first < 0.0 {
		return nil, fmt.Errorf("Gradient ray-sampler 'first' value cannot be negative")
	}
	if multiplicator <= 0.0 {
		return nil, fmt.Errorf("Gradient-ray-sampler 'multiplicator' value cannot be negative or 0")
	}
	if maxLimit <= 0.0 {
		return nil, fmt.Errorf("Gradient ray-sampler 'maxLimit' value cannot be negative or 0")
	}
	if wallStartColor < 0.0 {
		return nil, fmt.Errorf("Gradient ray-sampler 'wallStartColor' value cannot be negative")
	}
	if wallEndColor < 0.0 {
		return nil, fmt.Errorf("Gradient ray-sampler 'wallEndColor' value cannot be negative")
	}
	if screenHeight < 0 {
		return nil, fmt.Errorf("Gradient ray-sampler 'screenHeight' value cannot be negative")
	}
	if len(backgroundRange) == 0 {
		return nil, fmt.Errorf("Gradient ray-sampler 'backgroundRange' array cannot be empty")
	}
	previousBackgroundRangeValue := backgroundRange[0]
	for i := 1; i < len(backgroundRange); i++ {
		if previousBackgroundRangeValue > backgroundRange[i] {
			return nil, fmt.Errorf("Gradient ray-sampler 'backgroundRange' must be ordered from smallest to biggest value")
		}
		previousBackgroundRangeValue = backgroundRange[i]
	}
	if len(backgroundRange)+1 != len(backgroundColors) {
		return nil, fmt.Errorf("Gradient ray-sampler 'backgroundColors' length must be 'backgroundRange' length + 1")
	}
	result := &GradientRaySampler{
		backgroundRanges: backgroundRange,
	}
	currentLimit := first
	for currentLimit < maxLimit {
		result.depthRanges = append(result.depthRanges, currentLimit)
		currentLimit *= multiplicator
	}
	//result.depthRanges = append(result.depthRanges, math.Inf(1))
	result.wallStyles = result.getColorArrayFromDepthRange(wallStartColor, wallEndColor)
	result.backgroundRangesColors = make([]tcell.Style, len(backgroundColors))
	for index, element := range backgroundColors {
		result.backgroundRangesColors[index] = tcell.StyleDefault.Background(tcell.Color(element))
	}
	result.setBackgroundStyles(screenHeight)
	return result, nil
}

func (raySampler *GradientRaySampler) getColorArrayFromDepthRange(startColor, endColor int) []tcell.Style {
	styles := make([]tcell.Style, len(raySampler.depthRanges)+1)
	if startColor > endColor {
		colorStep := (startColor - endColor) / len(raySampler.depthRanges)
		for i := 0; i < len(raySampler.depthRanges); i++ {
			styles[i] = tcell.StyleDefault.Background(tcell.Color(endColor + colorStep*i))
		}
		styles[len(raySampler.depthRanges)] = tcell.StyleDefault.Background(tcell.Color(endColor))
	} else {
		colorStep := (endColor - startColor) / len(raySampler.depthRanges)
		for i := 0; i < len(raySampler.depthRanges); i++ {
			styles[i] = tcell.StyleDefault.Background(tcell.Color(startColor + colorStep*i))
		}
		styles[len(raySampler.depthRanges)] = tcell.StyleDefault.Background(tcell.Color(startColor))
	}
	return styles
}

//setBackgroundStyles reset the background-style and runes based from the new screen height.
func (raySampler *GradientRaySampler) setBackgroundStyles(screenHeight int) {
	raySampler.backgroundStyles = make([]tcell.Style, screenHeight)
	currentBackgroundRange := 0
	for i := 0; i < screenHeight; i++ {
		if i > int(float32(screenHeight)*raySampler.backgroundRanges[currentBackgroundRange]) && currentBackgroundRange < len(raySampler.backgroundRanges)-1 {
			currentBackgroundRange++
		}
		raySampler.backgroundStyles[i] = raySampler.backgroundRangesColors[currentBackgroundRange]
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
	return raySampler.backgroundStyles[rowIndex]
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
