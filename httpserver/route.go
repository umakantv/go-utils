package httpserver

// Route defines a route for the HTTP server
type Route struct {
	Name     string
	Method   string
	Path     string
	AuthType string // "none", "basic", "bearer"
}