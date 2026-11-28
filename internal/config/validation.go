package config

import "fmt"

func validate(config Config) error {
	if config.Timeout <= 0 {
		return fmt.Errorf("TIMEOUT must be > 0, got %s", config.Timeout)
	}

	if config.WorkersCount <= 0 {
		return fmt.Errorf("WORKERS must be > 0, got %d", config.WorkersCount)
	}

	return nil
}
