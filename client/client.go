package client

import (
	"francoisgergaud/3dGame/client/connector"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/environment/animatedelement/projectile"
	"francoisgergaud/3dGame/common/event"

	"github.com/gdamore/tcell"
)

//Engine represents an game engine which takes input and render on screen.
type Engine interface {
	//publisher.EventListener
	Action(eventKey *tcell.EventKey)
	Player() animatedelement.AnimatedElement
	OtherPlayers() map[string]animatedelement.AnimatedElement
	Projectiles() map[string]projectile.Projectile
	ReceiveEventsFromServer(events []event.Event)
	Shutdown()
	ConnectToServer(connectionToServer connector.ServerConnector)
}
