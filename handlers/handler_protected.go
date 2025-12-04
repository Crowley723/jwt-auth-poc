package handlers

import (
	"jwt-auth-poc/db"
	"jwt-auth-poc/middlewares"
	"net/http"
	"strconv"
)

// HandleProtectedDataGET returns protected data for an authenticated user
func HandleProtectedDataGET(ctx *middlewares.AppContext) {
	// Get the authenticated user ID from context (set by middleware)
	userIDStr := middlewares.GetUserID(ctx)
	if userIDStr == "" {
		ctx.SetJSONError(http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.Logger.Error("Invalid user ID from JWT", "user_id", userIDStr, "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	userQueries := db.NewUserQueries(ctx.DB)
	user, err := userQueries.GetByID(userID)
	if err != nil {
		ctx.Logger.Error("Failed to get user", "user_id", userID, "err", err)
		ctx.SetJSONError(http.StatusNotFound, "User not found")
		return
	}

	type DataStats struct {
		ItemsProcessed int      `json:"items_processed"`
		LastAccess     string   `json:"last_access"`
		Permissions    []string `json:"permissions"`
	}

	type Response struct {
		Message string    `json:"message"`
		User    *db.User  `json:"user"`
		Data    DataStats `json:"data"`
	}

	response := Response{
		Message: "This is protected data",
		User:    user,
		Data: DataStats{
			ItemsProcessed: 1234,
			LastAccess:     "2024-01-15",
			Permissions:    []string{"read", "write", "delete"},
		},
	}

	ctx.WriteJSON(http.StatusOK, response)
}

// HandleProtectedStatsGET returns statistics for an authenticated user
func HandleProtectedStatsGET(ctx *middlewares.AppContext) {
	// Get the authenticated user ID from context
	userIDStr := middlewares.GetUserID(ctx)
	if userIDStr == "" {
		ctx.SetJSONError(http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.Logger.Error("Invalid user ID from JWT", "user_id", userIDStr, "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	refreshTokenQueries := db.NewRefreshTokenQueries(ctx.DB)
	tokens, err := refreshTokenQueries.GetValidByUserID(userID)
	if err != nil {
		ctx.Logger.Error("Failed to get refresh tokens", "user_id", userID, "err", err)
		tokens = []db.RefreshToken{}
	}

	// Return stats
	type Response struct {
		UserID             int `json:"user_id"`
		ActiveRefreshToken int `json:"active_refresh_tokens"`
		RequestCount       int `json:"request_count"`
	}

	response := Response{
		UserID:             userID,
		ActiveRefreshToken: len(tokens),
		RequestCount:       42,
	}

	ctx.WriteJSON(http.StatusOK, response)
}
