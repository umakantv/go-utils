package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/umakantv/go-utils/logger"
)

// Server represents the HTTP server
type Server struct {
	router       *mux.Router
	port         string
	authCallback AuthCallback
}

// New creates a new HTTP server with authentication callback
func New(port string, authCallback AuthCallback) *Server {
	return &Server{
		router:       mux.NewRouter(),
		port:         port,
		authCallback: authCallback,
	}
}

// Register registers a route with its handler
func (s *Server) Register(route Route, handler Handler) {
	s.router.HandleFunc(route.Path, s.wrapHandler(route, handler)).Methods(route.Method).Name(route.Name)
}

// wrapHandler wraps the handler with authentication, logging, and context injection
func (s *Server) wrapHandler(route Route, handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var requestAuth *RequestAuth

		// Handle authentication
		if route.AuthType != "none" {
			if s.authCallback == nil {
				http.Error(w, "Authentication callback not configured", http.StatusInternalServerError)
				return
			}

			authenticated, auth := s.authCallback(r)
			if !authenticated {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			requestAuth = &auth
		}

		// Log the request
		logger.Info(fmt.Sprintf("Received request: %s - %s - %s", route.Name, route.Method, r.URL.Path))

		// Inject request details into context
		ctx = context.WithValue(ctx, RouteNameKey, route.Name)
		ctx = context.WithValue(ctx, RouteMethodKey, route.Method)
		ctx = context.WithValue(ctx, RoutePathKey, route.Path)
		ctx = context.WithValue(ctx, AuthTypeKey, route.AuthType)
		if requestAuth != nil {
			ctx = context.WithValue(ctx, RequestAuthKey, *requestAuth)
		}

		// Call the handler
		handler.Handle(ctx, w, r)
	}
}



// Start starts the HTTP server
func (s *Server) Start() error {
	fmt.Printf("Starting server on port %s\n", s.port)
	return http.ListenAndServe(":"+s.port, s.router)
}