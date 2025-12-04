package handlers

import (
	"encoding/json"
	"jwt-auth-poc/db"
	"jwt-auth-poc/middlewares"
	"jwt-auth-poc/utils"
	"net/http"
	"strconv"
	"strings"
)

// HandleRefreshTokenPost exchanges a refresh token for a new access token
func HandleRefreshTokenPost(ctx *middlewares.AppContext) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		ctx.SetJSONError(http.StatusBadRequest, "Invalid JSON")
		return
	}

	if strings.TrimSpace(request.RefreshToken) == "" {
		ctx.SetJSONError(http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokenHash := utils.HashToken(strings.TrimSpace(request.RefreshToken))

	refreshTokenQueries := db.NewRefreshTokenQueries(ctx.DB)
	refreshToken, err := refreshTokenQueries.GetByHashAndValidate(tokenHash)
	if err != nil {
		ctx.Logger.Debug("Invalid refresh token", "err", err)
		ctx.SetJSONError(http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	userID, err := strconv.Atoi(refreshToken.OwnerId)
	if err != nil {
		ctx.Logger.Error("Invalid user ID in refresh token", "owner_id", refreshToken.OwnerId, "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	userQueries := db.NewUserQueries(ctx.DB)
	user, err := userQueries.GetByID(userID)
	if err != nil {
		ctx.Logger.Error("Failed to get user", "user_id", userID, "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(ctx, user)
	if err != nil {
		ctx.Logger.Error("Failed to generate access token", "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}
	type Response struct {
		AccessToken string `json:"access_token"`
	}

	response := Response{
		AccessToken: newAccessToken,
	}

	ctx.WriteJSON(http.StatusOK, response)
}
