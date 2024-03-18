package response

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorReponse struct {
	ErrorMsg string `json:"error"`
}

type ErrorsReponse struct {
	ErrorsMsg []string `json:"errors"`
}

func JSONError(w http.ResponseWriter, code int, msg string, log *slog.Logger) {
	w.WriteHeader(code)

	b, err := json.Marshal(ErrorReponse{ErrorMsg: msg})
	if err != nil {
		log.Error(err.Error())
	}

	if _, err := w.Write(b); err != nil {
		log.Error(err.Error())
	}
}

func JSONErrors(w http.ResponseWriter, code int, errors []string, log *slog.Logger) {
	w.WriteHeader(code)

	b, err := json.Marshal(ErrorsReponse{ErrorsMsg: errors})
	if err != nil {
		log.Error(err.Error())
	}

	if _, err := w.Write(b); err != nil {
		log.Error(err.Error())
	}
}

func JSON(w http.ResponseWriter, code int, body any, log *slog.Logger) {
	w.WriteHeader(code)
	b, err := json.Marshal(body)
	if err != nil {
		log.Error(err.Error())
	}
	if _, err := w.Write([]byte(b)); err != nil {
		log.Error(err.Error())
	}
}
