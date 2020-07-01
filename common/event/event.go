package event

import (
	"encoding/json"
	"errors"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
)

type ExtradData map[string]interface{}

//Event is an event, with its publisher, the type of event, and the publisher's state.
type Event struct {
	PlayerID  string
	Action    string
	State     *state.AnimatedElementState
	TimeFrame uint32
	//data for init-event
	ExtraData ExtradData `json:"pointer,omitempty"`
}

//UnmarshalJSON applies custom deserialization on ExtraData
func (extradData *ExtradData) UnmarshalJSON(data []byte) error {
	jsonRawValueMap := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &jsonRawValueMap)
	if err != nil {
		return err
	}
	newExtradData := make(ExtradData)
	for key, jsonRawValue := range jsonRawValueMap {
		switch key {
		case "otherPlayers", "projectiles":
			s := make(map[string]*state.AnimatedElementState)
			err := json.Unmarshal(jsonRawValue, &s)
			if err != nil {
				return err
			}
			newExtradData[key] = s
		case "worldMap":
			c := world.WorldMapImpl{}
			err := json.Unmarshal(jsonRawValue, &c)
			if err != nil {
				return err
			}
			newExtradData[key] = &c
		case "playerID", "projectileID":
			stringValue := new(string)
			json.Unmarshal(jsonRawValue, stringValue)
			newExtradData[key] = *stringValue
		default:
			return errors.New("extra-data: " + key + " is not managed for JSON deserialization")
		}
	}
	*extradData = newExtradData
	return nil
}

//Clone create an event's deep-copy
func (event Event) Clone() (*Event, error) {
	var stateClone *state.AnimatedElementState
	if event.State != nil {
		stateClone = event.State.Clone()
	}
	var newExtradData ExtradData
	if event.ExtraData != nil {
		newExtradData = make(ExtradData)
		for key, value := range event.ExtraData {
			switch key {
			case "otherPlayers", "projectiles":
				animatedElementStates := make(map[string]*state.AnimatedElementState)
				for animatedElementID, animatedElementState := range value.(map[string]*state.AnimatedElementState) {
					animatedElementClone := animatedElementState.Clone()
					animatedElementStates[animatedElementID] = animatedElementClone
				}
				newExtradData[key] = animatedElementStates
			case "worldMap":
				newExtradData[key] = value.(world.WorldMap).Clone()
			case "playerID", "projectileID":
				newExtradData[key] = value
			default:
				return nil, errors.New("extra-data: " + key + " is not managed for Cloning")
			}
		}
	}
	result := &Event{
		Action:    event.Action,
		PlayerID:  event.PlayerID,
		State:     stateClone,
		TimeFrame: event.TimeFrame,
		ExtraData: newExtradData,
	}
	return result, nil
}
