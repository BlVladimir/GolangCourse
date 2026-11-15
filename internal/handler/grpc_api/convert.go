package grpc_api

import (
	moviepb "github.com/BlVladimir/GolangCourse/api/proto"
	"github.com/BlVladimir/GolangCourse/internal/domain"
)

func toProtoFilm(film domain.Movie) *moviepb.Movie {
	return &moviepb.Movie{
		Id:       int32(film.ID),
		Title:    film.Title,
		Year:     int32(film.Year),
		Director: film.Director,
	}
}

func toProtoResponse(movies []domain.Movie) *moviepb.LoadMoviesResponse {
	response := &moviepb.LoadMoviesResponse{
		Movies: make([]*moviepb.Movie, 0, len(movies)),
	}

	for _, movie := range movies {
		response.Movies = append(response.Movies, toProtoFilm(movie))
	}

	return response
}
