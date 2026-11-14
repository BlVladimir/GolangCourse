package grpc_api

import (
	"context"

	moviepb "github.com/BlVladimir/GolangCourse/api/proto"
	"github.com/BlVladimir/GolangCourse/internal/domain"
	// "github.com/BlVladimir/GolangCourse/internal/usecase"
)

type FilmService interface {
	LoadMovies(ctx context.Context, from, to int) ([]domain.Movie, error)
	GetLastMovie(ctx context.Context) (domain.Movie, error)
}

type Server struct {
	moviepb.UnimplementedMovieServiceServer
	svc FilmService
}

func New(svc FilmService) *Server {
	return &Server{svc: svc}
}

func (s *Server) LoadMovies(ctx context.Context, req *moviepb.LoadMoviesRequest) ([]domain.Movie, error) {
	from := int(req.GetFrom())
	to := int(req.GetTo())

	movies, err := s.svc.LoadMovies(ctx, from, to)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (s *Server) GetLastMovie(ctx context.Context, _ *moviepb.GetLastMovieRequest) (*moviepb.Film, error) {
	filmData, err := s.svc.GetLastMovie(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoFilm(filmData), nil
}

var _ moviepb.MovieServiceServer = (*Server)(nil)
