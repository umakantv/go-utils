# HTTP Client Module

The HTTP client module provides a configurable HTTP client with built-in support for authentication, retries, timeouts, and JSON request/response handling.

## Configuration

Configure the client using `ClientConfig`:

```go
type ClientConfig struct {
    Timeout     time.Duration // Request timeout (default: 30s)
    MaxRetries  int           // Max retry attempts
    BaseHeaders map[string]string // Headers added to all requests
}
```

## Creating a Client

```go
config := ClientConfig{
    Timeout:    10 * time.Second,
    MaxRetries: 3,
    BaseHeaders: map[string]string{
        "User-Agent": "MyMicroservice/1.0",
    },
}

client := New(config)
```

## Request Options

Customize individual requests with options:

```go
// Authentication
WithAuth("Bearer your-token-here")

// Custom headers
WithHeaders(map[string]string{"X-API-Key": "key"})

// Custom retries
WithRetries(5)

// Request body
WithBody(strings.NewReader("data"))
```

## Basic HTTP Methods

```go
// GET request
resp, err := client.Get("https://api.example.com/users")

// POST with body
resp, err := client.Post("https://api.example.com/users", WithBody(reader))

// PUT request
resp, err := client.Put("https://api.example.com/users/123", WithBody(reader))

// PATCH request
resp, err := client.Patch("https://api.example.com/users/123", WithBody(reader))

// DELETE request
resp, err := client.Delete("https://api.example.com/users/123")
```

## JSON Methods

### Sending JSON

```go
userData := []byte(`{"name": "John", "email": "john@example.com"}`)

// POST JSON
resp, err := client.PostJSON("https://api.example.com/users", userData)

// PUT JSON
resp, err := client.PutJSON("https://api.example.com/users/123", userData)

// PATCH JSON
resp, err := client.PatchJSON("https://api.example.com/users/123", userData)
```

### Receiving JSON

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var user User

// GET with JSON response
err := client.GetJSON("https://api.example.com/users/123", &user)

// POST with JSON request and response
userData := []byte(`{"name": "Jane"}`)
err := client.PostJSON("https://api.example.com/users", userData, &user)

// DELETE with JSON response (if any)
err := client.DeleteJSON("https://api.example.com/users/123", nil)
```

## Authentication Examples

### Bearer Token
```go
resp, err := client.Get("https://api.example.com/protected",
    WithAuth("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."))
```

### API Key
```go
resp, err := client.Get("https://api.example.com/data",
    WithHeaders(map[string]string{"X-API-Key": "your-api-key"}))
```

### Basic Auth
```go
import "encoding/base64"

auth := base64.StdEncoding.EncodeToString([]byte("username:password"))
resp, err := client.Get("https://api.example.com/basic-auth",
    WithHeaders(map[string]string{"Authorization": "Basic " + auth}))
```

## Retry Behavior

- Retries on network errors and 5xx HTTP status codes
- Uses exponential backoff: 1s, 2s, 3s, etc.
- Configurable per client and per request
- Client-level `MaxRetries` applies to all requests unless overridden

## Error Handling

JSON methods return errors for:
- Network failures
- Non-2xx HTTP status codes
- JSON parsing errors

```go
var user User
err := client.GetJSON("https://api.example.com/users/123", &user)
if err != nil {
    // Handle error (network, HTTP error, or JSON parsing)
    log.Printf("Failed to get user: %v", err)
    return
}
// Use user data
```

## Advanced Usage

### Custom Client Configuration
```go
config := ClientConfig{
    Timeout:    5 * time.Second,
    MaxRetries: 5,
    BaseHeaders: map[string]string{
        "Authorization": "Bearer default-token",
        "Content-Type":  "application/json",
    },
}
client := New(config)
```

### Per-Request Customization
```go
resp, err := client.Get("https://api.example.com/data",
    WithHeaders(map[string]string{"X-Custom": "value"}),
    WithRetries(1),
    WithAuth("Bearer request-specific-token"))
```

## Best Practices

1. Set reasonable timeouts (5-30 seconds)
2. Use appropriate retry counts (1-5)
3. Include authentication in base config or per-request
4. Handle errors properly in calling code
5. Use JSON methods for API integrations
6. Log requests/responses for debugging