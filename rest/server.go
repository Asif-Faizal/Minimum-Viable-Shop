package rest

import (
	"fmt"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/account"
	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
)

type Server struct {
	accountClient *account.AccountClient
	logger        util.Logger
}

func NewServer(accountGrpcURL string, logger util.Logger) (*Server, error) {
	if accountGrpcURL == "" {
		return nil, fmt.Errorf("ACCOUNT_GRPC_URL must be provided")
	}

	accountClient, err := account.NewAccountClient(accountGrpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to account service: %w", err)
	}

	return &Server{
		accountClient: accountClient,
		logger:        logger,
	}, nil
}

func (s *Server) Close() error {
	if s.accountClient != nil {
		s.accountClient.Close()
	}
	return nil
}
