package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AccountServiceURL   string        `envconfig:"ACCOUNT_SERVICE_URL" required:"true"`
	GraphQLGatewayURL   string        `envconfig:"GRAPHQL_GATEWAY_URL" required:"true"`
	MaxIdleConns        int           `envconfig:"PROXY_MAX_IDLE_CONNS" default:"100"`
	MaxIdleConnsPerHost int           `envconfig:"PROXY_MAX_IDLE_CONNS_PER_HOST" default:"100"`
	IdleConnTimeout     time.Duration `envconfig:"PROXY_IDLE_CONN_TIMEOUT" default:"90s"`
	RequestTimeout      time.Duration `envconfig:"PROXY_REQUEST_TIMEOUT" default:"30s"`
	Port                int           `envconfig:"PROXY_PORT" default:"80"`
}

func main() {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("[PROXY] Failed to process config: %v", err)
	}

	log.Printf("[PROXY] Starting with config: AccountServiceURL=%s, GraphQLGatewayURL=%s, MaxIdleConns=%d, MaxIdleConnsPerHost=%d, IdleConnTimeout=%v, RequestTimeout=%v",
		cfg.AccountServiceURL, cfg.GraphQLGatewayURL, cfg.MaxIdleConns, cfg.MaxIdleConnsPerHost, cfg.IdleConnTimeout, cfg.RequestTimeout)

	// Parse upstream targets
	authTarget, err := url.Parse(cfg.AccountServiceURL)
	if err != nil {
		log.Fatalf("[PROXY] Failed to parse account URL: %v", err)
	}

	gqlTarget, err := url.Parse(cfg.GraphQLGatewayURL)
	if err != nil {
		log.Fatalf("[PROXY] Failed to parse graphql URL: %v", err)
	}

	// Configure transport with connection pooling
	transport := &http.Transport{
		MaxIdleConns:        cfg.MaxIdleConns,
		MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
		IdleConnTimeout:     cfg.IdleConnTimeout,
	}

	// Create reverse proxies
	authProxy := httputil.NewSingleHostReverseProxy(authTarget)
	authProxy.Transport = transport
	authProxy.ErrorHandler = errorHandler

	gqlProxy := httputil.NewSingleHostReverseProxy(gqlTarget)
	gqlProxy.Transport = transport
	gqlProxy.ErrorHandler = errorHandler

	// Route handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log incoming request
		log.Printf("[PROXY] %s %s", r.Method, r.URL.Path)

		// Route based on path
		if strings.HasPrefix(r.URL.Path, "/auth") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, "/auth")
			authProxy.ServeHTTP(w, r)
		} else {
			// Everything else goes to GraphQL
			gqlProxy.ServeHTTP(w, r)
		}
	})

	// Add timeout wrapper
	timeoutHandler := http.TimeoutHandler(handler, cfg.RequestTimeout, "Request timeout")

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("[PROXY] Starting reverse proxy on %s", addr)
	log.Fatal(http.ListenAndServe(addr, timeoutHandler))
}

func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[PROXY ERROR] %s %s: %v", r.Method, r.URL.Path, err)
	http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
}
