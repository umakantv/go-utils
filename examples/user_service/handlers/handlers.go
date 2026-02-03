package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/umakantv/go-utils/cache"
	"github.com/umakantv/go-utils/errs"
	"github.com/umakantv/go-utils/examples/user_service/models"
	"github.com/umakantv/go-utils/httpserver"
	"github.com/umakantv/go-utils/logger"
)

// UserHandler handles user-related operations
type UserHandler struct {
	db    *sql.DB
	cache cache.Cache
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *sql.DB, cache cache.Cache) *UserHandler {
	return &UserHandler{
		db:    db,
		cache: cache,
	}
}

// logRequest logs the request with the specified format
func (h *UserHandler) logRequest(ctx context.Context, level string, message string, fields ...logger.Field) {
	routeName := httpserver.GetRouteName(ctx)
	method := httpserver.GetRouteMethod(ctx)
	path := httpserver.GetRoutePath(ctx)
	auth := httpserver.GetRequestAuth(ctx)

	// Build log message
	logMsg := time.Now().Format("2006-01-02 15:04:05") + " - " + routeName + " - " + method + " - " + path
	if auth != nil {
		logMsg += " - client:" + auth.Client
	}

	// Add custom fields
	allFields := append([]logger.Field{
		logger.String("route", routeName),
		logger.String("method", method),
		logger.String("path", path),
	}, fields...)

	switch level {
	case "info":
		logger.Info(logMsg, allFields...)
	case "error":
		logger.Error(logMsg, allFields...)
	case "debug":
		logger.Debug(logMsg, allFields...)
	}
}

// GetUsers handles GET /users - list all users
func (h *UserHandler) GetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.logRequest(ctx, "info", "Listing users")

	// Try cache first
	cacheKey := "users:list"
	if cached, err := h.cache.Get(cacheKey); err == nil {
		h.logRequest(ctx, "debug", "Serving from cache")
		w.Header().Set("Content-Type", "application/json")
		w.Write(cached.([]byte))
		return
	}

	// Query database
	rows, err := h.db.Query("SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC")
	if err != nil {
		h.logRequest(ctx, "error", "Failed to query users", logger.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errs.NewInternalServerError("Database error"))
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			h.logRequest(ctx, "error", "Failed to scan user", logger.Error(err))
			continue
		}
		users = append(users, user)
	}

	// Cache the result
	response, _ := json.Marshal(users)
	h.cache.Set(cacheKey, response, 5*time.Minute)

	h.logRequest(ctx, "info", "Users retrieved successfully", logger.Int("count", len(users)))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// GetUser handles GET /users/{id} - get user by ID
func (h *UserHandler) GetUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logRequest(ctx, "error", "Invalid user ID", logger.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs.NewValidationError("Invalid user ID"))
		return
	}

	h.logRequest(ctx, "info", "Getting user", logger.Int("user_id", id))

	// Try cache first
	cacheKey := "user:" + idStr
	if cached, err := h.cache.Get(cacheKey); err == nil {
		h.logRequest(ctx, "debug", "Serving user from cache", logger.Int("user_id", id))
		w.Header().Set("Content-Type", "application/json")
		w.Write(cached.([]byte))
		return
	}

	// Query database
	var user models.User
	err = h.db.QueryRow("SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		h.logRequest(ctx, "info", "User not found", logger.Int("user_id", id))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errs.NewNotFoundError("User not found"))
		return
	}
	if err != nil {
		h.logRequest(ctx, "error", "Failed to query user", logger.Error(err), logger.Int("user_id", id))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errs.NewInternalServerError("Database error"))
		return
	}

	// Cache the result
	response, _ := json.Marshal(user)
	h.cache.Set(cacheKey, response, 10*time.Minute)

	h.logRequest(ctx, "info", "User retrieved successfully", logger.Int("user_id", id))

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// CreateUser handles POST /users - create a new user
func (h *UserHandler) CreateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logRequest(ctx, "error", "Invalid request body", logger.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs.NewValidationError("Invalid JSON"))
		return
	}

	// Validate input
	if req.Name == "" || req.Email == "" {
		h.logRequest(ctx, "error", "Missing required fields", logger.String("name", req.Name), logger.String("email", req.Email))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs.NewValidationError("Name and email are required"))
		return
	}

	h.logRequest(ctx, "info", "Creating user", logger.String("name", req.Name), logger.String("email", req.Email))

	// Insert user
	result, err := h.db.Exec("INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)",
		req.Name, req.Email, time.Now(), time.Now())
	if err != nil {
		h.logRequest(ctx, "error", "Failed to create user", logger.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errs.NewInternalServerError("Failed to create user"))
		return
	}

	id, _ := result.LastInsertId()
	userID := int(id)

	// Clear users list cache
	h.cache.Delete("users:list")

	h.logRequest(ctx, "info", "User created successfully", logger.Int("user_id", userID))

	// Return created user
	user := models.User{
		ID:        userID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser handles PUT /users/{id} - update user
func (h *UserHandler) UpdateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logRequest(ctx, "error", "Invalid user ID", logger.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs.NewValidationError("Invalid user ID"))
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logRequest(ctx, "error", "Invalid request body", logger.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs.NewValidationError("Invalid JSON"))
		return
	}

	h.logRequest(ctx, "info", "Updating user", logger.Int("user_id", id))

	// Build update query dynamically
	setParts := []string{}
	args := []interface{}{}

	if req.Name != "" {
		setParts = append(setParts, "name = ?")
		args = append(args, req.Name)
	}
	if req.Email != "" {
		setParts = append(setParts, "email = ?")
		args = append(args, req.Email)
	}

	if len(setParts) == 0 {
		h.logRequest(ctx, "error", "No fields to update", logger.Int("user_id", id))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs.NewValidationError("No fields to update"))
		return
	}

	setParts = append(setParts, "updated_at = ?")
	args = append(args, time.Now())
	args = append(args, id)

	query := "UPDATE users SET " + strings.Join(setParts, ", ") + " WHERE id = ?"
	result, err := h.db.Exec(query, args...)
	if err != nil {
		h.logRequest(ctx, "error", "Failed to update user", logger.Error(err), logger.Int("user_id", id))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errs.NewInternalServerError("Failed to update user"))
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.logRequest(ctx, "info", "User not found for update", logger.Int("user_id", id))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errs.NewNotFoundError("User not found"))
		return
	}

	// Clear caches
	h.cache.Delete("users:list")
	h.cache.Delete("user:" + idStr)

	h.logRequest(ctx, "info", "User updated successfully", logger.Int("user_id", id))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

// DeleteUser handles DELETE /users/{id} - delete user
func (h *UserHandler) DeleteUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logRequest(ctx, "error", "Invalid user ID", logger.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errs.NewValidationError("Invalid user ID"))
		return
	}

	h.logRequest(ctx, "info", "Deleting user", logger.Int("user_id", id))

	// Delete user
	result, err := h.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		h.logRequest(ctx, "error", "Failed to delete user", logger.Error(err), logger.Int("user_id", id))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errs.NewInternalServerError("Failed to delete user"))
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		h.logRequest(ctx, "info", "User not found for deletion", logger.Int("user_id", id))
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errs.NewNotFoundError("User not found"))
		return
	}

	// Clear caches
	h.cache.Delete("users:list")
	h.cache.Delete("user:" + idStr)

	h.logRequest(ctx, "info", "User deleted successfully", logger.Int("user_id", id))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}