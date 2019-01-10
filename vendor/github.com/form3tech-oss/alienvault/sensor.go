package alienvault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

type applianceStatusResponse struct {
	Status applianceStatus `json:"status"`
}

type applianceStatus string

const (
	applianceStatusNotConnected applianceStatus = "notConnected"
)

// SensorStatus refers to whether or not the sensor is ready for jobs. "Ready" indicates that this is so.
type SensorStatus string

const (
	// SensorStatusReady indicates sensor is ready for configuration
	SensorStatusReady SensorStatus = "Ready"
	// SensorStatusConnectionLost refers to a sensor configuration which has lost contact with the actual appliance, possibly becuse the appliance no longer exists.
	SensorStatusConnectionLost SensorStatus = "Connection lost"
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

	// this usually takes 10-30 minutes so no need to poll that often
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {

		s, err := client.GetSensor(sensor.UUID)
		if err != nil {
			return err
		}

		if s.Status == SensorStatusReady {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}

}

func (client *Client) sweepSensors() error {

	sensors, err := client.GetSensors()
	if err != nil {
		return err
	}

	for _, sensor := range sensors {
		if sensor.Status == SensorStatusConnectionLost {
			if err := client.DeleteSensor(&sensor); err != nil {
				return err
			}
		}
	}

	return err
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

	log.Printf("[DEBUG] sweeping dead sensors...")

	// remove any dead sensors to free up license slots
	if err := client.sweepSensors(); err != nil {
		return err
	}

	// AV sometimes takes a few seconds to free up license slots after a sweep for some reason
	time.Sleep(time.Second * 5)

	log.Printf("[DEBUG] checking license...")
	if ok, err := client.HasSensorKeyAvailability(); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("the AlienVault license in use does not allow creation of more sensors")
	}

	log.Printf("[DEBUG] creating sensor key...")

	// first of all we need to make sure we can get our hands on an ath code (aka sensor key) to activate our new sensor
	// this may not be possible if we've maxed out the number of sensors on our license, so attempt this first and fail fast
	key, err := client.CreateSensorKey(false)
	if err != nil {
		return err
	}
	// ensure the key we create gets deleted if it isn't used for any reason
	defer func() {
		_ = client.DeleteSensorKey(key)
	}()

	log.Printf("[DEBUG] waiting for appliance to be created at %s...", ip.String())

	// wait until the sensor appliance has been created and is running an AV API over HTTP
	if err := client.waitForSensorApplianceCreation(ctx, ip); err != nil {
		return err
	}

	log.Printf("[DEBUG] activating sensor appliance...")

	// the sensor appliance is alive! cool, now we can activate it with our auth code
	if err := client.activateSensorAppliance(ctx, ip, sensor, key); err != nil {
		return err
	}

	// hacky wait to ensure sensor is registered on the AV side
	time.Sleep(time.Second * 10)

	log.Printf("[DEBUG] finding sensor to finish setup for...")

	// TODO: we don't actually  know the ID of our new sensor yet, so until we figure that out, let's just look for a sensor that has an incomplete setupStatus. This is risky...
	sensors, err := client.GetSensors()
	if err != nil {
		return err
	}

	count := 0
	var createdSensor Sensor
	for _, s := range sensors {
		if s.SetupStatus != SensorSetupStatusComplete && s.Name == sensor.Name {
			count++
			if count > 1 {
				return fmt.Errorf("failed to complete sensor setup as we found more than one sensor with the specified name being set up at the same time, and could differentiate between them")
			}
			createdSensor = s
		}
	}

	if count == 0 {
		return fmt.Errorf("no sensors found ready to be set up")
	}

	log.Printf("[DEBUG] completing setup...")

	// we need the ID of the created sensor to complete setup
	sensor.UUID = createdSensor.UUID

	if err := client.completeSetup(&createdSensor); err != nil {
		return err
	}

	log.Printf("[DEBUG] waiting for sensor to be live...")

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
		resp, err := anonymousClient.Get(url)
		if err == nil {
			b, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode == 200 {
				status := applianceStatusResponse{}
				if err := json.Unmarshal(b, &status); err == nil {
					if status.Status == applianceStatusNotConnected {
						break
					} else {
						return fmt.Errorf("Unexpected appliance status: %s", status.Status)
					}
				}

			} else {
				log.Printf("[ERROR] Status response code: %d", resp.StatusCode)
			}
		} else {
			log.Printf("[ERROR] Status check failed: %s", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}

	return nil
}

func (client *Client) activateSensorAppliance(ctx context.Context, ip net.IP, sensor *Sensor, key *SensorKey) error {
	anonymousClient := &http.Client{
		Timeout: time.Second * 5,
	}

	activationPayload := sensorActivation{
		Name:        sensor.Name,
		Description: sensor.Description,
		SensorKey:   key.ID,
		MasterNode:  client.fqdn,
	}

	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(activationPayload); err != nil {
			return err
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/1.0/connect", ip.String()), b)
		if err != nil {
			return err
		}
		req.Header.Set("Origin", fmt.Sprintf("https://%s", ip.String()))
		req.Header.Set("Referer", fmt.Sprintf("https://%s/", ip.String()))
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")

		if resp, err := anonymousClient.Do(req); err == nil {
			if resp.StatusCode == http.StatusOK {
				break
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
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
