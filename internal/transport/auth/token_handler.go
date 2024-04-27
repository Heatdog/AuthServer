package auth

import (
	"log/slog"
	"net/http"

	authservice "github.com/Heaterdog/AuthServer/internal/service/auth"
	"github.com/Heaterdog/AuthServer/internal/transport"
	"github.com/gorilla/mux"
)

type tokenHandler struct {
	logger      *slog.Logger
	authService authservice.AuthService
}

func NewTokenHandler(logger *slog.Logger, authService authservice.AuthService) transport.Handler {
	return &tokenHandler{
		logger:      logger,
		authService: authService,
	}
}

const (
	login          = "/login"
	register       = "/register"
	changePassword = "/pswd"
	confirm        = "/confirm/{id}/{role}"
)

func (handler *tokenHandler) Register(router *mux.Router) {
	router.HandleFunc(login, handler.loginHandler).Methods(http.MethodPost)
	router.HandleFunc(register, handler.registerHandler).Methods(http.MethodPost)
	router.HandleFunc(confirm, handler.confirmHandler).Methods(http.MethodPost)
}
