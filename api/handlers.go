package api

import (
	"encoding/json"
	"http-sqlite-template/db"
	"http-sqlite-template/middlewares"
	"net/http"
	"strconv"
	"strings"
)

// handleHealthGET returns "OK"
func handleHealthGET(ctx *middlewares.AppContext) {
	// Also check database health
	if err := ctx.DB.Health(); err != nil {
		ctx.Logger.Error("database health check failed", "err", err)
		ctx.SetJSONError(http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	ctx.WriteText(http.StatusOK, "OK")
}

// handleUsersGET lists all users - limit 100
func handleUsersGET(ctx *middlewares.AppContext) {
	userQueries := db.NewUserQueries(ctx.DB)

	users, err := userQueries.List(100, 0)
	if err != nil {
		ctx.Logger.Error("failed to list users", "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	ctx.WriteJSON(http.StatusOK, map[string]interface{}{
		"users": users,
		"count": len(users),
	})
}

// handleUsersPOST creates a new user
func handleUsersPOST(ctx *middlewares.AppContext) {
	var request struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		ctx.SetJSONError(http.StatusBadRequest, "Invalid JSON")
		return
	}

	if request.Email == "" || request.Name == "" {
		ctx.SetJSONError(http.StatusBadRequest, "Email and name are required")
		return
	}

	userQueries := db.NewUserQueries(ctx.DB)
	user, err := userQueries.Create(request.Email, request.Name)
	if err != nil {
		ctx.Logger.Error("failed to create user", "err", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			ctx.SetJSONError(http.StatusConflict, "User with this email already exists")
			return
		}
		ctx.SetJSONError(http.StatusInternalServerError, "Failed to create user")
		return
	}

	ctx.WriteJSON(http.StatusCreated, user)
}

// handleUserGET retrieves a single user by ID
func handleUserGET(ctx *middlewares.AppContext) {
	idStr := ctx.Request.PathValue("id")

	if idStr == "" {
		ctx.Logger.Debug("Empty path value", "url", ctx.Request.URL.Path)
		ctx.SetJSONError(http.StatusBadRequest, "User ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.SetJSONError(http.StatusBadRequest, "Invalid user ID")
		return
	}

	userQueries := db.NewUserQueries(ctx.DB)
	user, err := userQueries.GetByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.SetJSONError(http.StatusNotFound, "User not found")
			return
		}
		ctx.Logger.Error("failed to get user", "err", err, "id", id)
		ctx.SetJSONError(http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	ctx.WriteJSON(http.StatusOK, user)
}

// handleUserDELETE removes a user by ID
func handleUserDELETE(ctx *middlewares.AppContext) {
	idStr := ctx.Request.PathValue("id")

	if idStr == "" {
		ctx.Logger.Debug("Empty path value for delete", "url", ctx.Request.URL.Path)
		ctx.SetJSONError(http.StatusBadRequest, "User ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.SetJSONError(http.StatusBadRequest, "Invalid user ID")
		return
	}

	userQueries := db.NewUserQueries(ctx.DB)
	if err := userQueries.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.SetJSONError(http.StatusNotFound, "User not found")
			return
		}
		ctx.Logger.Error("failed to delete user", "err", err, "id", id)
		ctx.SetJSONError(http.StatusInternalServerError, "Failed to delete user")
		return
	}

	ctx.SetJSONStatus(http.StatusOK, "User deleted successfully")
}
