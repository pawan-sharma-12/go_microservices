package main

import (
	"log"
	"os"
	"time"

	"github.com/avast/retry-go"
	"github.com/joho/godotenv"
	"github.com/pawan-sharma-12/go_microservices/order"
)

type Config struct {
	OrderDatabaseURL string `envconfig:"ORDER_DATABASE_URL,required"`
	AccountURL       string `envconfig:"ACCOUNT_SERVICE_URL,required"`
	CatalogURL       string `envconfig:"CATALOG_SERVICE_URL,required"`
}

func main() {

	// -------------------------------
	// Load environment-specific config
	// -------------------------------
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("⚠️ Warning: %s not found, falling back to system env vars", envFile)
	} else {
		log.Printf("📦 Loaded environment config from %s", envFile)
	}

	// -------------------------------
	// Read configuration
	// -------------------------------
	cfg := Config{
		OrderDatabaseURL: os.Getenv("ORDER_DATABASE_URL"),
		AccountURL:       os.Getenv("ACCOUNT_SERVICE_URL"),
		CatalogURL:       os.Getenv("CATALOG_SERVICE_URL"),
	}

	// -------------------------------
	// Validations
	// -------------------------------
	if cfg.OrderDatabaseURL == "" {
		log.Fatal("❌ ORDER_DATABASE_URL not set")
	}
	if cfg.AccountURL == "" {
		log.Fatal("❌ ACCOUNT_SERVICE_URL not set")
	}
	if cfg.CatalogURL == "" {
		log.Fatal("❌ CATALOG_SERVICE_URL not set")
	}

	log.Println("🔗 ORDER_DATABASE_URL:", cfg.OrderDatabaseURL)
	log.Println("🔗 ACCOUNT_SERVICE_URL:", cfg.AccountURL)
	log.Println("🔗 CATALOG_SERVICE_URL:", cfg.CatalogURL)

	// -------------------------------
	// Retry DB connection
	// -------------------------------
	var r order.Repository
	err := retry.Do(
		func() error {
			var err error
			r, err = order.NewPostgresRepository(cfg.OrderDatabaseURL)
			if err != nil {
				log.Printf("❌ DB connection failed: %v", err)
			}
			return err
		},
		retry.Delay(5*time.Second),
		retry.Attempts(3),
	)

	if err != nil {
		log.Fatalf("💥 Could not connect to Order DB after retries: %v", err)
	}
	defer r.Close()

	// -------------------------------
	// Start gRPC server
	// -------------------------------
	log.Println("🚀 Order Service listening on port 8080...")
	s := order.NewService(r)
	log.Fatal(order.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, 8080))
}
