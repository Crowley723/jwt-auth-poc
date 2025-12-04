package api

import (
	"jwt-auth-poc/handlers"
	"jwt-auth-poc/middlewares"
	"net/http"
	"net/http/pprof"
)

func RegisterRoutes(mux *http.ServeMux, ctx *middlewares.AppContext) {
	// Serve static files
	//fs := http.FileServer(http.Dir("static/"))
	//mux.Handle("GET /", fs)

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// API routes
	mux.HandleFunc("GET /health", middlewares.Wrap(handlers.HandleHealthGET))

	mux.HandleFunc("GET /api/jwks.json", middlewares.Wrap(handlers.HandleJWKSPublicKeyGET))

	// Authentication routes
	mux.HandleFunc("POST /api/login", middlewares.Wrap(handlers.HandleUserLoginPost))
	mux.HandleFunc("POST /api/refresh", middlewares.Wrap(handlers.HandleRefreshTokenPost))

	// User management routes
	mux.HandleFunc("GET /api/users", middlewares.Wrap(handlers.HandleUsersGET))
	mux.HandleFunc("POST /api/users", middlewares.Wrap(handlers.HandleUsersPOST))
	mux.HandleFunc("GET /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// These don't use Wrap due to issues with path variables.
		appCtx := middlewares.GetOrCreateAppContext(r, w, ctx)
		handlers.HandleUserGET(appCtx)
	})
	mux.HandleFunc("DELETE /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// These don't use Wrap due to issues with path variables.
		appCtx := middlewares.GetOrCreateAppContext(r, w, ctx)
		handlers.HandleUserDELETE(appCtx)
	})

	// Protected routes (require JWT authentication)
	mux.HandleFunc("GET /api/protected/data", middlewares.Wrap(middlewares.RequireJWT(handlers.HandleProtectedDataGET)))
	mux.HandleFunc("GET /api/protected/stats", middlewares.Wrap(middlewares.RequireJWT(handlers.HandleProtectedStatsGET)))
}
