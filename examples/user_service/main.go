package main

import (
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/umakantv/go-utils/cache"
	"github.com/umakantv/go-utils/db"
	"github.com/umakantv/go-utils/examples/user_service/handlers"
	"github.com/umakantv/go-utils/httpserver"
	"github.com/umakantv/go-utils/logger"
)

// checkAuth implements authentication for the service
func checkAuth(r *http.Request) (bool, httpserver.RequestAuth) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return false, httpserver.RequestAuth{}
	}

	// Simple Bearer token check (in production, validate JWT)
	if len(auth) > 7 && auth[:7] == "Bearer " {
		token := auth[7:]
		if token == "secret-token" { // Simple check for demo
			return true, httpserver.RequestAuth{
				Type:   "bearer",
				Client: "user-service-client",
				Claims: map[string]interface{}{"service": "user-service"},
			}
		}
	}

	return false, httpserver.RequestAuth{}
}

func initializeDatabase() *sql.DB {
	// Database configuration for SQLite
	config := db.DatabaseConfig{
		DRIVER: "sqlite3",
		DB:     "./user_service.db",
	}

	dbConn := db.GetDBConnection(config)

	// Run schema
	schema, err := ioutil.ReadFile("./db/schema.sql")
	if err != nil {
		log.Fatal("Failed to read schema file:", err)
	}

	_, err = dbConn.Exec(string(schema))
	if err != nil {
		log.Fatal("Failed to execute schema:", err)
	}

	log.Println("Database initialized successfully")
	return dbConn
}

func initializeCache() cache.Cache {
	cache, err := cache.New(cache.Config{Type: "memory"})
	if err != nil {
		log.Fatal("Failed to initialize cache:", err)
	}
	return cache
}

func main() {
	// Initialize logger
	logger.Init(logger.LoggerConfig{
		CallerKey:  "file",
		TimeKey:    "timestamp",
		CallerSkip: 1,
	})

	logger.Info("Starting User Service...")

	// Initialize database
	dbConn := initializeDatabase()
	defer dbConn.Close()

	// Initialize cache
	cache := initializeCache()
	defer cache.Close()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(dbConn, cache)

	// Create HTTP server with authentication
	server := httpserver.New("8080", checkAuth)

	// Register routes
	server.Register(httpserver.Route{
		Name:     "HealthCheck",
		Method:   "GET",
		Path:     "/health",
		AuthType: "none",
	}, httpserver.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "user-service"}`))
	}))

	server.Register(httpserver.Route{
		Name:     "ListUsers",
		Method:   "GET",
		Path:     "/users",
		AuthType: "bearer",
	}, httpserver.HandlerFunc(userHandler.GetUsers))

	server.Register(httpserver.Route{
		Name:     "GetUser",
		Method:   "GET",
		Path:     "/users/{id}",
		AuthType: "bearer",
	}, httpserver.HandlerFunc(userHandler.GetUser))

	server.Register(httpserver.Route{
		Name:     "CreateUser",
		Method:   "POST",
		Path:     "/users",
		AuthType: "bearer",
	}, httpserver.HandlerFunc(userHandler.CreateUser))

	server.Register(httpserver.Route{
		Name:     "UpdateUser",
		Method:   "PUT",
		Path:     "/users/{id}",
		AuthType: "bearer",
	}, httpserver.HandlerFunc(userHandler.UpdateUser))

	server.Register(httpserver.Route{
		Name:     "DeleteUser",
		Method:   "DELETE",
		Path:     "/users/{id}",
		AuthType: "bearer",
	}, httpserver.HandlerFunc(userHandler.DeleteUser))

	logger.Info("User Service started on port 8080")
	logger.Info("Health check: GET /health")
	logger.Info("API endpoints: GET/POST/PUT/DELETE /users")

	// Start server
	if err := server.Start(); err != nil {
		logger.Error("Server failed to start", logger.Error(err))
	}
}
