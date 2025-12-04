package handlers

import (
	"encoding/json"
	"jwt-auth-poc/crypt_utils"
	"jwt-auth-poc/db"
	"jwt-auth-poc/middlewares"
	"net/http"
	"strconv"
	"strings"
)

// HandleUsersGET lists all users - limit 100
func HandleUsersGET(ctx *middlewares.AppContext) {
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

// HandleUsersPOST creates a new user
func HandleUsersPOST(ctx *middlewares.AppContext) {
	var request struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		ctx.SetJSONError(http.StatusBadRequest, "Invalid JSON")
		return
	}

	if request.Email == "" || request.Name == "" || request.Password == "" {
		ctx.SetJSONError(http.StatusBadRequest, "email, name, and password are required")
		return
	}

	hashedPassword, err := crypt_utils.HashPassword(request.Password)
	if err != nil {
		ctx.Logger.Error("failed to hash password", "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	userQueries := db.NewUserQueries(ctx.DB)
	user, err := userQueries.Create(request.Email, request.Name, hashedPassword)
	if err != nil {
		ctx.Logger.Error("failed to create user", "err", err)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			ctx.SetJSONError(http.StatusInternalServerError, "Failed to create user")
			return
		}
		ctx.SetJSONError(http.StatusInternalServerError, "Failed to create user")
		return
	}

	ctx.WriteJSON(http.StatusCreated, user)
}

// HandleUserGET retrieves a single user by ID
func HandleUserGET(ctx *middlewares.AppContext) {
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

// HandleUserDELETE removes a user by ID
func HandleUserDELETE(ctx *middlewares.AppContext) {
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
