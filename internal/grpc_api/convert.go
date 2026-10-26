package grpc_api

import (
	moviepb "github.com/BlVladimir/GolangCourse/api/proto"
	"github.com/BlVladimir/GolangCourse/internal/model"
)

func toProtoFilm(f model.Film) *moviepb.Movie {
	return &moviepb.Movie{
		Id:       int32(f.ID),
		Title:    f.Title,
		Year:     int32(f.Year),
		Director: f.Director,
	}
}

func toProtoResponse(report model.LoadReport) *moviepb.LoadMoviesResponse {
	resp := &moviepb.LoadMoviesResponse{
		Movies: make([]*moviepb.Movie, 0, len(report.FilmResults)),
		Errors: make([]*moviepb.FilmError, 0),
	}
	for _, r := range report.FilmResults {
		if r.Err != nil {
			resp.Errors = append(resp.Errors, &moviepb.FilmError{
				Id:     int32(r.ID),
				Reason: r.Err.Error(),
			})
			continue
		}
		resp.Movies = append(resp.Movies, toProtoFilm(r.FilmInfo))
	}
	return resp
}
