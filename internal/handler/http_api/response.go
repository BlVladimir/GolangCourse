package http_api

import (
	"github.com/BlVladimir/GolangCourse/internal/domain"
)

type movieDTO struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func toMoviesDTO(movies []domain.Movie) []movieDTO {
	out := make([]movieDTO, len(movies))
	for idx, movie := range movies {
		out[idx] = movieDTO{ID: movie.ID, Title: movie.Title, Year: movie.Year, Director: movie.Director}
	}
	return out
}
