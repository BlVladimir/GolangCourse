package grpc_api

import "os"

type Config struct {
	GRPCAddr string
}

const defaultGRPCAddr = ":9090"

func LoadGRPCConfig() (Config, error) {
	cnf := Config{
		GRPCAddr: envString("GRPC_ADDR", defaultGRPCAddr),
	}

	return cnf, nil
}

func envString(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
