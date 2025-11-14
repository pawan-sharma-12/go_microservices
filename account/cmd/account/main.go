package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/joho/godotenv"
	"github.com/pawan-sharma-12/go_microservices/account"
)
type Config struct{
	DatabaseURL string `envconfig:"DATABASE_URL,required"`
}
func main(){

	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}
	_ = godotenv.Load(".env." + env)

	cfg := Config{
		DatabaseURL: os.Getenv("ACCOUNT_DATABASE_URL"),
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL not set")
	}
	log.Println("Using Database URL:", cfg.DatabaseURL)
	var r account.Repository
	err := retry.Do(
		func() error {
			var err error
			r, err = account.NewPostgresRepository(cfg.DatabaseURL)
			if err != nil {
				log.Printf("❌ Attempt to connect to database failed: %v", err)
			}
			return err
		},
		retry.Delay(5*time.Second),
		retry.Attempts(3), // for debugging, use finite attempts
	)
	
	if err != nil {
		log.Fatalf("💥 Could not establish database connection after retries: %v", err)
	}
	defer r.Close()
	
	portStr := os.Getenv("ACCOUNT_SERVICE_PORT")
	if portStr == "" {
		portStr = "50051" // default for local dev
	}
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}
	
	log.Println("Account Service Listening at port", portInt)
	s := account.NewAccountService(r)
	log.Fatal(account.ListenAndServeGRPC(s, portInt))
	
}