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

	"github.com/BlVladimir/GolangCourse/internal/film"
	"github.com/BlVladimir/GolangCourse/internal/grpc_api"
	"github.com/BlVladimir/GolangCourse/internal/http_api"
)

func main() {
	cfg, err := film.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	svc := film.NewFilmService(cfg)

	// --- HTTP ---
	httpSrv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           http_api.NewRouter(http_api.New(svc)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// --- gRPC ---
	grpcSrv := grpc_api.NewGRPCServer(svc)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("grpc listen %s: %v", cfg.GRPCAddr, err)
	}

	go func() {
		log.Printf("HTTP listening on %s", cfg.HTTPAddr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	go func() {
		log.Printf("gRPC listening on %s", cfg.GRPCAddr)
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

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
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
