package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"jwt-auth-poc/middlewares"
)

func StartServer(ctx *middlewares.AppContext) error {
	mux := http.NewServeMux()

	RegisterRoutes(mux, ctx)

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
