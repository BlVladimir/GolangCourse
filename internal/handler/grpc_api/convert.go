package grpc_api

import (
	moviepb "github.com/BlVladimir/GolangCourse/api/proto"
	"github.com/BlVladimir/GolangCourse/internal/domain"
	// "github.com/BlVladimir/GolangCourse/internal/usecase"
)

func toProtoFilm(film domain.Movie) *moviepb.Film {
	return &moviepb.Film{
		Id:       int32(film.ID),
		Title:    film.Title,
		Year:     int32(film.Year),
		Director: film.Director,
	}
}

// func toProtoResponse(report usecase.LoadReport) *moviepb.LoadFilmsResponse {
// 	resp := &moviepb.LoadFilmsResponse{
// 		Movies: make([]*moviepb.Film, 0, len(report.MovieResults)),
// 		Errors: make([]*moviepb.FilmError, 0),
// 	}
// 	for _, r := range report.MovieResults {
// 		if r.Err != nil {
// 			resp.Errors = append(resp.Errors, &moviepb.FilmError{
// 				Id:     int32(r.ID),
// 				Reason: r.Err.Error(),
// 			})
// 			continue
// 		}
// 		resp.Movies = append(resp.Movies, toProtoFilm(r.MovieInfo))
// 	}
// 	return resp
// }
