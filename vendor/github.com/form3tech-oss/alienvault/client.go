package alienvault

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

// Client is an API client for interacting with AlienVault USM Anywhere
type Client struct {
	creds               Credentials
	fqdn                string
	urlPrefix           string
	httpClient          *http.Client
	skipTLSVerification bool
	version             int
}

// Credentials contain a username and password for accessing the AV USM system
type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

// New creates a new client using the provided FQDN and credentials
func New(fqdn string, creds Credentials, skipTLSVerification bool, version int) *Client {
	return &Client{
		version:             version,
		fqdn:                fqdn,
		creds:               creds,
		skipTLSVerification: skipTLSVerification,
		urlPrefix:           fmt.Sprintf("https://%s/api/%d.0", fqdn, version),
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

// Authenticate gives the client a session to use in subsequent calls.
func (client *Client) Authenticate() error {

	// Unfortunately job schedules and other things we need are not supported in the public v2 REST API,
	// so we have to use their internal one. The auth on this uses cookies, so we have to set this up here.

	credsData, err := json.Marshal(client.creds)
	if err != nil {
		return err
	}

	cookieJar, _ := cookiejar.New(nil)
	client.httpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: client.skipTLSVerification || strings.HasPrefix(client.fqdn, "127.0.0.1:"),
			},
		},
		Jar: cookieJar,
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
			d, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("Unexpected status code for auth: %d: %s", resp.StatusCode, string(d))
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
