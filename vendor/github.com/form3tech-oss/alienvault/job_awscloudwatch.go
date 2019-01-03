package alienvault

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// AWSCloudWatchJob is a job which retrieves logs from cloudwatch groups(s)/stream(s)
type AWSCloudWatchJob struct {
	job
	Params AWSCloudWatchJobParams `json:"params"` // Params allows you to specify which region/group/stream you wish to retrieve logs from, and which plugin should be used to process those logs
}

// AWSCloudWatchJobParams allows you to specify cloudwatch job parameters
type AWSCloudWatchJobParams struct {
	jobParams
	Region string `json:"regionName"` // The region to use when retrieving logs from cloudwatch
	Group  string `json:"groupName"`  // The group to use when retrieving logs from cloudwatch
	Stream string `json:"streamName"` // The stream to use when retrieving logs from cloudwatch
}

func (job *AWSCloudWatchJob) enforceTypeValues() {
	job.Custom = true
	job.App = JobApplicationAWS
	job.Action = JobActionMonitorCloudWatch
	job.Type = JobTypeCollection
}

// GetAWSCloudWatchJobs returns all AWS CloudWatch jobs
func (client *Client) GetAWSCloudWatchJobs() ([]AWSCloudWatchJob, error) {

	req, err := client.createRequest("GET", "/scheduler", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	jobs := []AWSCloudWatchJob{}

	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	outputJobs := []AWSCloudWatchJob{}

	for _, job := range jobs {
		if job.Action == JobActionMonitorCloudWatch {
			outputJobs = append(outputJobs, job)
		}
	}

	return outputJobs, nil
}

// GetAWSCloudWatchJob returns a particular *AWSCloudWatchJob as identified by the UUID parameter
func (client *Client) GetAWSCloudWatchJob(uuid string) (*AWSCloudWatchJob, error) {

	// there is no individual GET endpoint for this, so we have to return all jobs and filter

	jobs, err := client.GetAWSCloudWatchJobs()
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

// CreateAWSCloudWatchJob creates a new AWS cloudwatch job
func (client *Client) CreateAWSCloudWatchJob(j *AWSCloudWatchJob) error {

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

	createdJob := AWSCloudWatchJob{}
	if err := json.NewDecoder(resp.Body).Decode(&createdJob); err != nil {
		return err
	}

	if createdJob.UUID == "" {
		return fmt.Errorf("failed to create the job")
	}

	j.UUID = createdJob.UUID
	return nil
}

// UpdateAWSCloudWatchJob updates an existing AWS cloudwatch job
func (client *Client) UpdateAWSCloudWatchJob(j *AWSCloudWatchJob) error {

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

// DeleteAWSCloudWatchJob deletes an existing AWS cloudwatch job
func (client *Client) DeleteAWSCloudWatchJob(j *AWSCloudWatchJob) error {

	req, err := client.createRequest("DELETE", fmt.Sprintf("/scheduler/%s", j.UUID), nil)
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
