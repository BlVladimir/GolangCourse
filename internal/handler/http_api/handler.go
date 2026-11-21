package http_api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

type MovieService interface {
	LoadMovies(ctx context.Context, from, to int) ([]domain.Movie, error)
	SearchMovies(ctx context.Context, query string) ([]domain.Movie, error)
	GetLastMovie(ctx context.Context) (domain.Movie, error)
}

type Handler struct {
	services MovieService
}

func New(movieServise MovieService) *Handler {
	return &Handler{services: movieServise}
}

// LoadMovies godoc
//
// @Summary      Load movies in a range
// @Description  Downloads movies with ids in [from; to] from the stand.
//
//	Partial per-film failures are reported in the "errors" array
//	alongside HTTP 200 — the request itself succeeded.
//
// @Tags         movies
// @Produce      json
// @Param        from query int true "First movie id"
// @Param        to   query int true "Last movie id"
// @Success      200 {array} movieDTO
// @Failure      400 {object} ErrorResponse
// @Failure 	 500 {object} ErrorResponse
// @Router       /load-movies [get]
func (handler *Handler) LoadMovies(writer http.ResponseWriter, request *http.Request) {
	from, err := parseIntQuery(request, "from")
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	to, err := parseIntQuery(request, "to")
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}

	movies, err := handler.services.LoadMovies(request.Context(), from, to)
	if err != nil {
		status, msg := mapError(err)
		writeError(writer, status, msg)
		return
	}

	writeJSON(writer, http.StatusOK, manyToMoviesDTO(movies))
}

// GetLastMovie godoc
//
// @Summary      Get the last available movie
// @Description  Finds the film with the maximum id that exists on the stand.
// @Tags         movies
// @Produce      json
// @Success      200 {object} movieDTO
// @Failure      404 {object} ErrorResponse
// @Failure 	 500 {object} ErrorResponse
// @Router       /last-movie [get]
func (h *Handler) GetLastMovie(writer http.ResponseWriter, request *http.Request) {
	filmData, err := h.services.GetLastMovie(request.Context())
	if err != nil {
		status, msg := mapError(err)
		writeError(writer, status, msg)
		return
	}
	writeJSON(writer, http.StatusOK, oneToMoviesDTO(filmData))
}

// SearchMovies godoc
//
// @Summary      Search movies
// @Description  Searches movies by query across indexed fields (title, description, director, actors, genres).
// @Description  Results are cached in Redis for 1 hour.
// @Tags         movies
// @Produce      json
// @Param        q query string true "Search query"
// @Success      200 {array}  movieDTO
// @Failure      400 {object} ErrorResponse "Missing or empty q"
// @Failure      500 {object} ErrorResponse
// @Router       /search [get]
func (handler *Handler) SearchMovies(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("q")
	if strings.TrimSpace(query) == "" {
		writeError(writer, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	movies, err := handler.services.SearchMovies(request.Context(), query)
	if err != nil {
		status, msg := mapError(err)
		writeError(writer, status, msg)
		return
	}
	writeJSON(writer, http.StatusOK, manyToMoviesDTO(movies))
}

func parseIntQuery(r *http.Request, name string) (int, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return 0, fmt.Errorf("query parameter %q is required", name)
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("query parameter %q must be an integer", name)
	}
	return v, nil
}
