package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"http-sqlite-template/middlewares"
)

func StartServer(ctx *middlewares.AppContext) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", middlewares.Wrap(handleHealthGET))

	mux.HandleFunc("GET /api/users", middlewares.Wrap(handleUsersGET))
	mux.HandleFunc("POST /api/users", middlewares.Wrap(handleUsersPOST))
	mux.HandleFunc("GET /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// These don't use Wrap due to issues with path variables.
		appCtx := middlewares.GetOrCreateAppContext(r, w, ctx)
		handleUserGET(appCtx)
	})
	mux.HandleFunc("DELETE /api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		// These don't use Wrap due to issues with path variables.
		appCtx := middlewares.GetOrCreateAppContext(r, w, ctx)
		handleUserDELETE(appCtx)
	})

	handler := middlewares.AppContextMiddleware(ctx)(mux)

	address := "localhost:8080"

	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	ctx.Logger.Info("Listening on address", "addr", address)

	done := make(chan error, 1)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			done <- err
		}
		done <- nil
	}()

	<-ctx.Done()
	ctx.Logger.Info("Shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		ctx.Logger.Error("graceful shutdown failed", "err", err)
		return err
	}

	return <-done
}
