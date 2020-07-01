package render

import (
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/projectile"
	"francoisgergaud/3dGame/common/environment/world"

	"github.com/gdamore/tcell"
)

//Renderer provides the functionalities to render the environment's map.
type Renderer interface {
	Render(playerID string, worldMap world.WorldMap, player animatedelement.AnimatedElement, worldElements map[string]animatedelement.AnimatedElement, projectiles map[string]projectile.Projectile, screen tcell.Screen)
}
