package event

import (
	"encoding/json"
	"francoisgergaud/3dGame/common/environment/animatedelement/state"
	"francoisgergaud/3dGame/common/environment/world"
	"francoisgergaud/3dGame/common/math"
	testworld "francoisgergaud/3dGame/internal/testutils/common/environment/world"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUnmarshalMessage(t *testing.T) {
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
			"otherPlayers": map[string]*state.AnimatedElementState{
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
			"projectiles": map[string]*state.AnimatedElementState{
				"projectTest1": {
					Position: &math.Point2D{
						X: 10.0,
						Y: 2.5,
					},
					MoveDirection: state.Backward,
					Size:          0.25,
					StepAngle:     0.75,
					Velocity:      2.25,
					Angle:         0.05,
					Style:         tcell.StyleDefault.Foreground(tcell.Color110),
				},
			},
			"projectileID": "projectileIDTest",
			"playerID":     "playerIDTest",
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
	assert.Equal(t, eventToMarshal.ExtraData["otherPlayers"].(map[string]*state.AnimatedElementState)["otherPlayer"], eventToUnmarshal.ExtraData["otherPlayers"].(map[string]*state.AnimatedElementState)["otherPlayer"])
	assert.Equal(t, eventToMarshal.ExtraData["projectiles"].(map[string]*state.AnimatedElementState)["projectTest1"], eventToUnmarshal.ExtraData["projectiles"].(map[string]*state.AnimatedElementState)["projectTest1"])
	assert.Equal(t, eventToMarshal.ExtraData["projectileID"].(string), eventToUnmarshal.ExtraData["projectileID"].(string))
	assert.Equal(t, eventToMarshal.ExtraData["playerID"].(string), eventToUnmarshal.ExtraData["playerID"].(string))
}

func TestUnmarshalMessageWrongExtraData(t *testing.T) {
	eventToMarshal := Event{
		ExtraData: map[string]interface{}{
			"wrongKey": "playerIDTest",
		},
	}
	bytes, err := json.Marshal(eventToMarshal)
	assert.Nil(t, err)
	eventToUnmarshal := Event{}
	assert.Error(t, json.Unmarshal(bytes, &eventToUnmarshal))
}

func TestClone(t *testing.T) {
	worldMap := new(testworld.MockWorldMap)
	eventToClone := Event{
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
			"otherPlayers": map[string]*state.AnimatedElementState{
				"otherPlayerID": {
					Angle: 0.25,
				},
			},
			"projectiles": map[string]*state.AnimatedElementState{
				"projectileTest1": {
					Angle: 0.075,
				},
			},
			"worldMap":     worldMap,
			"playerID":     "playerIDTest",
			"projectileID": "projectileIDTest",
		},
	}

	worldMap.On("Clone").Return(new(testworld.MockWorldMap))
	result, err := eventToClone.Clone()
	assert.NotEqual(t, &eventToClone, &result)
	assert.False(t, eventToClone.State == result.State)
	assert.Equal(t, *eventToClone.State, *result.State)
	eventToCloneOtherPlayer := eventToClone.ExtraData["otherPlayers"].(map[string]*state.AnimatedElementState)["otherPlayerID"]
	resultOtherPlayer := result.ExtraData["otherPlayers"].(map[string]*state.AnimatedElementState)["otherPlayerID"]
	eventToCloneProjectile := eventToClone.ExtraData["projectiles"].(map[string]*state.AnimatedElementState)["projectileTest1"]
	resultProjectile := result.ExtraData["projectiles"].(map[string]*state.AnimatedElementState)["projectileTest1"]
	assert.Equal(t, eventToCloneOtherPlayer, resultOtherPlayer)
	assert.Equal(t, eventToCloneProjectile, resultProjectile)
	assert.Nil(t, err)
	mock.AssertExpectationsForObjects(t, worldMap)
}

func TestCloneWrongExtraData(t *testing.T) {
	eventToClone := Event{
		ExtraData: map[string]interface{}{
			"wrongKey": "wrongValue",
		},
	}

	result, err := eventToClone.Clone()
	assert.Nil(t, result)
	assert.Error(t, err)

}
