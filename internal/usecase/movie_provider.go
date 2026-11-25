package usecase

import (
	"context"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

type MovieProvider interface {
	GetMovie(ctx context.Context, id int) (domain.Movie, error)
}
