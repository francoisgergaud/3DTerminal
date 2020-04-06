package event

import "francoisgergaud/3dGame/common/environment/animatedelement/state"

//Event is an event, with its publisher, the type of event, and the publisher's state.
type Event struct {
	PlayerID  string
	Action    string
	State     *state.AnimatedElementState
	TimeFrame uint32
	ExtraData map[string]interface{}
}
