package http_api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/BlVladimir/GolangCourse/api/swagger"
)

func NewRouter(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /load-movies", handler.LoadMovies)
	mux.HandleFunc("GET /search", handler.SearchMovies)
	mux.HandleFunc("GET /last-movie", handler.GetLastMovie)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
