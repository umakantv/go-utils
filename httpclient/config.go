package httpclient

import "time"

// ClientConfig holds configuration for the HTTP client
type ClientConfig struct {
	// Timeout is the total timeout for requests
	Timeout time.Duration

	// MaxRetries is the default number of retries for requests at client level
	MaxRetries int

	// BaseHeaders are headers set on every request
	BaseHeaders map[string]string
}