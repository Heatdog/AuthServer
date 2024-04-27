package auth

import (
	"encoding/json"
	"io"
	"net/http"

	authmodel "github.com/Heaterdog/AuthServer/internal/model/auth"
	"github.com/Heaterdog/AuthServer/internal/transport"
	"github.com/go-playground/validator/v10"
)

func (handler *tokenHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	handler.logger.Debug("login handler")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		handler.logger.Debug(err.Error())
		transport.ResponseWriteError(w, http.StatusBadRequest, err.Error(), handler.logger)

		return
	}

	defer r.Body.Close()

	handler.logger.Debug("unmarshaling request body")

	var request authmodel.AuthRequest

	if err := json.Unmarshal(body, &request); err != nil {
		handler.logger.Debug(err.Error())
		transport.ResponseWriteError(w, http.StatusBadRequest, err.Error(), handler.logger)

		return
	}

	validate := validator.New()

	if err = validate.Struct(request); err != nil {
		handler.logger.Debug(err.Error())
		transport.ResponseWriteError(w, http.StatusBadRequest, err.Error(), handler.logger)

		return
	}

	access, refresh, err := handler.authService.SignIn(r.Context(), request)
	if err != nil {
		handler.logger.Debug(err.Error())
		transport.ResponseWriteError(w, http.StatusBadRequest, err.Error(), handler.logger)

		return
	}

	resp, err := json.Marshal(authmodel.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
	if err != nil {
		handler.logger.Debug(err.Error())
		transport.ResponseWriteError(w, http.StatusInternalServerError, err.Error(), handler.logger)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")

	if _, err = w.Write(resp); err != nil {
		handler.logger.Warn(err.Error())
		transport.ResponseWriteError(w, http.StatusInternalServerError, err.Error(), handler.logger)

		return
	}

	handler.logger.Debug(string(resp))
}
