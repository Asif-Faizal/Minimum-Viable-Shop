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
	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL" default:"localhost:50051"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL" default:"localhost:50052"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL" default:"localhost:50053"`
	Port       int    `envconfig:"PORT" default:"8080"`
	Env        string `envconfig:"ENVIRONMENT" default:"development"`
	LogLevel   string `envconfig:"LOG_LEVEL" default:"info"`
}

func main() {
	// Load configuration
	var cfg AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logger := util.NewLogger(cfg.LogLevel)

	// Initialize GraphQL server
	server, err := graphql.NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		logger.Service().Fatal().Err(err).Msg("failed to initialize GraphQL server")
	}
	defer func() {
		if err := server.Close(); err != nil {
			logger.Service().Error().Err(err).Msg("error closing server")
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
		logger.Service().Info().Str("addr", httpServer.Addr).Str("env", cfg.Env).Msg("Starting GraphQL server")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Service().Fatal().Err(err).Msg("server error")
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Service().Info().Msg("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Service().Fatal().Err(err).Msg("shutdown error")
	}

	logger.Service().Info().Msg("server stopped")
}

// healthCheckHandler returns OK if the server is running
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok"}`)
}
