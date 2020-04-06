package client

import (
	"francoisgergaud/3dGame/client/connector"
	"francoisgergaud/3dGame/client/player"
	"francoisgergaud/3dGame/common/environment/animatedelement"
	"francoisgergaud/3dGame/common/event"
)

//Engine represents an game engine which takes input and render on screen.
type Engine interface {
	Player() player.Player
	OtherPlayers() map[string]animatedelement.AnimatedElement
	ReceiveEventsFromServer(events []event.Event)
	Shutdown()
	ConnectToServer(connectionToServer connector.ServerConnector)
}
