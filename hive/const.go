package hive

const (
	typeMotionSensor   = "motionsensor"
	typeColourLight    = "colourtuneablelight"
	typeWarmWhiteLight = "warmwhitelight"

	colorCold = 6535
	colorWarm = 2700
)

var (
	statusON  = "ON"
	statusOFF = "OFF"

	colourModeWHITE  = "WHITE"
	colourModeCOLOUR = "COLOUR"

	// Some color constants in the HSV space.
	ColorWhite  = HSV{0, 0, 100}
	ColorGrey   = HSV{0, 0, 50}
	ColorBlack  = HSV{0, 0, 0}
	ColorRed    = HSV{0, 99, 100}
	ColorGreen  = HSV{120, 99, 50}
	ColorBlue   = HSV{250, 89, 100}
	ColorCyan   = HSV{180, 50, 100}
	ColorNavy   = HSV{250, 75, 50}
	ColorYellow = HSV{60, 50, 100}
	ColorOrange = HSV{14, 66, 100}
	ColorPurple = HSV{284, 75, 100}
)

// HSV contains the Hue, Saturation and Value used when setting the color of a
// colored light bulb.
type HSV struct {
	// Hue of the color. Valid values are [0, 359].
	Hue int
	// Saturation of the color. Valid values are [0, 99].
	Saturation int
	// Value of the color (brightness). Valid values are [0, 100].
	Value int
}

func intToTemperature(i int) int {
	return int(colorCold - (colorCold-colorWarm)*float64(i)/100)
}

func temperatureToInt(temperature int) int {
	return int(float64(colorCold-temperature)*100.0/float64(colorCold-colorWarm) + 0.5)
}
