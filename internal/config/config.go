package config

import "time"

type Config struct {
	MovieBaseURL string
	DatabaseURL  string
	RedisAddr    string
	Timeout      time.Duration

	HTTPAddr string
	GRPCAddr string

	WorkersCount int
}

const (
	defaultMovieBaseURL = "https://homeworksite.site"
	defaultDatabaseURL  = "postgresql://vblinov@localhost:5432/moviesdb?sslmode=disable"
	defaultRedisAddr    = "localhost:6379"
	defaultTimeout      = 5 * time.Second

	defaultHTTPAddr = ":8080"
	defaultGRPCAddr = ":9090"

	defaultWorkersCount = 30
)

func LoadConfig() (Config, error) {
	timeout, err := envDuration("TIMEOUT", defaultTimeout)
	if err != nil {
		return Config{}, err
	}

	workersCount, err := envInt("WORKERS", defaultWorkersCount)
	if err != nil {
		return Config{}, err
	}

	config := Config{
		MovieBaseURL: envString("MOVIE_BASE_URL", defaultMovieBaseURL),
		DatabaseURL:  envString("DATABASE_URL", defaultDatabaseURL),
		RedisAddr:    envString("REDIS_ADDR", defaultRedisAddr),
		Timeout:      timeout,

		HTTPAddr: envString("HTTP_ADDR", defaultHTTPAddr),
		GRPCAddr: envString("GRPC_ADDR", defaultGRPCAddr),

		WorkersCount: workersCount,
	}

	configValidationErr := validate(config)
	if configValidationErr != nil {
		return Config{}, configValidationErr
	}

	return config, nil
}
