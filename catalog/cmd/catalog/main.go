package main

import (
	"log"
	"time"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/catalog"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseUrl string `envconfig:"ELASTICSEARCH_URL"`
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
	var repository catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = catalog.NewElasticRepository(config.DatabaseUrl, logger)
		if err != nil {
			logger.Service().Error().Err(err).Msg("failed to connect to database")
			return err
		}
		return nil
	})
	defer repository.Close()
	logger.Service().Info().Msg("connected to database")
	service := catalog.NewCatalogService(repository)
	logger.Service().Info().Int("port", config.Port).Msg("starting catalog service")
	logger.Service().Fatal().Err(catalog.ListenGrpcServer(service, logger, config.Port)).Msg("failed to start gRPC server")
}
