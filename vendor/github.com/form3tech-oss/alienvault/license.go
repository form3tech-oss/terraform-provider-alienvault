package alienvault

import (
	"encoding/json"
	"time"
)

// License is an AV license subscription
type License struct {
	ControlNodeLimit int   `json:"controlNodesAllowed"`
	SensorNodeLimit  int   `json:"sensorNodesAllowed"`
	MonthlyStorageKB int64 `json:"monthlyKBStorage"`
	Expiration       int64 `json:"expiration"`
}

// IsExpired returns true if the license in use has expired
func (license *License) IsExpired() bool {
	return time.Unix(license.Expiration, 0).Before(time.Now())
}

// GetLicense returns the license in use by the current account
func (client *Client) GetLicense() (*License, error) {

	req, err := client.createRequest("GET", "/license", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	license := License{}

	if err := json.NewDecoder(resp.Body).Decode(&license); err != nil {
		return nil, err
	}

	return &license, nil
}

// HasSensorAvailability tells us whether we have room to create new sensors using the current license
func (client *Client) HasSensorAvailability() (bool, error) {

	sensors, err := client.GetSensors()
	if err != nil {
		return false, err
	}

	license, err := client.GetLicense()
	if err != nil {
		return false, err
	}

	return len(sensors) < license.SensorNodeLimit, nil
}

// HasSensorKeyAvailability tells us whether we have room to create new sensor keys using the current license
func (client *Client) HasSensorKeyAvailability() (bool, error) {

	sensors, err := client.GetSensors()
	if err != nil {
		return false, err
	}

	keys, err := client.GetSensorKeys()
	if err != nil {
		return false, err
	}

	license, err := client.GetLicense()
	if err != nil {
		return false, err
	}

	return len(sensors)+len(keys) < license.SensorNodeLimit, nil
}
