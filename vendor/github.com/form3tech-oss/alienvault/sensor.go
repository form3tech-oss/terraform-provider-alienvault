package alienvault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Sensor is a machine which gathers event data from your infrastrcture and absorbs it into the AV system
type Sensor struct {
	ID          string       `json:"id,omitempty"`
	UUID        string       `json:"uuid,omitempty"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Status      SensorStatus `json:"status"`
}

// SensorStatus refers to whether or not the sensor is ready for jobs. "Ready" indicates that this is so.
type SensorStatus string

const (
	// SensorStatusReady indicates sensor is ready for configuration
	SensorStatusReady SensorStatus = "Ready"
)

// WaitForSensorToBeReady blocks until the given sensor is ready. Pass a context with timeout to abort after a set time.
func (client *Client) WaitForSensorToBeReady(ctx context.Context, sensor *Sensor) error {

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context expired, no longer waiting for sensor to be ready")
		case <-ticker.C:
			s, err := client.GetSensor(sensor.ID)
			if err != nil {
				return err
			}

			if s.Status == SensorStatusReady {
				return nil
			}
		}
	}

}

// GetSensor returns a specific sensor as identified by the id parameter
func (client *Client) GetSensor(id string) (*Sensor, error) {

	req, err := client.createRequest("GET", "/sensors", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	sensors := []Sensor{}

	if err := json.NewDecoder(resp.Body).Decode(&sensors); err != nil {
		return nil, err
	}

	for _, sensor := range sensors {
		if sensor.ID == id {
			return &sensor, nil
		}
	}

	return nil, fmt.Errorf("Sensor %s could not be found", id)
}

// GetSensors returns a list of all sensors
func (client *Client) GetSensors() ([]Sensor, error) {

	req, err := client.createRequest("GET", "/sensors", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	sensors := []Sensor{}

	if err := json.NewDecoder(resp.Body).Decode(&sensors); err != nil {
		return nil, err
	}

	return sensors, nil
}

// CreateSensor creates a new sensor
func (client *Client) CreateSensor(sensor *Sensor) error {

	data, err := json.Marshal(sensor)
	if err != nil {
		return err
	}

	req, err := client.createRequest("POST", "/sensors", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	createdSensor := Sensor{}
	if err := json.NewDecoder(resp.Body).Decode(&createdSensor); err != nil {
		return err
	}

	sensor.ID = createdSensor.ID
	return nil
}

// UpdateSensor updates an existing sensor
func (client *Client) UpdateSensor(sensor *Sensor) error {

	data, err := json.Marshal(sensor)
	if err != nil {
		return err
	}

	req, err := client.createRequest("PATCH", fmt.Sprintf("/sensors/%s", sensor.ID), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	createdSensor := Sensor{}
	if err := json.NewDecoder(resp.Body).Decode(&createdSensor); err != nil {
		return err
	}

	sensor.ID = createdSensor.ID
	return nil
}

// DeleteSensor deletes an existing sensor
func (client *Client) DeleteSensor(sensor *Sensor) error {

	req, err := client.createRequest("DELETE", fmt.Sprintf("/sensors/%s", sensor.ID), nil)
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Unexpected status code on delete: %d", resp.StatusCode)
	}

	return nil
}
