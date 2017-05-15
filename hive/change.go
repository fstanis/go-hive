package hive

import "errors"

// Change represents a single update to the device's current state. For example,
// for a light bulb, a change may include both turning it on and setting the
// brightness to a certain level, which will be done in a single request via the
// Do method of a device.
//
// Change methods are built to be chained, so an example use would be:
//     someDevice.Do(hive.NewChange().TurnOn().Brightness(50))
type Change struct {
	state jsonState
}

// NewChange simply returns an empty Change object.
func NewChange() *Change {
	return &Change{}
}

// Brightness makes this change set the brightness to the given value. Valid
// values are numbers between 0 and 100 (inclusive). An invalid value will
// result in a panic.
func (c *Change) Brightness(brightness int) *Change {
	if brightness < 0 || brightness > 100 {
		panic(errors.New("brightness must be between 0 and 100"))
	}
	c.state.Brightness = &brightness
	c.resetHSV()
	return c
}

// Name sets the device name this change will apply.
func (c *Change) Name(name string) *Change {
	c.state.Name = &name
	return c
}

// Color makes this change set the color to the given value. It will also change
// the lightbulb mode to *color* if set to *color temperature*. This method will
// panic if either of hue, saturation or value are set to an invalid value.
func (c *Change) Color(hsv HSV) *Change {
	if hsv.Hue < 0 || hsv.Hue > 359 {
		panic(errors.New("hue must be between 0 and 359"))
	}
	if hsv.Saturation < 0 || hsv.Saturation > 99 {
		panic(errors.New("saturation must be between 0 and 99"))
	}
	if hsv.Value < 0 || hsv.Value > 100 {
		panic(errors.New("value must be between 0 and 100"))
	}
	c.state.ColourMode = &colourModeCOLOUR
	c.state.Hue = &hsv.Hue
	c.state.Saturation = &hsv.Saturation
	c.state.Value = &hsv.Value
	c.state.Brightness = nil
	c.state.ColourTemperature = nil
	return c
}

// ColorTemperature makes this change set the color temperature to the given
// value. It will also change the lightbulb mode to *color temperature* if set
// to *color*.Valid values are numbers between 0 and 100 (inclusive). An invalid
// value will result in a panic.
func (c *Change) ColorTemperature(temperature int) *Change {
	if temperature < 0 || temperature > 100 {
		panic(errors.New("temperature must be between 0 and 1"))
	}

	teperatureVal := intToTemperature(temperature)
	c.state.ColourMode = &colourModeWHITE
	c.state.ColourTemperature = &teperatureVal
	c.resetHSV()
	return c
}

// TurnOn makes this change turn the light on.
func (c *Change) TurnOn() *Change {
	c.state.Status = &statusON
	return c
}

// TurnOff makes this change turn the light off.
func (c *Change) TurnOff() *Change {
	c.state.Status = &statusOFF
	return c
}

func (c *Change) resetHSV() {
	c.state.Hue = nil
	c.state.Saturation = nil
	c.state.Value = nil
}
