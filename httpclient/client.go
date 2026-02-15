package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the HTTP client with configurable options
type Client struct {
	httpClient *http.Client
	config     ClientConfig
}

// New creates a new HTTP client with the given config
func New(config ClientConfig) *Client {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries < 0 {
		config.MaxRetries = 0
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
	}
}

// RequestOption is a function to configure request options
type RequestOption func(*RequestOptions)

// RequestOptions holds options for a single request
type RequestOptions struct {
	Headers map[string]string
	Retries int // overrides client MaxRetries if set to positive
	Body    io.Reader
}

// WithHeaders adds custom headers to the request
func WithHeaders(headers map[string]string) RequestOption {
	return func(opts *RequestOptions) {
		if opts.Headers == nil {
			opts.Headers = make(map[string]string)
		}
		for k, v := range headers {
			opts.Headers[k] = v
		}
	}
}

// WithAuth adds authorization header
func WithAuth(token string) RequestOption {
	return WithHeaders(map[string]string{"Authorization": token})
}

// WithRetries sets the number of retries for this request
func WithRetries(retries int) RequestOption {
	return func(opts *RequestOptions) {
		opts.Retries = retries
	}
}

// WithBody sets the request body
func WithBody(body io.Reader) RequestOption {
	return func(opts *RequestOptions) {
		opts.Body = body
	}
}

// Do performs the HTTP request with retries and options
func (c *Client) Do(method, url string, opts ...RequestOption) (*http.Response, error) {
	reqOpts := &RequestOptions{}
	for _, opt := range opts {
		opt(reqOpts)
	}

	var req *http.Request
	var err error

	if reqOpts.Body != nil {
		req, err = http.NewRequest(method, url, reqOpts.Body)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}

	// Set base headers
	for k, v := range c.config.BaseHeaders {
		req.Header.Set(k, v)
	}

	// Override with request headers
	for k, v := range reqOpts.Headers {
		req.Header.Set(k, v)
	}

	retries := c.config.MaxRetries
	if reqOpts.Retries > 0 {
		retries = reqOpts.Retries
	}

	var resp *http.Response
	for attempt := 0; attempt <= retries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		// Retry on network errors or 5xx status codes
		if attempt < retries {
			time.Sleep(time.Duration(attempt+1) * time.Second) // exponential backoff
		}
	}

	return resp, err
}

// Get performs a GET request
func (c *Client) Get(url string, opts ...RequestOption) (*http.Response, error) {
	return c.Do(http.MethodGet, url, opts...)
}

// Post performs a POST request
func (c *Client) Post(url string, opts ...RequestOption) (*http.Response, error) {
	return c.Do(http.MethodPost, url, opts...)
}

// Put performs a PUT request
func (c *Client) Put(url string, opts ...RequestOption) (*http.Response, error) {
	return c.Do(http.MethodPut, url, opts...)
}

// Patch performs a PATCH request
func (c *Client) Patch(url string, opts ...RequestOption) (*http.Response, error) {
	return c.Do(http.MethodPatch, url, opts...)
}

// Delete performs a DELETE request
func (c *Client) Delete(url string, opts ...RequestOption) (*http.Response, error) {
	return c.Do(http.MethodDelete, url, opts...)
}

// GetJSON performs a GET request and unmarshals the JSON response into result
func (c *Client) GetJSON(url string, result interface{}, opts ...RequestOption) error {
	resp, err := c.Do(http.MethodGet, url, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

// DeleteJSON performs a DELETE request and unmarshals the JSON response into result
func (c *Client) DeleteJSON(url string, result interface{}, opts ...RequestOption) error {
	resp, err := c.Do(http.MethodDelete, url, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

// PostJSON performs a POST request with JSON body and unmarshals the JSON response into result
func (c *Client) PostJSON(url string, jsonBody []byte, result interface{}, opts ...RequestOption) error {
	opts = append(opts, WithHeaders(map[string]string{"Content-Type": "application/json"}), WithBody(bytes.NewReader(jsonBody)))
	resp, err := c.Do(http.MethodPost, url, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

// PutJSON performs a PUT request with JSON body and unmarshals the JSON response into result
func (c *Client) PutJSON(url string, jsonBody []byte, result interface{}, opts ...RequestOption) error {
	opts = append(opts, WithHeaders(map[string]string{"Content-Type": "application/json"}), WithBody(bytes.NewReader(jsonBody)))
	resp, err := c.Do(http.MethodPut, url, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}

// PatchJSON performs a PATCH request with JSON body and unmarshals the JSON response into result
func (c *Client) PatchJSON(url string, jsonBody []byte, result interface{}, opts ...RequestOption) error {
	opts = append(opts, WithHeaders(map[string]string{"Content-Type": "application/json"}), WithBody(bytes.NewReader(jsonBody)))
	resp, err := c.Do(http.MethodPatch, url, opts...)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if result != nil {
		return json.Unmarshal(body, result)
	}
	return nil
}