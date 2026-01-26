package account

import (
	"net/http"

	"github.com/Asif-Faizal/Minimum-Viable-Shop/util"
)

func NewRestHandler(service Service, logger util.Logger) http.Handler {
	server := &restServer{
		service: service,
		logger:  logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", server.handleHealth)
	mux.HandleFunc("/accounts/check-email", server.handleCheckEmail)
	mux.HandleFunc("/accounts/login", server.handleLogin)
	mux.HandleFunc("/accounts/logout", server.handleLogout)
	mux.HandleFunc("/accounts/refresh", server.handleRefreshToken)

	return mux
}
