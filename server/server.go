package server

import (
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/event"
	"francoisgergaud/3dGame/server/connector"
)

//Server defines the requirement for the server part of the game: in charge of managing the environment and communicate with the players
// - providing the environment to new player
// - listen to the event from player and update the environment accordingly
// - update the environment inertnally (bots)
// - communicate environment changes to players
type Server interface {
	RegisterPlayer(clientConnection connector.ClientConnection) (playerID string, worldMap world.WorldMap, state state.AnimatedElementState, otherPlayers map[string]state.AnimatedElementState)
	UnregisterClient(playerID string)
	ReceiveEventsFromClient(events []event.Event)
}
