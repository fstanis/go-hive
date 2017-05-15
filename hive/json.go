package hive

import (
	"encoding/json"
	"time"
)

type jsonError struct {
	ErrorText *string `json:"error"`
}

type jsonTimestamp time.Time

func (j *jsonTimestamp) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}
	t := time.Unix(timestamp/1000, timestamp%1000*1000000)
	*j = jsonTimestamp(t)
	return nil
}

type jsonSession struct {
	Token    string       `json:"token"`
	Status   string       `json:"status"`
	User     jsonUser     `json:"user"`
	Platform jsonPlatform `json:"platform"`
	Alerts   jsonAlerts   `json:"alerts"`
	Devices  []jsonEntity `json:"devices"`
	Products []jsonEntity `json:"products"`
}

type jsonUser struct {
	ID              string `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Country         string `json:"country"`
	CountryCode     string `json:"countryCode"`
	Locale          string `json:"locale"`
	Postcode        string `json:"postcode"`
	Mobile          string `json:"mobile"`
	TemperatureUnit string `json:"temperatureUnit"`
	Timezone        string `json:"timezone"`
}

type jsonLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Devices  bool   `json:"devices"`
	Products bool   `json:"products"`
}

type jsonPlatform struct {
	Endpoint  string `json:"endpoint"`
	Honeycomb string `json:"honeycomb"`
	Name      string `json:"name"`
}

type jsonAlerts struct {
	FailuresEmail bool `json:"failuresEmail"`
	FailuresSMS   bool `json:"failuresSMS"`
	NightAlerts   bool `json:"nightAlerts"`
	WarningsEmail bool `json:"warningsEmail"`
	WarningsSMS   bool `json:"warningsSMS"`
}

type jsonEntity struct {
	ID        string        `json:"id"`
	Type      string        `json:"type"`
	Created   jsonTimestamp `json:"created"`
	LastSeen  jsonTimestamp `json:"lastSeen"`
	Parent    string        `json:"parent"`
	Props     jsonProps     `json:"props"`
	State     jsonState     `json:"state"`
	SortOrder int           `json:"sortOrder"`
}

type jsonProps struct {
	Manufacturer *string     `json:"manufacturer"`
	Model        *string     `json:"model"`
	Online       *bool       `json:"online"`
	Version      *string     `json:"version"`
	Power        *string     `json:"power"`
	Signal       *int        `json:"signal"`
	Connection   *string     `json:"connection"`
	IPAddress    *string     `json:"ipAddress"`
	Migrating    *bool       `json:"migrating"`
	PMZ          *string     `json:"pmz"`
	Uptime       *int        `json:"uptime"`
	Motion       *jsonMotion `json:"motion"`
}

type jsonState struct {
	Name      *string       `json:"name,omitempty"`
	Mode      *string       `json:"mode,omitempty"`
	Status    *string       `json:"status,omitempty"`
	Discovery *bool         `json:"discovery,omitempty"`
	Schedule  *jsonSchedule `json:"schedule,omitempty"`

	// Light
	Brightness        *int    `json:"brightness,omitempty"`
	ColourMode        *string `json:"colourMode,omitempty"`
	ColourTemperature *int    `json:"colourTemperature,omitempty"`
	Hue               *int    `json:"hue,omitempty"`
	Saturation        *int    `json:"saturation,omitempty"`
	Value             *int    `json:"value,omitempty"`
}

type jsonMotion struct {
	Status bool          `json:"status"`
	Start  jsonTimestamp `json:"start"`
	End    jsonTimestamp `json:"end"`
}

type jsonSchedule struct {
	Monday    []jsonScheduleEntry `json:"monday"`
	Tuesday   []jsonScheduleEntry `json:"tuesday"`
	Wednesday []jsonScheduleEntry `json:"wednesday"`
	Thursday  []jsonScheduleEntry `json:"thursday"`
	Friday    []jsonScheduleEntry `json:"friday"`
	Saturday  []jsonScheduleEntry `json:"saturday"`
	Sunday    []jsonScheduleEntry `json:"sunday"`
}

type jsonScheduleEntry struct {
	Start int                  `json:"start"`
	Value jsonScheduleDayValue `json:"value"`
}

type jsonScheduleDayValue struct {
	Brightness int    `json:"brightness"`
	Status     string `json:"status"`
}
