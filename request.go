package request

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"time"
)

type Config struct {
	TimeoutSeconds     int
	InsecureSkipVerify bool
}

var ConfigDefault = Config{
	TimeoutSeconds:     60,
	InsecureSkipVerify: true,
}

type httpRequest struct {
	Config Config
	Client http.Client
}

// Creates a new http client with the given configuration
// Use default config if not given any
func NewHttpRequestClient(config ...Config) httpRequest {
	if len(config) < 1 {
		return httpRequest{
			Config: ConfigDefault,
			Client: newHttpClient(ConfigDefault.TimeoutSeconds),
		}
	}

	cfg := config[0]

	return httpRequest{
		Config: cfg,
		Client: newHttpClient(cfg.TimeoutSeconds),
	}
}

// Sends a http GET request with given parameters
func (r httpRequest) Get(url string, headers map[string]string) (*http.Response, error) {
	return sendRequest(r.Client, http.MethodGet, url, headers, nil)
}

// Sends a POST request to the specified URL with the specified headers and body.
func (r httpRequest) Post(url string, headers map[string]string, body []byte) (*http.Response, error) {
	return sendRequest(r.Client, http.MethodPost, url, headers, body)
}

// Sends a PUT request to the specified URL with the specified headers and body.
func (r httpRequest) Put(url string, headers map[string]string, body []byte) (*http.Response, error) {
	return sendRequest(r.Client, http.MethodPut, url, headers, body)
}

// Sends a PATCH request to the specified URL with the specified headers and body.
func (r httpRequest) Patch(url string, headers map[string]string, body []byte) (*http.Response, error) {
	return sendRequest(r.Client, http.MethodPatch, url, headers, body)
}

// Sends a DELETE request to the specified URL with the specified headers and body.
func (r httpRequest) Delete(url string, headers map[string]string) (*http.Response, error) {
	return sendRequest(r.Client, http.MethodDelete, url, headers, nil)
}

func sendRequest(client http.Client, method string, url string, headers map[string]string, body []byte) (*http.Response, error) {

	reqBody := new(bytes.Buffer)
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	// Create a new request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	setResponseHeaders(req, headers)

	// Send the request
	resp, err := client.Do(req)

	return resp, err
}

// NewHttpClient creates a new http client with the specified timeout.
func newHttpClient(timeoutSeconds int) http.Client {
	return http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
}

// SetResponseHeaders adds the specified headers to the request.
func setResponseHeaders(req *http.Request, headers map[string]string) {
	// Add headers to the request
	has_content_type := false
	for key, value := range headers {
		req.Header.Add(key, value)
		if key == "Content-Type" {
			has_content_type = true
		}
	}

	// Add default content type if not set
	if !has_content_type {
		req.Header.Add("Content-Type", "application/json")
	}
}
