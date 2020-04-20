package event

import (
	"encoding/json"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/math"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJoinMessage(t *testing.T) {
	eventToMarshal := Event{
		Action:   "action",
		PlayerID: "playerID",
		State: &state.AnimatedElementState{
			Position: &math.Point2D{
				X: 1.0,
				Y: 3.0,
			},
			Angle:           1.5,
			MoveDirection:   state.Forward,
			RotateDirection: state.Right,
			Size:            0.25,
			StepAngle:       0.025,
			Style:           tcell.StyleDefault.Foreground(tcell.Color106),
			Velocity:        0.75,
		},
		TimeFrame: uint32(98),
		ExtraData: map[string]interface{}{
			"worldMap": world.NewWorldMap(
				[][]int{
					{0, 1},
					{1, 0}}),
			"otherPlayers": map[string]state.AnimatedElementState{
				"otherPlayer": {
					Position: &math.Point2D{
						X: 20.0,
						Y: 12.5,
					},
					MoveDirection: state.Backward,
					Size:          0.75,
					StepAngle:     0.0025,
					Velocity:      2.5,
					Angle:         0.005,
					Style:         tcell.StyleDefault.Foreground(tcell.Color108),
				},
			},
		},
	}
	bytes, err := json.Marshal(eventToMarshal)
	assert.Nil(t, err)
	eventToUnmarshal := Event{}
	json.Unmarshal(bytes, &eventToUnmarshal)
	assert.True(t, true)
	assert.Equal(t, eventToMarshal.Action, eventToUnmarshal.Action)
	assert.Equal(t, eventToMarshal.PlayerID, eventToUnmarshal.PlayerID)
	assert.Equal(t, eventToMarshal.State.Position.X, eventToUnmarshal.State.Position.X)
	assert.Equal(t, eventToMarshal.State.MoveDirection, eventToUnmarshal.State.MoveDirection)
	assert.Equal(t, eventToMarshal.State.Style, eventToUnmarshal.State.Style)
	assert.Equal(t, eventToMarshal.TimeFrame, eventToUnmarshal.TimeFrame)
	assert.Equal(t, eventToMarshal.ExtraData["worldMap"].(world.WorldMap).GetCellValue(0, 0), eventToUnmarshal.ExtraData["worldMap"].(world.WorldMap).GetCellValue(0, 0))
	assert.Equal(t, eventToMarshal.ExtraData["otherPlayers"].(map[string]state.AnimatedElementState)["otherPlayer"].MoveDirection, eventToUnmarshal.ExtraData["otherPlayers"].(map[string]state.AnimatedElementState)["otherPlayer"].MoveDirection)
}
