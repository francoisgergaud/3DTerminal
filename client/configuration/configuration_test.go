package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	worldUpdateRate := 1
	configuration := NewConfiguration(worldUpdateRate)
	assert.Equal(t, worldUpdateRate, configuration.WorlUpdateRate)
	assert.Greater(t, configuration.FrameRate, 1)
	assert.Greater(t, len(configuration.GradientRSBackgroundColors), 1)
	assert.Greater(t, len(configuration.GradientRSBackgroundRange), 1)
	assert.Greater(t, configuration.GradientRSFirst, 0.0)
	assert.Greater(t, configuration.GradientRSLimit, 0.0)
	assert.Greater(t, configuration.GradientRSMultiplicator, 0.0)
	assert.Greater(t, configuration.GradientRSWallEndColor, 0)
	assert.Greater(t, configuration.GradientRSWallStartColor, 0)
	assert.Greater(t, configuration.PlayerFieldOfViewAngle, 0.1)
	assert.Greater(t, configuration.ScreenHeight, 0)
	assert.Greater(t, configuration.ScreenWidth, 0)
	assert.Greater(t, configuration.Visibility, 1.0)
}
