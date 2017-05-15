package hive

import (
	"fmt"
	"time"
)

// Device represents a reference to a single device. It can be a light, sensor
// thermostat or anything else.
type Device struct {
	entity *jsonEntity
	client *Client
}

// Do sends the request to apply the given change to this device.
func (d *Device) Do(c *Change) error {
	return d.client.modifyDeviceState(d, &c.state)
}

// ID returns the unique ID of this device.
func (d *Device) ID() string {
	return d.entity.ID
}

// Type returns the type of this device.
func (d *Device) Type() string {
	return d.entity.Type
}

// Name returns the user-given name of this device, or an empty string if no
// name is given.
func (d *Device) Name() string {
	if d.entity.State.Name == nil {
		return ""
	}
	return *d.entity.State.Name
}

// String returns a string representation of this device, containing the ID,
// name and type.
func (d *Device) String() string {
	name := d.Name()
	if name == "" {
		name = "(unnamed)"
	}
	return fmt.Sprintf("[%s] %s (%s)", d.ID(), name, d.Type())
}

// Created returns the time when this device was added.
func (d *Device) Created() time.Time {
	return time.Time(d.entity.Created)
}

// LastSeen returns the time when this device was last online.
func (d *Device) LastSeen() time.Time {
	return time.Time(d.entity.LastSeen)
}

// IsOnline returns true if this device is currently powered on and connected,
// false otherwise.
func (d *Device) IsOnline() bool {
	return d.entity.Props.Online != nil && *d.entity.Props.Online
}

// Getters specific to motion sensors

// IsMotionSensor checks if this device is a motion sensor.
func (d *Device) IsMotionSensor() bool {
	return d.Type() == typeMotionSensor
}

// HasMotion returns true if this device is a motion sensor and is currently
// detecting motion.
func (d *Device) HasMotion() bool {
	if d.entity.Props.Motion == nil {
		return false
	}
	return d.entity.Props.Motion.Status
}

// LastMotionStart returns the start time of the last detected motion by this
// device, if it's a motion sensor.
func (d *Device) LastMotionStart() time.Time {
	if d.entity.Props.Motion == nil {
		return time.Time{}
	}
	return time.Time(d.entity.Props.Motion.Start)
}

// LastMotionEnd returns the end time of the last detected motion by this
// device, if it's a motion sensor.
func (d *Device) LastMotionEnd() time.Time {
	if d.entity.Props.Motion == nil {
		return time.Time{}
	}
	return time.Time(d.entity.Props.Motion.End)
}

// Getters specific to lights

// IsLight returns true if this device is a light bulb of any kind.
func (d *Device) IsLight() bool {
	return d.Type() == typeWarmWhiteLight || d.IsColorLight()
}

// IsColorLight returns true if this device is a color light bulb.
func (d *Device) IsColorLight() bool {
	return d.Type() == typeColourLight
}

// IsOn returns true if this device is a light bulb and is currently turned on.
func (d *Device) IsOn() bool {
	return d.entity.State.Status != nil && *d.entity.State.Status == statusON
}

// Brightness returns the current brightness level of this light bulb, between
// 0 and 100.
func (d *Device) Brightness() int {
	if d.entity.State.Brightness == nil {
		return 0
	}
	return *d.entity.State.Brightness
}

// Color returns the current color set on this colored light bulb. Will return
// the last used color if the device is not currently in color mode or turned
// off.
func (d *Device) Color() HSV {
	if d.entity.State.Hue == nil || d.entity.State.Saturation == nil || d.entity.State.Value == nil {
		return HSV{}
	}
	return HSV{
		*d.entity.State.Hue,
		*d.entity.State.Saturation,
		*d.entity.State.Value,
	}
}

// ColorTemperature returns the current temperature of this colored light bulb
// which can be between 0 and 100, with 100 being the warmest possible setting.
// Will return the last temperature value if the device is off or in color mode.
func (d *Device) ColorTemperature() int {
	if d.entity.State.ColourTemperature == nil {
		return 0
	}
	temp := *d.entity.State.ColourTemperature
	return temperatureToInt(temp)
}

func (d *Device) Mode() string {
	if d.entity.State.Mode == nil {
		return ""
	}
	return *d.entity.State.Mode
}
