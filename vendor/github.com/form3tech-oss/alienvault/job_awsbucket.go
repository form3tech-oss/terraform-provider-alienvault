package alienvault

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// AWSBucketJob is a scheduled job for retrieving logs from an S3 bucket
type AWSBucketJob struct {
	job
	Params AWSBucketJobParams `json:"params"` // Params allows you to dictate which bucket and path to use for the job, and specify which plugin should be used to process the logs.
}

// AWSBucketJobParams are parameters for an AWSBucketJob
type AWSBucketJobParams struct {
	jobParams
	BucketName string `json:"bucketName"` // The name of the bucket to use when retrieving logs for this job
	Path       string `json:"path"`       // The path to use when looking for logs in the specified bucket
}

func (job *AWSBucketJob) enforceTypeValues() {
	job.Custom = true
	job.App = JobApplicationAWS
	job.Action = JobActionMonitorBucket
	job.Type = JobTypeCollection
}

// GetAWSBucketJobs returns a slice of all AWS Bucket jobs
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

// GetAWSBucketJob returns a particular *AWSBucketJob as identified by the UUID parameter
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

// CreateAWSBucketJob creates a new bucket job
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

// UpdateAWSBucketJob updates an AWS bucket job
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

// DeleteAWSBucketJob deletes a bucket job
func (client *Client) DeleteAWSBucketJob(j *AWSBucketJob) error {

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
