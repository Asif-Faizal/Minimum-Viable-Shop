package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/graphql"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL" default:"localhost:50051"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL" default:"localhost:50052"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL" default:"localhost:50053"`
	Port       int    `envconfig:"PORT" default:"8080"`
	Env        string `envconfig:"ENVIRONMENT" default:"development"`
}

func main() {
	// Load configuration
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize GraphQL server
	server, err := graphql.NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatalf("failed to initialize GraphQL server: %v", err)
	}
	defer func() {
		if err := server.Close(); err != nil {
			log.Printf("error closing server: %v", err)
		}
	}()

	// Create HTTP server
	mux := http.NewServeMux()

	// GraphQL endpoint
	mux.Handle("/graphql", handler.NewDefaultServer(server.ToExecutableSchema()))

	// Playground endpoint
	mux.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	// Health check endpoint
	mux.HandleFunc("/health", healthCheckHandler)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting GraphQL server on %s (environment: %s)", httpServer.Addr, cfg.Env)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}

	log.Println("server stopped")
}

// healthCheckHandler returns OK if the server is running
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok"}`)
}
