package filmstand

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	MovieBaseURL string
	Timeout      time.Duration
}

const (
	defaultMovieBaseURL = "https://homeworksite.site"
	defaultTimeout      = 5 * time.Second
)

func LoadClientConfig() (Config, error) {
	timeout, err := envDuration("TIMEOUT", defaultTimeout)
	if err != nil {
		return Config{}, err
	}

	cnf := Config{
		MovieBaseURL: envString("MOVIE_BASE_URL", defaultMovieBaseURL),
		Timeout:      timeout,
	}

	configValidationErr := configValidate(cnf.Timeout)
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

func configValidate(timeout time.Duration) error {
	if timeout <= 0 {
		return fmt.Errorf("TIMEOUT must be > 0, got %s", timeout)
	}

	return nil
}
