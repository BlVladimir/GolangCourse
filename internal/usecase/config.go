package usecase

import (
	"fmt"
	"os"
	"strconv"
)

type UseCaseConfig struct {
	WorkersCount int
}

const defaultWorkers = 10

func LoadUseCaseConfig() (UseCaseConfig, error) {
	workersCount, err := envInt("WORKERS", defaultWorkers)
	if err != nil {
		return UseCaseConfig{}, err
	}

	validateErr := configValidate(workersCount)
	if validateErr != nil {
		return UseCaseConfig{}, err
	}

	return UseCaseConfig{WorkersCount: workersCount}, nil
}

func envInt(key string, def int) (int, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return def, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s: expected integer, got %q: %w", key, raw, err)
	}
	return v, nil
}

func configValidate(workers int) error {
	if workers <= 0 {
		return fmt.Errorf("WORKERS must be > 0, got %d", workers)
	}

	return nil
}
