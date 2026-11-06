package grpc_api

import (
	"context"
	"errors"
	"log"

	"github.com/BlVladimir/GolangCourse/internal/film"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, film.ErrInvalidRange):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, film.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, film.ErrUpstream):
		return status.Error(codes.Unavailable, err.Error())
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "client closed request")
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "upstream timeout")
	default:
		log.Printf("grpcapi: internal error: %v", err)
		return status.Error(codes.Internal, "internal error")
	}
}
