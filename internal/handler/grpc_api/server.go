package grpc_api

import (
	"context"

	moviepb "github.com/BlVladimir/GolangCourse/api/proto"
	"github.com/BlVladimir/GolangCourse/internal/domain"
)

type FilmService interface {
	LoadMovies(ctx context.Context, from, to int) ([]domain.Movie, error)
	GetLastMovie(ctx context.Context) (domain.Movie, error)
}

type Server struct {
	moviepb.UnimplementedMovieServiceServer
	service FilmService
}

func New(service FilmService) *Server {
	return &Server{service: service}
}

func (server *Server) LoadMovies(ctx context.Context, request *moviepb.LoadMoviesRequest) (*moviepb.LoadMoviesResponse, error) {
	from := int(request.GetFrom())
	to := int(request.GetTo())

	movies, err := server.service.LoadMovies(ctx, from, to)
	if err != nil {
		return nil, mapError(err)
	}

	return toProtoResponse(movies), nil
}

func (server *Server) GetLastMovie(ctx context.Context, _ *moviepb.GetLastMovieRequest) (*moviepb.Movie, error) {
	filmData, err := server.service.GetLastMovie(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoFilm(filmData), nil
}

var _ moviepb.MovieServiceServer = (*Server)(nil)
