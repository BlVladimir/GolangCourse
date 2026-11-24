package http_api

import (
	"github.com/BlVladimir/GolangCourse/internal/model"
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

func oneToMoviesDTO(movie model.Film) movieDTO {
	return movieDTO{ID: movie.ID, Title: movie.Title, Year: movie.Year, Director: movie.Director}
}

func manyToMoviesDTO(movies []model.Film) []movieDTO {
	out := make([]movieDTO, len(movies))
	for idx, movie := range movies {
		out[idx] = oneToMoviesDTO(movie)
	}
	return out
}
