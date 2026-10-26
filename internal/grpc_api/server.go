package grpc_api

import (
	"context"

	moviepb "github.com/BlVladimir/GolangCourse/api/proto"
	"github.com/BlVladimir/GolangCourse/internal/model"
)

type FilmService interface {
	LoadMovies(ctx context.Context, from, to int) (model.LoadReport, error)
	GetLastMovie(ctx context.Context) (model.Film, error)
}

type Server struct {
	moviepb.UnimplementedMovieServiceServer
	svc FilmService
}

func New(svc FilmService) *Server {
	return &Server{svc: svc}
}

func (s *Server) LoadMovies(ctx context.Context, req *moviepb.LoadMoviesRequest) (*moviepb.LoadMoviesResponse, error) {
	from := int(req.GetFrom())
	to := int(req.GetTo())

	report, err := s.svc.LoadMovies(ctx, from, to)
	if err != nil {
		return nil, mapError(err)
	}

	return toProtoResponse(report), nil
}

func (s *Server) GetLastMovie(ctx context.Context, _ *moviepb.GetLastMovieRequest) (*moviepb.Movie, error) {
	filmData, err := s.svc.GetLastMovie(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoFilm(filmData), nil
}

var _ moviepb.MovieServiceServer = (*Server)(nil)
