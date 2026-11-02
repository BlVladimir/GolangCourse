package usecase

import (
	"fmt"

	"github.com/BlVladimir/GolangCourse/internal/domain"
)

func requestValidate(from, to int) error {
	if from <= 0 {
		return fmt.Errorf("%w: from must be > 0, got %d", domain.ErrInvalidRange, from)
	}
	if to <= 0 {
		return fmt.Errorf("%w: to must be > 0, got %d", domain.ErrInvalidRange, to)
	}
	if from > to {
		return fmt.Errorf("%w: from (%d) > to (%d)", domain.ErrInvalidRange, from, to)
	}
	return nil
}
