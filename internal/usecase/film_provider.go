package usecase

import (
	"context"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

type FilmProvider interface {
	GetMovie(ctx context.Context, id int) (domain.Film, error)
}
