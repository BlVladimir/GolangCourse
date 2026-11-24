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

func toProtoResponse(movies []model.Film) *moviepb.LoadMoviesResponse {
	resp := &moviepb.LoadMoviesResponse{
		Movies: make([]*moviepb.Movie, 0, len(movies)),
	}
	for _, movie := range movies {
		resp.Movies = append(resp.Movies, toProtoFilm(movie))
	}
	return resp
}
