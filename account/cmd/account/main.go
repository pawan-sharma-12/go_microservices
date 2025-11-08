package main

import (
	"log"
	"time"
	"github.com/kelseyhightower/envconfig"
	"github.com/pawan-sharma-12/go_microservices/account"
	"github.com/avast/retry-go"
	"github.com/joho/godotenv"
)
type Config struct{
	DatabaseURL string `env:"DATABASE_URL,required"`
}
func main(){
	_ = godotenv.Load()
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	var r account.Repository
	err = retry.Do(
		func() error {
			var err error
			r, err = account.NewPostgresRepository(cfg.DatabaseURL)
			if err != nil {
				log.Printf("Failed to connect to database: %v", err)
			}
			return err
		},
		retry.Delay(5*time.Second),
		retry.Attempts(0), // Infinite retries
	)
	if err != nil {
		log.Fatal("Could not establish database connection: ", err)
	}
	defer r.Close()
	log.Println("Server is Listening at port 8080...")
	s := account.NewAccountService(r)
	log.Fatal(account.ListenAndServeGRPC(s, 8080 ))
}