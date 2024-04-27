package auth

import (
	"encoding/json"
	"io"
	"net/http"

	authmodel "github.com/Heaterdog/AuthServer/internal/model/auth"
	"github.com/Heaterdog/AuthServer/internal/transport"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

func (handler *tokenHandler) registerHandler(w http.ResponseWriter, r *http.Request) {
	handler.logger.Debug("register handler")

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

	err = handler.authService.SignUp(r.Context(), request)

	if err == pgx.ErrNoRows {
		handler.logger.Debug(err.Error())
		transport.ResponseWriteError(w, http.StatusNotFound, err.Error(), handler.logger)

		return
	}

	if err != nil {
		handler.logger.Debug(err.Error())
		transport.ResponseWriteError(w, http.StatusBadRequest, err.Error(), handler.logger)

		return
	}

	w.WriteHeader(http.StatusCreated)
}
