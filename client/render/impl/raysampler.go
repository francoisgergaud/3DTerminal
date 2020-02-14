package impl

import "github.com/gdamore/tcell"

//RaySampler is the Ray-sampler interface
type RaySampler interface {
	GetBackgroundRune(rowIndex int) rune
	GetWallRune(rowIndex int) rune
	GetBackgroundStyle(rowIndex int) tcell.Style
	GetWallStyleFromDistance(distance float64) tcell.Style
}
