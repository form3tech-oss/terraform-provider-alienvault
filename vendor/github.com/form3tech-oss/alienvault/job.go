package alienvault

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Job struct {
	UUID        string                 `json:"uuid,omitempty"`
	SensorID    string                 `json:"sensor"`
	App         string                 `json:"app"`
	Action      string                 `json:"action"`
	Schedule    string                 `json:"schedule"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"`
	Disabled    bool                   `json:"disabled"`
	Params      map[string]interface{} `json:"params,omitempty"`
	Running     bool                   `json:"running,omitempty"`
	Custom      bool                   `json:"custom"`
	LastRun     int                    `json:"lastRun,omitempty"`
	NextRun     int                    `json:"nextRun,omitempty"`
}

func (client *Client) GetJob(uuid string) (*Job, error) {

	req, err := client.createRequest("GET", "/scheduler", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	jobs := []Job{}

	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	for _, job := range jobs {
		if job.UUID == uuid {
			return &job, nil
		}
	}

	return nil, fmt.Errorf("Job %s could not be found", uuid)
}

func (client *Client) CreateJob(job *Job) error {

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	req, err := client.createRequest("POST", "/scheduler", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	createdJob := Job{}
	if err := json.NewDecoder(resp.Body).Decode(&createdJob); err != nil {
		return err
	}

	job.UUID = createdJob.UUID
	return nil
}

func (client *Client) UpdateJob(job *Job) error {

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	req, err := client.createRequest("PUT", fmt.Sprintf("/scheduler/%s", job.UUID), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	createdJob := Job{}
	if err := json.NewDecoder(resp.Body).Decode(&createdJob); err != nil {
		return err
	}

	job.UUID = createdJob.UUID
	return nil
}

func (client *Client) DeleteJob(uuid string) error {

	req, err := client.createRequest("DELETE", fmt.Sprintf("/scheduler/%s", uuid), nil)
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
