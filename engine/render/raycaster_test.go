package render

import (
	"francoisgergaud/3dGame/common"
	"francoisgergaud/3dGame/internal/testutils"
	"testing"
)

// map is not mocked. It is easier to visualizer the orgin and destination of the ray with a visualization of the map.
var grid = [][]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

const WALL = 1
const NOWALL = 0

var world testutils.MockWorldMap = testutils.MockWorldMap{}

func TestRayCast1(t *testing.T) {
	raycaster := new(RayCasterImpl)
	startPoint := common.Point2D{X: 5, Y: 5}
	angle := 0.0
	visibility := 10.0
	for offset := 0; offset < 4; offset++ {
		world.On("GetCellValue", 5+offset, 5).Return(NOWALL)
	}
	world.On("GetCellValue", 9, 5).Return(WALL)
	impact := raycaster.rayCast(&startPoint, &world, angle, visibility)
	expectedImpact := common.Point2D{X: 9, Y: 5}
	if impact == nil {
		t.Errorf("RayCast is incorrect, got nil.")
	} else if !impact.AlmostEquals(expectedImpact) {
		t.Errorf("RayCast was incorrect, expected %s, got: %s.", expectedImpact, impact)
	}
}

func TestRayCast2(t *testing.T) {
	raycaster := new(RayCasterImpl)
	startPoint := common.Point2D{X: 5, Y: 5}
	angle := 0.25
	visibility := 10.0
	for offset := 0; offset < 4; offset++ {
		world.On("GetCellValue", 5+offset, 5+offset).Return(NOWALL)
	}
	world.On("GetCellValue", 9, 9).Return(WALL)
	impact := raycaster.rayCast(&startPoint, &world, angle, visibility)
	expectedImpact := common.Point2D{X: 9, Y: 9}
	if impact == nil {
		t.Errorf("RayCast is incorrect, got nil.")
	} else if !impact.AlmostEquals(expectedImpact) {
		t.Errorf("RayCast was incorrect, expected %s, got: %s.", expectedImpact, impact)
	}
}

func TestRayCast3(t *testing.T) {
	raycaster := new(RayCasterImpl)
	startPoint := common.Point2D{X: 5, Y: 5}
	angle := 0.5
	visibility := 10.0
	world.On("GetCellValue", 4, 5).Return(NOWALL)
	for offset := 0; offset < 4; offset++ {
		world.On("GetCellValue", 5, 5+offset).Return(NOWALL)
	}
	world.On("GetCellValue", 5, 9).Return(WALL)
	impact := raycaster.rayCast(&startPoint, &world, angle, visibility)
	expectedImpact := common.Point2D{X: 5, Y: 9}
	if impact == nil {
		t.Errorf("RayCast is incorrect, got nil.")
	} else if !impact.AlmostEquals(expectedImpact) {
		t.Errorf("RayCast was incorrect, expected %s, got: %s.", expectedImpact, impact)
	}
}

func TestRayCast4(t *testing.T) {
	raycaster := new(RayCasterImpl)
	startPoint := common.Point2D{X: 5, Y: 5}
	angle := 0.75
	visibility := 10.0
	for offset := 0; offset < 4; offset++ {
		world.On("GetCellValue", 5-offset, 5+offset).Return(NOWALL)
		world.On("GetCellValue", 4-offset, 5+offset).Return(NOWALL)
	}
	world.On("GetCellValue", 0, 9).Return(WALL)
	impact := raycaster.rayCast(&startPoint, &world, angle, visibility)
	expectedImpact := common.Point2D{X: 1, Y: 9}
	if impact == nil {
		t.Errorf("RayCast is incorrect, got nil.")
	} else if !impact.AlmostEquals(expectedImpact) {
		t.Errorf("RayCast was incorrect, expected %s, got: %s.", expectedImpact, impact)
	}
}

func TestRayCast5(t *testing.T) {
	raycaster := new(RayCasterImpl)
	startPoint := common.Point2D{X: 5, Y: 5}
	angle := 1.0
	visibility := 10.0
	for offset := 0; offset < 5; offset++ {
		world.On("GetCellValue", 5-offset, 5).Return(NOWALL)
		//edge case: both the cells above and bellow the x axis are checked
		world.On("GetCellValue", 5-offset, 4).Return(NOWALL)
	}
	world.On("GetCellValue", 0, 5).Return(WALL)
	impact := raycaster.rayCast(&startPoint, &world, angle, visibility)
	expectedImpact := common.Point2D{X: 1, Y: 5}
	if impact == nil {
		t.Errorf("RayCast is incorrect, got nil.")
	} else if !impact.AlmostEquals(expectedImpact) {
		t.Errorf("RayCast was incorrect, expected %s, got: %s.", expectedImpact, impact)
	}
}

