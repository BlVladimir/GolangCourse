package film

import (
	"fmt"
)

func requestValidate(from, to int) error {
	if from <= 0 {
		return fmt.Errorf("%w: from must be > 0, got %d", ErrInvalidRange, from)
	}
	if to <= 0 {
		return fmt.Errorf("%w: to must be > 0, got %d", ErrInvalidRange, to)
	}
	if from > to {
		return fmt.Errorf("%w: from (%d) > to (%d)", ErrInvalidRange, from, to)
	}
	return nil
}
