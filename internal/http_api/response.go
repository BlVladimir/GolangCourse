package http_api

import "github.com/BlVladimir/GolangCourse/internal/model"

type FilmError struct {
	ID     int    `json:"id"     example:"5"`
	Reason string `json:"reason" example:"status 404"`
}

type LoadMoviesResponse struct {
	Movies []model.Film `json:"movies"`
	Errors []FilmError  `json:"errors"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid range"`
}

func toLoadMoviesResponse(report model.LoadReport) LoadMoviesResponse {
	response := LoadMoviesResponse{
		Movies: make([]model.Film, 0, len(report.FilmResults)),
		Errors: make([]FilmError, 0),
	}
	for _, r := range report.FilmResults {
		if r.Err != nil {
			response.Errors = append(response.Errors, FilmError{ID: r.ID, Reason: r.Err.Error()})
			continue
		}
		response.Movies = append(response.Movies, r.FilmInfo)
	}
	return response
}
