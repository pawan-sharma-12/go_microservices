package main

import (
	"log"
	"os"
	"time"
	"github.com/avast/retry-go"
	"github.com/joho/godotenv"
	"github.com/pawan-sharma-12/go_microservices/catalog"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL,required"`
}
func main() {
	var config Config
	_ = godotenv.Load("../../../.env")
	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}
	log.Println("Using Database URL:", cfg.DatabaseURL)
	var r catalog.Repository 
	err := retry.Do(
		func() error {
			var err error
			r, err = catalog.NewElasticRepository(config.DatabaseURL)
			if err != nil {
				log.Printf("❌ Attempt to connect to database failed: %v", err)
			}
			return err
		},
		retry.Delay(2*time.Second),
		retry.Attempts(3), // 0 means infinite retries
	)
	if err != nil {
		log.Fatalf("💥 Could not establish database connection after retries: %v", err)
	}
	
	defer r.Close()
	log.Println("Server is Listening at port 8080...")
	s := catalog.NewService(r)
	log.Fatal(catalog.ListenAndServeGRPC(s, 8080))

}
