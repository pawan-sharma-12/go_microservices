package main

import (
	"net/http"
	"log"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/kelseyhightower/envconfig"
	"github.com/99designs/gqlgen/graphql/playground"
)
type AppConfig struct{
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	ProductURL string `envconfig:"PRODUCT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
}
func main(){
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	server, err := NewGraphQlServer(cfg.AccountURL, cfg.CatalogURL, cfg.ProductURL)
	if err != nil {	
		log.Fatal(err)
	}
	http.Handle("/graphql", handler.New(server.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("Siddharth", "/graphql"))
	log.Fatal(http.ListenAndServe(":8080", nil))
}