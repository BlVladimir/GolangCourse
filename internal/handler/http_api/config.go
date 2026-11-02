package http_api

import "os"

type Config struct {
	HTTPAddr string
}

const defaultHTTPAddr = ":8080"

func LoadHTTPConfig() (Config, error) {
	cnf := Config{
		HTTPAddr: envString("GRPC_ADDR", defaultHTTPAddr),
	}

	return cnf, nil
}

func envString(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
