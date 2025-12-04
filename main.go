package main

import (
	"context"
	"http-sqlite-template/api"
	"http-sqlite-template/db"
	"http-sqlite-template/middlewares"
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

	// Initialize database
	database, err := db.New("./data/app.db", logger)
	if err != nil {
		logger.Error("failed to initialize database", "err", err)
		return
	}
	defer database.Close()

	// Run migrations
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
