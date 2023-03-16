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

// newHttpClient creates a new http client with the given configuration
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

func (r httpRequest) Get(url string, headers map[string]string) (*http.Response, error) {
	// Create client for the service
	client := newHttpClient(r.Config.TimeoutSeconds)

	// Create a new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	setResponseHeaders(req, headers)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SendRequestPost sends a POST request to the specified URL with the specified headers and body.
func (r httpRequest) Post(url string, headers map[string]string, body []byte) (*http.Response, error) {

	// Create client for the service
	client := newHttpClient(r.Config.TimeoutSeconds)

	// Create a new request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
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
