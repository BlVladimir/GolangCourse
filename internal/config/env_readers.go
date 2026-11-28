package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func envString(key, def string) string {
	v := os.Getenv(key)
	if v != "" {
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
