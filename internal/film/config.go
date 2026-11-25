package film

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	MovieBaseURL string
	HTTPAddr     string
	GRPCAddr     string
	Workers      int
	Timeout      time.Duration
}

const (
	defaultMovieBaseURL = "https://homeworksite.site"
	defaultHTTPAddr     = ":8080"
	defaultGRPCAddr     = ":9090"
	defaultWorkers      = 10
	defaultTimeout      = 5 * time.Second
)

func Load() (Config, error) {
	workers, err := envInt("WORKERS", defaultWorkers)
	if err != nil {
		return Config{}, err
	}

	timeout, err := envDuration("TIMEOUT", defaultTimeout)
	if err != nil {
		return Config{}, err
	}

	cnf := Config{
		MovieBaseURL: envString("MOVIE_BASE_URL", defaultMovieBaseURL),
		HTTPAddr:     envString("HTTP_ADDR", defaultHTTPAddr),
		GRPCAddr:     envString("GRPC_ADDR", defaultGRPCAddr),
		Workers:      workers,
		Timeout:      timeout,
	}

	configValidationErr := configValidate(cnf.Workers, cnf.Timeout)
	if configValidationErr != nil {
		return Config{}, configValidationErr
	}

	return cnf, nil
}

func envString(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
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

func envDuration(key string, def time.Duration) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return def, nil
	}
	v, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s: expected duration (e.g. 5s), got %q: %w", key, raw, err)
	}
	return v, nil
}

func configValidate(workers int, timeout time.Duration) error {
	if workers <= 0 {
		return fmt.Errorf("WORKERS must be > 0, got %d", workers)
	}
	if timeout <= 0 {
		return fmt.Errorf("TIMEOUT must be > 0, got %s", timeout)
	}

	return nil
}
