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
		case "otherPlayers":
			s := make(map[string]state.AnimatedElementState)
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
		default:
			return errors.New("extra-data: " + key + "is not managed for JSON deserialization")
		}
	}
	*extradData = newExtradData
	return nil
}
