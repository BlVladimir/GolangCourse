package film_service

import (
	"fmt"
	"time"
)

func configValidate(workers int, timeout time.Duration) error {
	if workers <= 0 {
		return fmt.Errorf("--workers must be > 0, got %d", workers)
	}
	if timeout <= 0 {
		return fmt.Errorf("--timeout must be > 0, got %s", timeout)
	}

	return nil
}

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
