// Command server é o binário único da RoosterLabs: serve a landing hoje
// e o produto (motor de extração, área do cliente) amanhã.
//
// É um servidor HTTP comum. Em produção roda dentro do AWS Lambda via
// Lambda Web Adapter, mas o código não sabe disso — e não deve saber.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/roosterlabs/roosterlabs-engineering/internal/leads"
	"github.com/roosterlabs/roosterlabs-engineering/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default do Lambda Web Adapter
	}

	ctx, cancelConnect := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancelConnect()

	var store leads.Store
	var err error
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL != "" {
		store, err = leads.NewPostgresStore(ctx, databaseURL)
		if err != nil {
			logger.Error("failed to connect to postgres", "err", err)
			os.Exit(1)
		}
		defer func() {
			if err := store.Close(); err != nil {
				logger.Error("failed to close store", "err", err)
			}
		}()
	} else {
		logger.Warn("DATABASE_URL not set, using memory store")
		store = leads.NewMemoryStore()
	}

	handler, err := server.NewHandler(server.Config{
		Store:        store,
		BaseURL:      os.Getenv("BASE_URL"),
		ContactEmail: os.Getenv("CONTACT_EMAIL"),
	})
	if err != nil {
		logger.Error("failed to build handler", "err", err)
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "err", err)
	}
}
