package main

import (
	"context"
	"fmt"
	"os"

	"citadel/highway17/internal/app"
	"citadel/highway17/internal/config"
	"citadel/highway17/internal/database"
	"citadel/highway17/internal/logger"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// Connect to database
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Sugar().Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(ctx); err != nil {
		log.Sugar().Fatalf("Failed to run migrations: %v", err)
	}

	log.Sugar().Info("Database migrations completed successfully")

	// Create and start application
	echoApp, err := app.New(cfg, db, log)
	if err != nil {
		log.Sugar().Fatalf("Failed to create application: %v", err)
	}

	listenAddr := fmt.Sprintf(":%d", cfg.AppPort)
	log.Sugar().Infof("Starting Highway 17 Dashboard on %s", listenAddr)

	if err := echoApp.Start(listenAddr); err != nil {
		log.Sugar().Fatalf("Server error: %v", err)
	}
}
