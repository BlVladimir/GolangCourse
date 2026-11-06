package http_api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/BlVladimir/GolangCourse/internal/film"
)

const statusClientClosedRequest = 499

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("httpapi: encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, film.ErrInvalidRange):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, film.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, film.ErrUpstream):
		return http.StatusBadGateway, err.Error()
	case errors.Is(err, context.Canceled):
		return statusClientClosedRequest, "client closed request"
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout, "upstream timeout"
	default:
		log.Printf("httpapi: internal error: %v", err)
		return http.StatusInternalServerError, "internal error"
	}
}
