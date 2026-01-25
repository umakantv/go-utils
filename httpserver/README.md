# HTTP Server Module

The HTTP server module provides a standardized way to build HTTP servers for microservices with built-in routing, authentication, logging, and context injection.

## Dependencies

- `github.com/gorilla/mux` for advanced routing with path parameters

## Route Definition

Define routes using the `Route` struct:

```go
type Route struct {
    Name     string // Unique route identifier
    Method   string // HTTP method: "GET", "POST", "PUT", "PATCH", "DELETE"
    Path     string // URL path with optional parameters (e.g., "/users/{id}")
    AuthType string // Authentication type: "none", "basic", "bearer"
}
```

## RequestAuth Structure

Authentication details are provided via `RequestAuth`:

```go
type RequestAuth struct {
    Type   string      // Authentication type ("basic", "bearer")
    Client string      // Client/microservice identifier
    Claims interface{} // Authentication claims (JWT payload, user info, etc.)
}
```

## Handler Interface

Handlers must implement the `Handler` interface:

```go
type Handler interface {
    Handle(ctx context.Context, w http.ResponseWriter, r *http.Request)
}
```

Use `HandlerFunc` for function-based handlers:

```go
type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)
```

## Authentication Callback

Define an authentication callback function:

```go
func checkAuth(r *http.Request) (bool, httpserver.RequestAuth) {
    auth := r.Header.Get("Authorization")
    if auth == "" {
        return false, httpserver.RequestAuth{}
    }

    // Example: Bearer token validation
    if strings.HasPrefix(auth, "Bearer ") {
        token := strings.TrimPrefix(auth, "Bearer ")
        // Validate token and extract claims
        claims, err := validateJWT(token)
        if err != nil {
            return false, httpserver.RequestAuth{}
        }

        return true, httpserver.RequestAuth{
            Type:   "bearer",
            Client: claims.ClientID,
            Claims: claims,
        }
    }

    // Example: Basic auth
    user, pass, ok := r.BasicAuth()
    if !ok {
        return false, httpserver.RequestAuth{}
    }

    // Validate credentials
    if validateUser(user, pass) {
        return true, httpserver.RequestAuth{
            Type:   "basic",
            Client: user,
            Claims: nil, // or user info
        }
    }

    return false, httpserver.RequestAuth{}
}
```

## Creating a Server

```go
server := httpserver.New("8080", checkAuth) // Port and auth callback
```

## Registering Routes

```go
server.Register(httpserver.Route{
    Name:     "GetUser",
    Method:   "GET",
    Path:     "/users/{id}",
    AuthType: "bearer",
}, httpserver.HandlerFunc(getUserHandler))

server.Register(httpserver.Route{
    Name:     "CreateUser",
    Method:   "POST",
    Path:     "/users",
    AuthType: "basic",
}, httpserver.HandlerFunc(createUserHandler))
```

## Starting the Server

```go
err := server.Start()
if err != nil {
    log.Fatal("Server failed to start:", err)
}
```

## Complete Example

```go
package main

import (
    "context"
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/umakantv/go-utils/httpserver"
    "github.com/umakantv/go-utils/logger"
)

func checkAuth(r *http.Request) (bool, httpserver.RequestAuth) {
    // Implement your authentication logic here
    auth := r.Header.Get("Authorization")
    if auth == "" {
        return false, httpserver.RequestAuth{}
    }

    // Example JWT validation
    if strings.HasPrefix(auth, "Bearer ") {
        token := strings.TrimPrefix(auth, "Bearer ")
        claims, err := validateJWT(token) // Your JWT validation function
        if err != nil {
            return false, httpserver.RequestAuth{}
        }

        return true, httpserver.RequestAuth{
            Type:   "bearer",
            Client: claims.ClientID,
            Claims: claims,
        }
    }

    return false, httpserver.RequestAuth{}
}

func main() {
    // Initialize logger
    logger.Init(logger.LoggerConfig{
        CallerKey:  "file",
        TimeKey:    "timestamp",
        CallerSkip: 1,
    })

    // Create server with auth callback
    server := httpserver.New("8080", checkAuth)

    // Register routes
    server.Register(httpserver.Route{
        Name:     "HealthCheck",
        Method:   "GET",
        Path:     "/health",
        AuthType: "none",
    }, httpserver.HandlerFunc(healthCheckHandler))

    server.Register(httpserver.Route{
        Name:     "GetUser",
        Method:   "GET",
        Path:     "/users/{id}",
        AuthType: "bearer",
    }, httpserver.HandlerFunc(getUserHandler))

    server.Register(httpserver.Route{
        Name:     "CreateUser",
        Method:   "POST",
        Path:     "/users",
        AuthType: "basic",
    }, httpserver.HandlerFunc(createUserHandler))

    // Start server
    logger.Info("Starting server on port 8080")
    if err := server.Start(); err != nil {
        logger.Error("Server failed to start", logger.Error(err))
    }
}

func healthCheckHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "healthy"}`))
}

func getUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    // Access route metadata from context
    routeName := httpserver.GetRouteName(ctx)
    method := httpserver.GetRouteMethod(ctx)
    path := httpserver.GetRoutePath(ctx)
    authType := httpserver.GetAuthType(ctx)

    // Access authentication details
    requestAuth := httpserver.GetRequestAuth(ctx)
    if requestAuth != nil {
        logger.Info("Authenticated request",
            logger.String("client", requestAuth.Client),
            logger.String("auth_type", requestAuth.Type))
    }

    // Get path parameters
    vars := mux.Vars(r)
    userID := vars["id"]

    logger.Info("Processing get user request",
        logger.String("user_id", userID),
        logger.String("route", routeName),
        logger.String("client", requestAuth.Client))

    // Mock response
    user := map[string]interface{}{
        "id":   userID,
        "name": "John Doe",
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}

func createUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    // Access metadata
    routeName := httpserver.GetRouteName(ctx)

    var user map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        logger.Error("Failed to decode user data", logger.Error(err))
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    logger.Info("Creating new user",
        logger.String("route", routeName),
        logger.Any("user_data", user))

    // Mock response
    user["id"] = "123"
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

## Authentication Types

### None Authentication
```go
Route{
    Name:     "PublicEndpoint",
    Method:   "GET",
    Path:     "/public",
    AuthType: "none",
}
```
No authentication required.

### Basic Authentication
```go
Route{
    Name:     "AdminEndpoint",
    Method:   "POST",
    Path:     "/admin",
    AuthType: "basic",
}
```
Expects `Authorization: Basic <base64-encoded-credentials>` header.

### Bearer Token Authentication
```go
Route{
    Name:     "APIEndpoint",
    Method:   "GET",
    Path:     "/api/data",
    AuthType: "bearer",
}
```
Expects `Authorization: Bearer <token>` header.

## Context Metadata

Every request automatically injects metadata into the context:

- `RouteName`: The route's name
- `RouteMethod`: HTTP method
- `RoutePath`: Route path template
- `AuthType`: Authentication type used
- `RequestAuth`: Authentication details (only for authenticated requests)

Access using helper functions:

```go
routeName := httpserver.GetRouteName(ctx)
method := httpserver.GetRouteMethod(ctx)
path := httpserver.GetRoutePath(ctx)
authType := httpserver.GetAuthType(ctx)
requestAuth := httpserver.GetRequestAuth(ctx)
if requestAuth != nil {
    client := requestAuth.Client
    claims := requestAuth.Claims
}
```

## Automatic Request Logging

Every incoming request is automatically logged in the format:
```
"Received request: {route name} - {METHOD} - {actual path}"
```

Example log output:
```
Received request: GetUser - GET - /users/123
```

## Path Parameters

Use Gorilla Mux syntax for path parameters:

```go
// Single parameter
"/users/{id}"

// Multiple parameters
"/users/{userId}/posts/{postId}"

// Optional parameters
"/users/{id:[0-9]+}" // Only numeric IDs
```

Access parameters in handlers:

```go
vars := mux.Vars(r)
userID := vars["id"]
postID := vars["postId"]
```

## Middleware Chain

The server applies middleware in this order:
1. **Authentication**: Calls auth callback for non-"none" routes, injects `RequestAuth`
2. **Logging**: Logs the incoming request
3. **Context Injection**: Adds route metadata and auth details to context
4. **Handler Execution**: Calls your handler function

## Error Handling

Authentication failures return HTTP 401 Unauthorized.

Handle other errors in your handlers:

```go
func myHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    if err := processRequest(r); err != nil {
        logger.Error("Request processing failed", logger.Error(err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
```

## Best Practices

1. Initialize logger before starting server
2. Implement robust authentication callback
3. Use descriptive route names
4. Choose appropriate authentication types per route
5. Access context metadata for logging and authorization
6. Handle authentication errors properly
7. Use path parameters for RESTful APIs
8. Keep handlers focused and testable
9. Validate authentication claims in handlers
10. Use JSON for request/response bodies

## Extending Authentication

The authentication logic is placeholder. Extend the `authenticate` method in `server.go` to integrate with your authentication system:

```go
func (s *Server) authenticate(authType string, r *http.Request) error {
    switch authType {
    case "bearer":
        token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
        return validateJWT(token)
    case "basic":
        user, pass, ok := r.BasicAuth()
        if !ok {
            return errors.New("invalid basic auth")
        }
        return validateCredentials(user, pass)
    }
    return nil
}
```