package main

import (
	"log"
	"time"

	"github.com/Shridhar2104/logilo/shopify"

	"github.com/tinrab/retry"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_SHOPIFY_URL"`
	
}

func main() {

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to process environment variables: %v", err)
	}

	var r shopify.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {

		r, err = shopify.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			return err
		}
		return nil
	})
	defer r.Close()
	log.Println("server starting on port 8080 ...")

	s := shopify.NewShopifyService(r)
	log.Fatal(shopify.NewGRPCServer(s, 8080))
}
