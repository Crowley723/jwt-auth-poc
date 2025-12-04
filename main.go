package main

import (
	"context"
	"jwt-auth-poc/api"
	"jwt-auth-poc/db"
	"jwt-auth-poc/middlewares"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.New("./app/data/app.db", logger)
	if err != nil {
		logger.Error("failed to initialize database", "err", err)
		return
	}
	defer database.Close()

	if err := database.RunMigrations(); err != nil {
		logger.Error("failed to run migrations", "err", err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal")
		cancel()
	}()

	appCtx := middlewares.NewAppContext(ctx, logger, database)

	err = api.StartServer(appCtx)
	if err != nil {
		logger.Error("failed to start server", "err", err)
		return
	}
}
