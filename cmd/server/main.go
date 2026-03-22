package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"tsuskills-skills/config"
	router "tsuskills-skills/internal/delivery/http"
	"tsuskills-skills/internal/delivery/http/handler"
	"tsuskills-skills/internal/infra/postgres"
	"tsuskills-skills/internal/logger"
	"tsuskills-skills/internal/repository"
	"tsuskills-skills/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger, err := logger.New(&cfg.Logger.Logger)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	appLogger.Info(ctx, "Starting skills service...")

	pool, err := postgres.Connect(ctx, &cfg.Postgres)
	if err != nil {
		appLogger.Fatal(ctx, fmt.Sprintf("Failed to connect to Postgres: %v", err))
	}
	defer pool.Close()
	appLogger.Info(ctx, "Connected to PostgreSQL")

	connString := cfg.Postgres.Pool.ConnConfig.ConnString()
	if err := postgres.RunMigrations(connString, cfg.Postgres.MigrationsPath); err != nil {
		appLogger.Fatal(ctx, fmt.Sprintf("Migrations failed: %v", err))
	}
	appLogger.Info(ctx, "Migrations applied")

	repo := repository.New(pool)
	svc := service.New(repo, appLogger)
	h := handler.NewHandler(svc, appLogger)
	r := router.NewRouter(h, appLogger)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		appLogger.Info(ctx, "Shutting down...")
		shutCtx, shutCancel := context.WithTimeout(context.Background(), cfg.Server.ShutDownTimeOut)
		defer shutCancel()
		httpServer.Shutdown(shutCtx)
		cancel()
	}()

	appLogger.Info(ctx, fmt.Sprintf("Server starting on %s:%d", cfg.Server.Host, cfg.Server.Port))
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		appLogger.Fatal(ctx, fmt.Sprintf("Server failed: %v", err))
	}
	appLogger.Info(ctx, "Server stopped")
}
