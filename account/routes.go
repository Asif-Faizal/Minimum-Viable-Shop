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
	mux.HandleFunc("/health", server.HealthCheck)
	mux.HandleFunc("/accounts/check-email", server.HandleCheckEmail)
	mux.HandleFunc("/accounts/login", server.HandleLogin)
	mux.HandleFunc("/accounts/logout", server.HandleLogout)
	mux.HandleFunc("/accounts/refresh", server.HandleRefreshToken)

	return mux
}