func TestRayCast6(t *testing.T) {
	raycaster := new(RayCasterImpl)
	startPoint := common.Point2D{X: 5, Y: 5}
	angle := 1.25
	visibility := 10.0
	for offset := 0; offset < 4; offset++ {
		//edge case, the 3 cells on the diagonal axis are checked.
		world.On("GetCellValue", 5-offset, 5-offset).Return(NOWALL)
		world.On("GetCellValue", 5-offset, 4-offset).Return(NOWALL)
		world.On("GetCellValue", 4-offset, 5-offset).Return(NOWALL)
	}
	world.On("GetCellValue", 0, 1).Return(WALL)
	impact := raycaster.rayCast(&startPoint, &world, angle, visibility)
	expectedImpact := common.Point2D{X: 1, Y: 1}
	if impact == nil {
		t.Errorf("RayCast is incorrect, got nil.")
	} else if !impact.AlmostEquals(expectedImpact) {
		t.Errorf("RayCast was incorrect, expected %s, got: %s.", expectedImpact, impact)
	}
}

func TestRayCast7(t *testing.T) {
	raycaster := new(RayCasterImpl)
	startPoint := common.Point2D{X: 5, Y: 5}
	angle := 1.5
	visibility := 10.0
	for offset := 0; offset < 5; offset++ {
		world.On("GetCellValue", 5, 5-offset).Return(NOWALL)
	}
	world.On("GetCellValue", 5, 0).Return(WALL)
	impact := raycaster.rayCast(&startPoint, &world, angle, visibility)
	expectedImpact := common.Point2D{X: 5, Y: 1}
	if impact == nil {
		t.Errorf("RayCast is incorrect, got nil.")
	} else if !impact.AlmostEquals(expectedImpact) {
		t.Errorf("RayCast was incorrect, expected %s, got: %s.", expectedImpact, impact)
	}
}

func TestRayCast8(t *testing.T) {
	raycaster := new(RayCasterImpl)
	startPoint := common.Point2D{X: 5, Y: 5}
	angle := 1.75
	visibility := 10.0
	for offset := 0; offset < 4; offset++ {
		world.On("GetCellValue", 5+offset, 5-offset).Return(NOWALL)
		world.On("GetCellValue", 5+offset, 4-offset).Return(NOWALL)
	}
	world.On("GetCellValue", 9, 0).Return(WALL)
	impact := raycaster.rayCast(&startPoint, &world, angle, visibility)
	expectedImpact := common.Point2D{X: 9, Y: 1}
	if impact == nil {
		t.Errorf("RayCast is incorrect, got nil.")
	} else if !impact.AlmostEquals(expectedImpact) {
		t.Errorf("RayCast was incorrect, expected %s, got: %s.", expectedImpact, impact)
	}
}

func TestRayCast9(t *testing.T) {
	raycaster := new(RayCasterImpl)
	startPoint := common.Point2D{X: 6.3, Y: 7}
	angle := 1.40
	visibility := 10.0
	world.On("GetCellValue", 6, 6).Return(NOWALL)
	world.On("GetCellValue", 5, 6).Return(NOWALL)
	world.On("GetCellValue", 5, 5).Return(NOWALL)
	world.On("GetCellValue", 5, 4).Return(NOWALL)
	world.On("GetCellValue", 5, 3).Return(NOWALL)
	world.On("GetCellValue", 5, 2).Return(NOWALL)
	world.On("GetCellValue", 4, 2).Return(NOWALL)
	world.On("GetCellValue", 4, 1).Return(NOWALL)
	world.On("GetCellValue", 4, 0).Return(WALL)
	impact := raycaster.rayCast(&startPoint, &world, angle, visibility)
	expectedImpact := common.Point2D{X: 4.350, Y: 1}
	if impact == nil {
		t.Errorf("RayCast is incorrect, got nil.")
	} else if !impact.AlmostEquals(expectedImpact) {
		t.Errorf("RayCast was incorrect, expected %s, got: %s.", expectedImpact, impact)
	}
}
