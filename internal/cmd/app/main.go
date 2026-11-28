package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BlVladimir/GolangCourse/internal/adapter/filmstand"
	"github.com/BlVladimir/GolangCourse/internal/adapter/postgres"
	"github.com/BlVladimir/GolangCourse/internal/adapter/redis_searcher"
	"github.com/BlVladimir/GolangCourse/internal/config"
	"github.com/BlVladimir/GolangCourse/internal/handler/grpc_api"
	"github.com/BlVladimir/GolangCourse/internal/handler/http_api"
	"github.com/BlVladimir/GolangCourse/internal/usecase"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cnf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	repository, err := postgres.New(ctx, cnf.DatabaseURL)
	if err != nil {
		log.Fatalf("repository: %v", err)
	}

	err = repository.Migrate(ctx)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: cnf.RedisAddr,
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}

	searcher := redis_searcher.NewMoviesSearcher(redisClient)
	provider := filmstand.NewClient(cnf)

	service := usecase.NewFilmService(cnf, provider, repository, searcher)
	err = service.SyncDatabase(ctx)
	if err != nil {
		log.Fatalf("sync: %v", err)
	}

	if err := service.BuildIndex(ctx); err != nil {
		log.Fatalf("build index: %v", err)
	}

	// --- HTTP ---
	httpServer := &http.Server{
		Addr:              cnf.HTTPAddr,
		Handler:           http_api.NewRouter(http_api.New(service)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("HTTP listening on %s", cnf.HTTPAddr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	// --- gRPC ---
	grpcSrv := grpc_api.NewGRPCServer(service)

	lis, err := net.Listen("tcp", cnf.GRPCAddr)
	if err != nil {
		log.Fatalf("grpc listen %s: %v", cnf.GRPCAddr, err)
	}

	go func() {
		log.Printf("gRPC listening on %s", cnf.GRPCAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("grpc: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}

	done := make(chan struct{})
	go func() {
		grpcSrv.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
		log.Println("gRPC stopped gracefully")
	case <-shutdownCtx.Done():
		log.Println("gRPC shutdown timeout, forcing stop")
		grpcSrv.Stop()
	}
}
