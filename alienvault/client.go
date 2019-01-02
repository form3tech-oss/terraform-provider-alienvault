package alienvault

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

type Client struct {
	creds      Credentials
	fqdn       string
	urlPrefix  string
	httpClient *http.Client
}

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

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

func NewClient(fqdn string, creds Credentials) *Client {
	return &Client{
		fqdn:      fqdn,
		creds:     creds,
		urlPrefix: fmt.Sprintf("https://%s/api/1.0", fqdn),
	}
}

func (client *Client) createRequest(method string, path string, body io.Reader) (*http.Request, error) {

	// The 1.0 API requires the specific content type below and an X-XSRF-TOKEN header set to the value of the XSRF-TOKEN cookie

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", client.urlPrefix, path), body)
	if err != nil {
		return nil, err
	}
	cookies := client.httpClient.Jar.Cookies(req.URL)
	for i := range cookies {
		cookie := cookies[i]
		if cookie.Name == "XSRF-TOKEN" {
			req.Header.Set("X-XSRF-TOKEN", cookie.Value)
		}
	}
	req.Header.Set("Origin", fmt.Sprintf("https://%s", client.fqdn))
	req.Header.Set("Referer", fmt.Sprintf("https://%s/", client.fqdn))
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return req, nil
}

// Authenticate gives us a session for use in the AV API. Unfortunately schedules and other things we need are not supported in the public v2 REST API, so we have to use their internal one. The auth on this uses cookies, so we have to set this up here.
func (client *Client) Authenticate() error {

	credsData, err := json.Marshal(client.creds)
	if err != nil {
		return err
	}

	cookieJar, _ := cookiejar.New(nil)
	client.httpClient = &http.Client{
		Jar: cookieJar,
	}

	// skip TLS verification when running locally e.g. for testing
	if strings.HasPrefix(client.fqdn, "127.0.0.1:") {
		client.httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	// grab XSRF token etc.
	{
		_, err := client.httpClient.Get(fmt.Sprintf("https://%s/#/login", client.fqdn))
		if err != nil {
			return err
		}
	}

	// do login
	{
		req, err := client.createRequest("POST", "/login", bytes.NewBuffer(credsData))
		if err != nil {
			return err
		}

		resp, err := client.httpClient.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("Unexpected status code for auth: %d", resp.StatusCode)
		}
	}

	// get new csrf post-login
	{
		req, err := client.createRequest("GET", "/", nil)
		if err != nil {
			return err
		}

		_, err = client.httpClient.Do(req)
		if err != nil {
			return err
		}
	}

	return nil
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
