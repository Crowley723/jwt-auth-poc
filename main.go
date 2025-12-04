package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"jwt-auth-poc/api"
	"jwt-auth-poc/crypt_utils"
	"jwt-auth-poc/db"
	"jwt-auth-poc/middlewares"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jwtProvider := ReadOrGenerateJWTKeys(logger)
	if jwtProvider == nil {
		logger.Error("failed to initialize jwt provider")
		return
	}

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

	appCtx := middlewares.NewAppContext(ctx, logger, database, jwtProvider)

	err = api.StartServer(appCtx)
	if err != nil {
		logger.Error("failed to start server", "err", err)
		return
	}
}

func ReadOrGenerateJWTKeys(logger *slog.Logger) crypt_utils.JWTProvider {
	keyFile := crypt_utils.GetJWTPrivateKeyPath()

	var privateKey *ecdsa.PrivateKey
	var err error

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		fmt.Println("Generating private key")
		privateKey, err = crypt_utils.CreateSigningKeys()
		if err != nil {
			logger.Error("failed to create jwt signing key", "err", err)
			return nil
		}
	} else {
		logger.Debug("Loading signing key...")
		privateKey, err = crypt_utils.LoadECDSAPrivateKeyFromPEM()
		if err != nil {
			logger.Error("failed to load private key", "err", err)
			return nil
		}
	}

	jwtProvider, err := crypt_utils.NewECDSAJWTProvider(privateKey)
	if err != nil {
		logger.Error("failed to initialize jwt provider", "err", err)
		return nil
	}

	return jwtProvider
}
