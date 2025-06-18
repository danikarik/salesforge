package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/danikarik/salesforge/internal/app"
	"github.com/danikarik/salesforge/internal/model/pg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	ctx := context.Background()

	// Read environment variables into the Specification struct
	var spec app.Specification
	err := envconfig.Process("api", &spec)
	if err != nil {
		log.Fatalf("Failed to process environment variables: %v", err)
	}

	// Connect to the database using the provided URL
	pool, err := pgxpool.New(ctx, spec.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer pool.Close()

	// Create a new store instance
	store, err := pg.NewStore(pool)
	if err != nil {
		log.Fatalf("Failed to create store: %v", err)
	}

	// Create a new service instance with the store
	srv := app.NewService(app.Config{
		Store: store,
	})

	// Create an HTTP server with the service's handler
	httpServer := &http.Server{
		Addr:    spec.Address,
		Handler: srv.Handler(),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting service on ", spec.Address)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down...")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
