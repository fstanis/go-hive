package hive

import (
	"testing"
)

func TestDoChange(t *testing.T) {
	mock := &mockEndpoint{}
	token := "abc"
	client := &Client{client: mock, Token: token}
	device := &Device{&jsonEntity{}, client}

	device.Do(NewChange().TurnOn().Brightness(55))
	if mock.token != token {
		t.Errorf("Do sent token %q, want %q", mock.token, token)
	}
	payload := mock.parsePayload()
	if payload["status"] != statusON {
		t.Errorf("State status set to %q, want %q", payload["status"], statusON)
	}
	if payload["brightness"] != 55.0 {
		t.Errorf("State brightness set to %f, want 55", payload["brightness"])
	}
}
