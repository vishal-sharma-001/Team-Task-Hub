package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/launchventures/team-task-hub-backend/internal/app"
	"github.com/launchventures/team-task-hub-backend/internal/config"
)

func main() {
	godotenv.Load()

	cfg := config.New()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	if err := runMigrations(dbURL); err != nil {
		// Log but don't fail - migrations might already be run
		log.Printf("Warning: Migration skipped: %v\n", err)
	}

	log.Println("Initializing application...")
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v\n", err)
	}
	log.Println("Application initialized successfully")
	defer application.Close()

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s\n", addr)
	if err := http.ListenAndServe(addr, application.Router); err != nil {
		log.Fatalf("Server error: %v\n", err)
	}
}

func runMigrations(dbURL string) error {
	// First, verify database connectivity with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer conn.Close(ctx)

	// If connected, run migrations
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return fmt.Errorf("migration instance creation failed: %w", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("getting migration version failed: %w", err)
	}

	log.Printf("Database migration completed. Current version: %d (dirty: %v)\n", version, dirty)
	return nil
}
