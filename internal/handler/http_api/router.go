package http_api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/BlVladimir/GolangCourse/api/swagger"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /load-movies", h.LoadMovies)
	mux.HandleFunc("GET /last-movie", h.GetLastMovie)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
