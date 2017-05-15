package hive

import "testing"

func TestChange(t *testing.T) {
	change := NewChange().Brightness(55).TurnOff().TurnOn().Name("test")
	if *change.state.Brightness != 55 {
		t.Errorf("Change has brightness set to %d, expected 55", *change.state.Brightness)
	}
	if *change.state.Status != statusON {
		t.Errorf("Change has status set to %q, expected %q", *change.state.Status, statusON)
	}
	if *change.state.Name != "test" {
		t.Errorf("Change has name set to %s, expected test", *change.state.Name)
	}
}
