package main

import (
	"log"
	"time"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASE_URL"`
	Port        int    `envconfig:"GRPC_PORT" default:"8080"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	var repository catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = catalog.NewElasticRepository(config.DatabaseUrl)
		if err != nil {
			log.Printf("failed to connect to database: %v", err)
			return err
		}
		return nil
	})
	defer repository.Close()
	log.Println("connected to database")
	service := catalog.NewCatalogService(repository)
	log.Fatal(catalog.ListenGrpcServer(service, config.Port))
	log.Println("catalog service started")
}
