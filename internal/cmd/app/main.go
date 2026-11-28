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
	"github.com/BlVladimir/GolangCourse/internal/handler/grpc_api"
	"github.com/BlVladimir/GolangCourse/internal/handler/http_api"
	"github.com/BlVladimir/GolangCourse/internal/usecase"
)

func main() {
	clientConfig, err := filmstand.LoadClientConfig()
	if err != nil {
		log.Fatalf("client config: %v", err)
	}

	serviceConfig, err := usecase.LoadUseCaseConfig()
	if err != nil {
		log.Fatalf("usecase config: %v", err)
	}

	provider := filmstand.NewClient(clientConfig)
	service := usecase.NewFilmService(serviceConfig, provider)

	// --- HTTP ---
	httpConfig, err := http_api.LoadHTTPConfig()
	if err != nil {
		log.Fatalf("http config: %v", err)
	}

	httpServer := &http.Server{
		Addr:              httpConfig.HTTPAddr,
		Handler:           http_api.NewRouter(http_api.New(service)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Printf("HTTP listening on %s", httpConfig.HTTPAddr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	// // --- gRPC ---
	grpcConfig, err := grpc_api.LoadGRPCConfig()
	if err != nil {
		log.Fatalf("http config: %v", err)
	}
	grpcSrv := grpc_api.NewGRPCServer(service)

	lis, err := net.Listen("tcp", grpcConfig.GRPCAddr)
	if err != nil {
		log.Fatalf("grpc listen %s: %v", grpcConfig.GRPCAddr, err)
	}

	go func() {
		log.Printf("gRPC listening on %s", grpcConfig.GRPCAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("grpc: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("shutting down...")

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
