package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	// "github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	// "github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	// Determine environment
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	// Load environment-specific .env file
	envFile := "../.env." + env
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("⚠️ Warning: %s not found, falling back to system env vars", envFile)
	} else {
		log.Printf("📦 Loaded environment config from %s", envFile)
	}

	// Load config from environment
	cfg := AppConfig{
		AccountURL: os.Getenv("ACCOUNT_SERVICE_URL"),
		CatalogURL: os.Getenv("CATALOG_SERVICE_URL"),
		OrderURL:   os.Getenv("ORDER_SERVICE_URL"),
	}
	log.Println("URL print : ",cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	// Validate
	if cfg.AccountURL == "" || cfg.CatalogURL == "" || cfg.OrderURL == "" {
		log.Fatalf("❌ Missing required environment variables: ACCOUNT_SERVICE_URL, CATALOG_SERVICE_URL, ORDER_SERVICE_URL")
	}

	// Create GraphQL server
	server, err := NewGraphQlServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create gqlgen handler with introspection enabled
	srv := handler.NewDefaultServer(server.ToExecutableSchema())

	// HTTP handlers
	http.Handle("/graphql", srv)
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	log.Println("GraphQL server running on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
