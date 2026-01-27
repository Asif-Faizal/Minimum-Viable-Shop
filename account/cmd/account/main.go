package main

import (
	"log"
	"time"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/account"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseUrl        string        `envconfig:"DATABASE_URL"`
	Port               int           `envconfig:"GRPC_PORT" default:"8080"`
	LogLevel           string        `envconfig:"LOG_LEVEL" default:"info"`
	JwtSecret          string        `envconfig:"JWT_SECRET" default:"my-secret-key"`
	AccessTokenExpiry  time.Duration `envconfig:"ACCESS_TOKEN_EXPIRY" default:"45m"`
	RefreshTokenExpiry time.Duration `envconfig:"REFRESH_TOKEN_EXPIRY" default:"168h"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	logger := util.NewLogger(config.LogLevel)

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

	service := account.NewAccountService(
		repository,
		config.JwtSecret,
		config.AccessTokenExpiry,
		config.RefreshTokenExpiry,
	)

	// Start gRPC server (blocks)
	logger.Service().Info().Int("port", config.Port).Msg("starting account service")

	if err := account.ListenGrpcServer(service, logger, config.Port); err != nil {
		logger.Service().Fatal().Err(err).Msg("failed to start gRPC server")
	}
}
