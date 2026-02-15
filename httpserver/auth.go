package httpserver

import "net/http"

// RequestAuth contains authentication information for a request
type RequestAuth struct {
	Type   string      // Authentication type: "basic", "bearer", etc.
	Client string      // Client/microservice identifier
	Claims interface{} // Authentication claims (JWT payload, user info, etc.)
}

// AuthCallback is the function signature for authentication checking
type AuthCallback func(r *http.Request) (bool, RequestAuth)