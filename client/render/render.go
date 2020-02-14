package render

import (
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/world"

	"github.com/gdamore/tcell"
)

//Renderer provides the functionalities to render the environment's map.
type Renderer interface {
	Render(worldMap world.WorldMap, player player.Player, worldElements map[string]animatedelement.AnimatedElement, screen tcell.Screen)
}
