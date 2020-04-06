package configuration

//NewConfiguration is the default engine-configuration factory
func NewConfiguration(worldUpdateRate int) *Configuration {
	return &Configuration{
		FrameRate:                  20,
		WorlUpdateRate:             worldUpdateRate,
		ScreenHeight:               40,
		ScreenWidth:                120,
		PlayerFieldOfViewAngle:     0.4,
		Visibility:                 20.0,
		GradientRSFirst:            1.0,
		GradientRSMultiplicator:    2.0,
		GradientRSLimit:            10.0,
		GradientRSWallStartColor:   255,
		GradientRSWallEndColor:     240,
		GradientRSBackgroundRange:  []float32{0.5, 0.55, 0.65},
		GradientRSBackgroundColors: []int{63, 58, 64, 70},
	}
}

//Configuration contains the required parametrable parameters for the engine.
type Configuration struct {
	//The frame-rate per second.
	FrameRate int
	//the world-update's rate.
	WorlUpdateRate int
	//The screen's height.
	ScreenHeight int
	//The screen's width.
	ScreenWidth int
	//The player's (or camera) field-of-view angle in Pie radian.
	PlayerFieldOfViewAngle float64
	//The player's (or camera) maximum's visibility.
	Visibility float64
	//The gradient-ray-sampler first distance upper-range value.
	GradientRSFirst float64
	//The gradient-ray-sampler distance exponential multiplicator.
	GradientRSMultiplicator float64
	//The gradient-ray-sampler distance maximum upper-range. After this range, the last gradient-color will be used until infinit.
	GradientRSLimit float64
	//The gradient-ray-sampler start-color value (closer color).
	GradientRSWallStartColor int
	//The gradient-ray-sampler end-color value (farest color).
	GradientRSWallEndColor int
	//The gradient-ray-sampler background-column-index ratio, from 0 to 1. The last value must be 1.0, and the values must be increasing
	GradientRSBackgroundRange []float32
	//The gradient-ray-sampler background-colors, which apply to the upper-range ratio of the row defined in GradientRSBackgroundRange.
	GradientRSBackgroundColors []int
}
