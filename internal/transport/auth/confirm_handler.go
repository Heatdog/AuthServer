package auth

import (
	"fmt"
	"net/http"

	"github.com/Heaterdog/AuthServer/internal/transport"
)

func (handler *tokenHandler) confirmHandler(w http.ResponseWriter, r *http.Request) {
	handler.logger.Debug("login handler")

	id := r.URL.Query().Get("id")
	role := r.URL.Query().Get("role")

	if id == "" || role == "" {
		err := fmt.Errorf("null params")

		handler.logger.Debug(err.Error())
		transport.ResponseWriteError(w, http.StatusBadRequest, err.Error(), handler.logger)

		return
	}

	if role != "HeadOfDepartment" || role != "Worker" {
		err := fmt.Errorf("inorrect role")

		handler.logger.Debug(err.Error())
		transport.ResponseWriteError(w, http.StatusBadRequest, err.Error(), handler.logger)

		return
	}

}
