package main

import (
	"log"
	"time"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/order"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASE_URL"`
	AccountUrl  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogUrl  string `envconfig:"CATALOG_SERVICE_URL"`
	Port        int    `envconfig:"GRPC_PORT" default:"8080"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	logger := util.NewLogger(config.LogLevel)
	var repository order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = order.NewPostgresRepository(config.DatabaseUrl, logger)
		if err != nil {
			logger.Service().Error().Err(err).Msg("failed to connect to database")
			return err
		}
		return nil
	})
	defer repository.Close()
	logger.Service().Info().Msg("connected to database")
	service := order.NewOrderService(repository)
	logger.Service().Info().Int("port", config.Port).Msg("starting order service")
	log.Fatal(order.ListenGrpcServer(service, config.AccountUrl, config.CatalogUrl, logger, config.Port))
}
