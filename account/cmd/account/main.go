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
	RestPort    int    `envconfig:"REST_PORT" default:"8081"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"info"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	logger := account.NewLogger(config.LogLevel)

	var repository account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = account.NewPostgresRepository(config.DatabaseUrl, logger)
		if err != nil {
			logger.Service().Error().Err(err).Msg("failed to connect to database")
			return err
		}
		return nil
	})
	defer repository.Close()

	logger.Service().Info().Msg("connected to database")

	service := account.NewAccountService(repository)

	// Start REST server for health check in a goroutine
	go func() {
		logger.Service().Info().Int("port", config.RestPort).Msg("starting REST server")
		if err := account.ListenRestServer(service, logger, config.RestPort); err != nil {
			logger.Service().Fatal().Err(err).Msg("failed to start REST server")
		}
	}()

	// Start gRPC server (blocks)
	logger.Service().Info().Int("port", config.Port).Msg("starting account service")

	if err := account.ListenGrpcServer(service, logger, config.Port); err != nil {
		logger.Service().Fatal().Err(err).Msg("failed to start gRPC server")
	}
}
