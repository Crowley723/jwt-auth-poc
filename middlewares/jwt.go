package middlewares

import (
	"jwt-auth-poc/crypt_utils"
	"net/http"
	"strings"
)

// RequireJWT is a middleware that validates JWT tokens
func RequireJWT(next func(*AppContext)) func(*AppContext) {
	return func(ctx *AppContext) {
		// Extract token from Authorization header
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.SetJSONError(http.StatusUnauthorized, "Missing authorization header")
			return
		}

		// Check for Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.SetJSONError(http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := parts[1]

		// Validate the token
		claims, err := ctx.JWTProvider.ValidateToken(token)
		if err != nil {
			ctx.Logger.Debug("JWT validation failed", "err", err)
			ctx.SetJSONError(http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Extract user ID from claims
		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			ctx.SetJSONError(http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Store user ID in context for handler use
		ctx.Set("user_id", userID)

		// Call the next handler
		next(ctx)
	}
}

// GetUserID retrieves the authenticated user ID from the context
func GetUserID(ctx *AppContext) string {
	if userID, ok := ctx.Get("user_id").(string); ok {
		return userID
	}
	return ""
}
