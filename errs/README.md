# Errors Module

The errors module provides standardized error handling for microservices with predefined error types and HTTP status code mappings.

## AppError Struct

The core error type is `AppError`:

```go
type AppError struct {
    Code    int    `json:",omitempty"`
    Message string
}
```

## Predefined Error Constructors

### HTTP 404 - Not Found
```go
err := NewNotFoundError("Resource not found")
```

### HTTP 500 - Internal Server Error
```go
err := NewInternalServerError("Database connection failed")
```

### HTTP 422 - Validation Error
```go
err := NewValidationError("Invalid input data")
```

### HTTP 401 - Authentication Error
```go
err := NewAuthenticationError("Invalid credentials")
```

### HTTP 403 - Authorization Error
```go
err := NewAuthorizationError("Access denied")
```

## Usage in Handlers

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    user, err := findUser(r.URL.Query().Get("id"))
    if err != nil {
        if dbErr, ok := err.(*AppError); ok {
            w.WriteHeader(dbErr.Code)
            json.NewEncoder(w).Encode(dbErr)
            return
        }
        // Handle other errors
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}
```

## Error Interface

`AppError` implements the `error` interface:

```go
func (e AppError) Error() string {
    return e.Message
}
```

## AsMessage Method

Convert error to message-only format:

```go
err := NewValidationError("Email is required")
msgOnly := err.AsMessage() // Returns AppError with only Message field
```