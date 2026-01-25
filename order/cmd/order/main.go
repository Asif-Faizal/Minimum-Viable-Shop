package main

import (
	"log"
	"time"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASE_URL"`
	AccountUrl  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogUrl  string `envconfig:"CATALOG_SERVICE_URL"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	var repository order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = order.NewPostgresRepository(config.DatabaseUrl)
		if err != nil {
			log.Printf("failed to connect to database: %v", err)
			return err
		}
		return nil
	})
	defer repository.Close()
	log.Println("connected to database")
	service := order.NewOrderService(repository)
	log.Fatal(order.ListenGrpcServer(service, config.AccountUrl, config.CatalogUrl, 8080))
	log.Println("order service started")
}
