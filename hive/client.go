/*
A simple library that allows you to interface with a REST API compatible with
Hive Home smart devices.
*/
package hive

import (
	"encoding/json"
	"errors"
	"path"
	"strings"
)

const (
	refreshDevicesTarget = "products?after="
	deviceTarget         = "nodes"
)

var (
	// Returned when the login fails due to the server not returning a valid token.
	ErrNoToken = errors.New("no token returned from server")

	// Returned when the login fails to the server not returning an endpoint.
	ErrNoEndpoint = errors.New("no endpoint URL returned from server")
)

type endpoint interface {
	PostJSON(url string, jsonStr []byte, token string) ([]byte, error)
	Get(url string, token string) ([]byte, error)
}

// Client is the main object used to obtain credentials and interface with the
// REST API.
type Client struct {
	// Token represents an authentication token used for all calls to the API.
	Token string

	// EndpointURL is the URL to the given API endpoint all calls will be sent to.
	EndpointURL string

	client  endpoint
	devices map[string]*Device
}

// NewClient returns a new client that's ready to use the Login method or use
// a saved token.
func NewClient() *Client {
	return &Client{client: &httpClient{}}
}

// Login uses the given credentials object to authenticate and obtain a token
// and an endpoint URL. It will also load the initial list of devices.
func (c *Client) Login(creds *Credentials) error {
	credsJSON, err := creds.toJSON(false, true)
	if err != nil {
		return err
	}
	resp, err := c.client.PostJSON(creds.URL, credsJSON, "")
	if err != nil {
		return err
	}
	auth, err := parseSession(resp)
	if err != nil {
		return err
	}

	c.Token = auth.Token
	c.EndpointURL = trailingSlash(auth.Platform.Endpoint)

	c.parseDevices(auth.Products)
	return nil
}

// RefreshDevices updates the devices available and their current states.
func (c *Client) RefreshDevices() error {
	url := c.buildURL(refreshDevicesTarget)
	resp, err := c.client.Get(url, c.Token)
	if err != nil {
		return err
	}

	var devices []jsonEntity
	if err := json.Unmarshal(resp, &devices); err != nil {
		return err
	}

	c.parseDevices(devices)
	return nil
}

// Device returns the device with the given ID or null if no such device could
// be found.
func (c *Client) Device(id string) *Device {
	if c.devices == nil {
		return nil
	}
	return c.devices[id]
}

// Devices returns a slice with all the devices present.
func (c *Client) Devices() []*Device {
	if c.devices == nil {
		return nil
	}

	result := make([]*Device, 0, len(c.devices))
	for _, device := range c.devices {
		result = append(result, device)
	}
	return result
}

func (c *Client) parseDevices(devices []jsonEntity) {
	if c.devices == nil {
		c.devices = make(map[string]*Device, len(devices))
	}

	for i := range devices {
		id := devices[i].ID
		device := c.Device(id)
		if device != nil {
			device.entity = &devices[i]
		} else {
			c.devices[id] = &Device{
				entity: &devices[i],
				client: c,
			}
		}
	}
}

func (c *Client) modifyDeviceState(device *Device, state *jsonState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	url := c.buildURL(deviceTarget, device.Type(), device.ID())
	_, err = c.client.PostJSON(url, data, c.Token)
	return err
}

func (c *Client) buildURL(params ...string) string {
	return c.EndpointURL + path.Join(params...)
}

// Credentials contains the user's username, password and URL used to perform
// the login operation.
type Credentials struct {
	Username string
	Password string
	URL      string
}

var (
	// Returned when the Credentials object contains no username.
	ErrCredsNoUsername = errors.New("no username specified")

	// Returned when the Credentials object contains no password.
	ErrCredsNoPassword = errors.New("no password specified")

	// Returned when the Credentials object contains no URL.
	ErrCredsNoURL = errors.New("no URL specified")
)

func (c *Credentials) toJSON(getDevices bool, getProducts bool) ([]byte, error) {
	if c.Username == "" {
		return nil, ErrCredsNoUsername
	}
	if c.Password == "" {
		return nil, ErrCredsNoPassword
	}
	if c.URL == "" {
		return nil, ErrCredsNoURL
	}

	jsonObj := jsonLogin{
		Username: c.Username,
		Password: c.Password,
		Devices:  getDevices,
		Products: getProducts,
	}

	return json.Marshal(jsonObj)
}

func parseSession(data []byte) (*jsonSession, error) {
	var session jsonSession
	err := json.Unmarshal(data, &session)
	if err != nil {
		return nil, err
	}

	if session.Token == "" {
		return nil, ErrNoToken
	}
	if session.Platform.Endpoint == "" {
		return nil, ErrNoEndpoint
	}

	return &session, nil
}

func trailingSlash(input string) string {
	return strings.TrimRight(input, "/") + "/"
}
