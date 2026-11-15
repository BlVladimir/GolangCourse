package http_api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

type MovieService interface {
	LoadMovies(ctx context.Context, from, to int) ([]domain.Movie, error)
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
func (handler *Handler) LoadMovies(w http.ResponseWriter, r *http.Request) {
	from, err := parseIntQuery(r, "from")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	to, err := parseIntQuery(r, "to")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	movies, err := handler.services.LoadMovies(r.Context(), from, to)
	if err != nil {
		status, msg := mapError(err)
		writeError(w, status, msg)
		return
	}

	writeJSON(w, http.StatusOK, manyToMoviesDTO(movies))
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
func (h *Handler) GetLastMovie(w http.ResponseWriter, r *http.Request) {
	filmData, err := h.services.GetLastMovie(r.Context())
	if err != nil {
		status, msg := mapError(err)
		writeError(w, status, msg)
		return
	}
	writeJSON(w, http.StatusOK, oneToMoviesDTO(filmData))
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
