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
	DatabaseURL string `envconfig:"DATABASE_URL,required"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL,required"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL,required"`
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("⚠️ Warning: .env file not found, relying on system env vars")
	}

	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		AccountURL:  os.Getenv("ACCOUNT_SERVICE_URL"),
		CatalogURL:  os.Getenv("CATALOG_SERVICE_URL"),
	}

	// Validate configuration
	if cfg.DatabaseURL == "" {
		log.Fatal("❌ DATABASE_URL not set")
	}
	if cfg.AccountURL == "" {
		log.Fatal("❌ ACCOUNT_SERVICE_URL not set")
	}
	if cfg.CatalogURL == "" {
		log.Fatal("❌ CATALOG_SERVICE_URL not set")
	}

	log.Println("✅ Using Database URL:", cfg.DatabaseURL)
	log.Println("✅ Using Account Service URL:", cfg.AccountURL)
	log.Println("✅ Using Catalog Service URL:", cfg.CatalogURL)

	// Retry DB connection
	var r order.Repository
	err := retry.Do(
		func() error {
			var err error
			r, err = order.NewPostgresRepository(cfg.DatabaseURL)
			if err != nil {
				log.Printf("❌ Attempt to connect to database failed: %v", err)
			}
			return err
		},
		retry.Delay(5*time.Second),
		retry.Attempts(3),
	)
	if err != nil {
		log.Fatalf("💥 Could not establish database connection after retries: %v", err)
	}
	defer r.Close()

	// Start the gRPC server
	log.Println("🚀 Order Server is Listening at port 8080...")
	s := order.NewService(r)
	log.Fatal(order.ListenGRPC(s, cfg.AccountURL, cfg.CatalogURL, 8080))
}
