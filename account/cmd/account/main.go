package main

import (
	"log"
	"time"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/account"
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
	var repository account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = account.NewPostgresRepository(config.DatabaseUrl)
		if err != nil {
			log.Printf("failed to connect to database: %v", err)
			return err
		}
		return nil
	})
	defer repository.Close()
	log.Println("connected to database")
	service := account.NewAccountService(repository)
	log.Fatal(account.ListenGrpcServer(service, config.Port))
	log.Println("account service started")
}
