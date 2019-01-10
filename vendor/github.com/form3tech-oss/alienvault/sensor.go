package alienvault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Sensor is a machine which gathers event data from your infrastrcture and absorbs it into the AV system
type Sensor struct {
	// Annoyingly, AV have two fields ID and UUID which both appear to be a primary key - but it is actually UUID that is used in APi calls and referenced in other resources. ID appears unused.
	UUID        string            `json:"uuid,omitempty"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      SensorStatus      `json:"status"`
	SetupStatus SensorSetupStatus `json:"setupStatus"`
}

type sensorActivation struct {
	//{"key":"${alienvault_sensor_key.main.id}","masterNode":"form3.alienvault.cloud","name":"${var.stack_name}-sensor","description":"${var.stack_name} sensor created by terraform"}
	SensorKey   string `json:"key"`
	MasterNode  string `json:"masterNode"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SensorStatus refers to whether or not the sensor is ready for jobs. "Ready" indicates that this is so.
type SensorStatus string

const (
	// SensorStatusReady indicates sensor is ready for configuration
	SensorStatusReady SensorStatus = "Ready"
)

type sensorSetupPatch struct {
	SetupStatus SensorSetupStatus `json:"setupStatus"`
}

// SensorSetupStatus refers to whether or not the sensor has had it's configuration finalised
type SensorSetupStatus string

const (
	// SensorSetupStatusComplete indicates sensor has had it's configuration finalised
	SensorSetupStatusComplete SensorSetupStatus = "Complete"
)

// waitForSensorToBeReady blocks until the given sensor is ready. Pass a context with timeout to abort after a set time.
func (client *Client) waitForSensorToBeReady(ctx context.Context, sensor *Sensor) error {

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context expired, no longer waiting for sensor to be ready")
		case <-ticker.C:
			s, err := client.GetSensor(sensor.UUID)
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
		if sensor.UUID == id {
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

// CreateSensorViaAppliance creates a new sensor via the sensor appliance referenced by the provided IP address
func (client *Client) CreateSensorViaAppliance(ctx context.Context, sensor *Sensor, ip net.IP) error {

	// first of all we need to make sure we can get our hands on an ath code (aka sensor key) to activate our new sensor
	// this may not be possible if we've maxed out the number of sensors on our license, so attempt this first and fail fast
	key, err := client.CreateSensorKey(false)
	if err != nil {
		return err
	}

	// wait until the sensor appliance has been created and is running an AV API over HTTP
	if err := client.waitForSensorApplianceCreation(ctx, ip); err != nil {
		return err
	}

	// the sensor appliance is alive! cool, now we can activate it with our auth code
	if err := client.activateSensorAppliance(ip, sensor, key); err != nil {
		return err
	}

	// hacky wait to ensure sensor is registered on the AV side
	time.Sleep(time.Second * 10)

	// TODO: we don't actually  know the ID of our new sensor yet, so until we figure that out, let's just look for a sensor that has an incomplete setupStatus. This is risky...
	sensors, err := client.GetSensors()
	if err != nil {
		return err
	}

	count := 0
	var createdSensor Sensor
	for _, s := range sensors {
		if s.SetupStatus != SensorSetupStatusComplete {
			count++
			if count > 1 {
				return fmt.Errorf("failed to complete sensor setup as we found more than one sensor being set up at the same time, and could differentiate between them")
			}
			createdSensor = s
		}
	}

	// we need the ID of the created sensor to complete setup
	sensor.UUID = createdSensor.UUID

	if err := client.completeSetup(&createdSensor); err != nil {
		return err
	}

	return client.waitForSensorToBeReady(ctx, sensor)
}

func (client *Client) waitForSensorApplianceCreation(ctx context.Context, ip net.IP) error {
	anonymousClient := &http.Client{
		Timeout: time.Second * 5,
	}

	url := fmt.Sprintf("http://%s/api/1.0/status", ip.String())

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	//keep hitting the sensor appliance every 10 seconds until it responds over http, or until context ends
	for {
		if resp, err := anonymousClient.Get(url); err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (client *Client) activateSensorAppliance(ip net.IP, sensor *Sensor, key *SensorKey) error {
	anonymousClient := &http.Client{
		Timeout: time.Second * 5,
	}

	activationPayload := sensorActivation{
		Name:        sensor.Name,
		Description: sensor.Description,
		SensorKey:   key.ID,
		MasterNode:  client.fqdn,
	}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(activationPayload); err != nil {
		return err
	}

	resp, err := anonymousClient.Post(fmt.Sprintf("http://%s/api/1.0/connect", ip.String()), "application/json;charset=UTF-8", b)
	if err != nil {
		return err
	}

	// todo remove this debug!
	data, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected HTTP status code on sensor activation: %d", resp.StatusCode)
	}

	return nil
}

// UpdateSensor updates an existing sensor
func (client *Client) UpdateSensor(sensor *Sensor) error {
	return fmt.Errorf("Not implemented")
}

// completeSetup marks a sensor as having it's setup finalised
func (client *Client) completeSetup(sensor *Sensor) error {

	sensorPatch := sensorSetupPatch{
		SetupStatus: SensorSetupStatusComplete,
	}

	data, err := json.Marshal(sensorPatch)
	if err != nil {
		return err
	}

	req, err := client.createRequest("PATCH", fmt.Sprintf("/sensors/%s", sensor.UUID), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Unexpected status code for sensor setup finalisation: %d", resp.StatusCode)
	}

	return nil
}

// DeleteSensor deletes an existing sensor
func (client *Client) DeleteSensor(sensor *Sensor) error {

	req, err := client.createRequest("DELETE", fmt.Sprintf("/sensors/%s", sensor.UUID), nil)
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
