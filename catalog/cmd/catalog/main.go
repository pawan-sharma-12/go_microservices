package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/joho/godotenv"
	"github.com/pawan-sharma-12/go_microservices/catalog"
)

type Config struct {
	ElasticURL string
}

func main() {
	// Determine environment
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	// Load environment-specific .env file
	_ = godotenv.Load(".env." + env)

	cfg := Config{
		ElasticURL: os.Getenv("CATALOG_ELASTIC_URL"),
	}
	if cfg.ElasticURL == "" {
		log.Fatal("CATALOG_ELASTIC_URL not set")
	}

	log.Println("Using Elastic URL:", cfg.ElasticURL)

	// Connect to Elasticsearch with retries
	var r catalog.Repository
	err := retry.Do(
		func() error {
			var err error
			r, err = catalog.NewElasticRepository(cfg.ElasticURL)
			if err != nil {
				log.Printf("❌ Attempt to connect to Elasticsearch failed: %v", err)
			}
			return err
		},
		retry.Delay(2*time.Second),
		retry.Attempts(3),
	)

	if err != nil {
		log.Fatalf("💥 Could not establish Elasticsearch connection after retries: %v", err)
	}
	defer r.Close()

	// Determine port
	portStr := os.Getenv("CATALOG_SERVICE_PORT")
	if portStr == "" {
		portStr = "50052" // default local catalog service port
	}

	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	log.Println("Catalog Service Listening at port", portInt)

	s := catalog.NewService(r)
	log.Fatal(catalog.ListenAndServeGRPC(s, portInt))
}
