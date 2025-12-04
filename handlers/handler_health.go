package handlers

import (
	"jwt-auth-poc/middlewares"
	"net/http"
)

// HandleHealthGET returns "OK"
func HandleHealthGET(ctx *middlewares.AppContext) {
	if err := ctx.DB.Health(); err != nil {
		ctx.Logger.Error("database health check failed", "err", err)
		ctx.SetJSONError(http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	ctx.WriteText(http.StatusOK, "OK")
}
