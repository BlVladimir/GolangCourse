package film_service

import (
	"fmt"
)

func requestValidate(from int, to int) error {

	if from <= 0 {
		return fmt.Errorf("--from must be > 0, got %d", from)
	}
	if to <= 0 {
		return fmt.Errorf("--to must be > 0, got %d", to)
	}
	if from > to {
		return fmt.Errorf("--from (%d) must not be greater than --to (%d)", from, to)
	}

	return nil
}
