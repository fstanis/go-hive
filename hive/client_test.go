package hive

import (
	"encoding/json"
	"testing"
)

type mockEndpoint struct {
	url     string
	token   string
	payload string
	result  string
	err     error
}

func (c *mockEndpoint) PostJSON(url string, jsonStr []byte, token string) ([]byte, error) {
	c.url = url
	c.token = token
	c.payload = string(jsonStr)
	return []byte(c.result), c.err
}

func (c *mockEndpoint) Get(url string, token string) ([]byte, error) {
	return c.PostJSON(url, nil, token)
}

func (c *mockEndpoint) parsePayload() map[string]interface{} {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(c.payload), &m)
	return m
}

func (c *mockEndpoint) setResult(jsonData map[string]string) {
	result, _ := json.Marshal(jsonData)
	c.result = string(result)
}

func TestLogin(t *testing.T) {
	mock := &mockEndpoint{}
	client := &Client{client: mock}

	creds := &Credentials{"user", "secret", "http://example.com/"}
	mock.result = `
  {
    "token": "1234567890",
    "platform": {
      "endpoint": "https://example.com/api/version"
    }
  }
  `
	if err := client.Login(creds); err != nil {
		t.Errorf("client.Login returned error: %v", err)
	}
	if client.Token != "1234567890" {
		t.Errorf("client assigned token %q, want %q", client.Token, "1234567890")
	}
	if client.endpointBase != "https://example.com/api/version/" {
		t.Errorf("client assigned endpoint %q, want %q", client.endpointBase, "https://example.com/api/version")
	}

	if mock.url != creds.URL {
		t.Errorf("endpoint URL requested %q, want %q", mock.url, creds.URL)
	}
	payload := mock.parsePayload()
	if payload["username"] != creds.Username {
		t.Errorf("username sent %q, want %q", payload["username"], creds.Username)
	}
	if payload["password"] != creds.Password {
		t.Errorf("password sent %q, want %q", payload["password"], creds.Password)
	}
}

func TestLoginFail(t *testing.T) {
	mock := &mockEndpoint{}
	client := &Client{client: mock}
	creds := &Credentials{"user", "secret", "http://example.com/"}

	if err := client.Login(creds); err == nil {
		t.Error("client.Login returned no error for invalid response")
	}

	mock.result = "{}"
	if err := client.Login(creds); err == nil {
		t.Error("client.Login returned no error for invalid response")
	}
}

func TestRefreshDevices(t *testing.T) {
	mock := &mockEndpoint{}
	client := &Client{client: mock}

	mock.result = `
	[
		{"id":"12345678-abcd","type":"warmwhitelight"},
		{"id":"01234567-abcd","type":"warmwhitelight"}
	]
	`
	err := client.RefreshDevices()
	if err != nil {
		t.Errorf("client.RefreshDevices returned error: %v", err)
	}
	if len(client.devices) != 2 {
		t.Errorf("Client parsed %d devices, expected 2", len(client.devices))
	}
	if client.Device("01234567-abcd") == nil {
		t.Error("Device 01234567-abcd not found after parsing")
	}
}
