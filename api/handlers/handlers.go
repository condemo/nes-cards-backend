package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ApiError struct {
	Err    error  `json:"-"`
	Msg    string `json:"error"`
	Status int    `json:"code"`
}

func (e ApiError) Error() string {
	return e.Err.Error()
}

func NewApiError(err error, msg string, status int) ApiError {
	return ApiError{
		Err:    err,
		Msg:    msg,
		Status: status,
	}
}

type CustomHandler func(w http.ResponseWriter, r *http.Request) error

var internalError = map[string]any{
	"status": http.StatusInternalServerError,
	"msg":    "internal server error",
}

func MakeHandler(f CustomHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if apiErr, ok := err.(ApiError); ok {
				SendJSON(w, apiErr.Status, apiErr)
			} else {
				SendJSON(w, http.StatusInternalServerError, internalError)
			}
			slog.Error("API ERROR", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func SendJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
