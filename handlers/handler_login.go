package handlers

import (
	"encoding/json"
	"jwt-auth-poc/db"
	"jwt-auth-poc/middlewares"
	"jwt-auth-poc/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-crypt/crypt"
)

func HandleUserLoginPost(ctx *middlewares.AppContext) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		ctx.SetJSONError(http.StatusBadRequest, "Invalid JSON")
		return
	}

	if strings.TrimSpace(request.Email) == "" || strings.TrimSpace(request.Password) == "" {
		ctx.SetJSONError(http.StatusBadRequest, "Email and password are required")
		return
	}

	userQueries := db.NewUserQueries(ctx.DB)
	userDetails, err := userQueries.GetUserDetailsByEmail(strings.TrimSpace(request.Email))
	if err != nil {
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	if valid, err := crypt.CheckPassword(strings.TrimSpace(request.Password), userDetails.PasswordHash); err != nil || !valid {
		ctx.Logger.Debug("Failed login attempt", "email", userDetails.Email, "err", err)
		ctx.SetJSONError(http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, hash, err := utils.GenerateRefreshToken()
	if err != nil {
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	refreshTokenQueries := db.NewRefreshTokenQueries(ctx.DB)
	newRefreshToken, err := refreshTokenQueries.Create(strconv.Itoa(userDetails.ID), hash)
	if err != nil {
		ctx.Logger.Error("failed to save new refresh token", "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(ctx, userDetails)
	if err != nil {
		ctx.Logger.Error("failed to generate access token", "err", err)
		ctx.SetJSONError(http.StatusInternalServerError, "Internal server error")
		return
	}

	type Response struct {
		RefreshToken       string `json:"refresh_token"`
		RefreshTokenExpiry int64  `json:"refresh_token_expiry"`
		AccessToken        string `json:"access_token"`
	}
	var response = Response{
		RefreshToken:       token,
		RefreshTokenExpiry: newRefreshToken.ExpiresAt.Unix(),
		AccessToken:        newAccessToken,
	}

	ctx.WriteJSON(http.StatusOK, response)
}
