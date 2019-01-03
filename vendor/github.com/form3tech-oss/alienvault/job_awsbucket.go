package alienvault

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type AWSBucketJob struct {
	job
	Params AWSBucketJobParams `json:"params"`
}

type AWSBucketJobParams struct {
	jobParams
	BucketName string `json:"bucketName"`
	Path       string `json:"path"`
}

func (job *AWSBucketJob) enforceTypeValues() {
	job.Custom = true
	job.App = JobApplicationAWS
	job.Action = JobActionMonitorBucket
	job.Type = JobTypeCollection
}

func (client *Client) GetAWSBucketJobs() ([]AWSBucketJob, error) {

	req, err := client.createRequest("GET", "/scheduler", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	jobs := []AWSBucketJob{}

	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	outputJobs := []AWSBucketJob{}

	for _, job := range jobs {
		if job.Action == JobActionMonitorBucket {
			outputJobs = append(outputJobs, job)
		}
	}

	return outputJobs, nil
}

func (client *Client) GetAWSBucketJob(uuid string) (*AWSBucketJob, error) {

	// there is no individual GET endpoint for this, so we have to return all jobs and filter

	jobs, err := client.GetAWSBucketJobs()
	if err != nil {
		return nil, err
	}

	for _, job := range jobs {
		if job.UUID == uuid {
			return &job, nil
		}
	}

	return nil, fmt.Errorf("Job %s could not be found", uuid)
}

func (client *Client) CreateAWSBucketJob(j *AWSBucketJob) error {

	if j.UUID != "" {
		return fmt.Errorf("you cannot specify a UUID when creating a job")
	}

	// force values for this subtype of job
	j.enforceTypeValues()

	data, err := json.Marshal(j)
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

	createdJob := AWSBucketJob{}
	if err := json.NewDecoder(resp.Body).Decode(&createdJob); err != nil {
		return err
	}

	if createdJob.UUID == "" {
		return fmt.Errorf("failed to create the job")
	}

	j.UUID = createdJob.UUID
	return nil
}

func (client *Client) UpdateAWSBucketJob(j *AWSBucketJob) error {

	// force values for this subtype of job
	j.enforceTypeValues()

	data, err := json.Marshal(j)
	if err != nil {
		return err
	}

	req, err := client.createRequest("PUT", fmt.Sprintf("/scheduler/%s", j.UUID), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	createdJob := job{}
	if err := json.NewDecoder(resp.Body).Decode(&createdJob); err != nil {
		return err
	}

	j.UUID = createdJob.UUID
	return nil
}

func (client *Client) DeleteAWSBucketJob(uuid string) error {

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
