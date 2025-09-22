package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/config"
	"backend/metrics"
	"backend/server"
)

func main() {
	fmt.Println("Starting ML-driven gRPC Load Balancer Backend Service...")

	// Load configuration
	cfg := config.LoadConfig()
	fmt.Printf("Pod ID: %s\n", cfg.PodID)
	fmt.Printf("gRPC Port: %s\n", cfg.GRPCPort)
	fmt.Printf("Metrics Port: %s\n", cfg.MetricsPort)
	fmt.Printf("CPU Factor: %.2f\n", cfg.CPUFactor)
	fmt.Printf("Latency Base: %d ms\n", cfg.LatencyBase)

	// Initialize metrics
	metrics.Initialize(cfg.PodID)
	fmt.Println("Metrics initialized")

	// Create gRPC server
	grpcServer, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create gRPC server: %v", err)
	}

	// Start metrics server in a goroutine
	go func() {
		fmt.Printf("Starting metrics server on port %s...\n", cfg.MetricsPort)
		if err := metrics.StartMetricsServer(cfg.MetricsPort); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Start gRPC server in a goroutine
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	fmt.Println("Backend service started successfully!")
	fmt.Printf("gRPC server listening on :%s\n", cfg.GRPCPort)
	fmt.Printf("Metrics server listening on :%s/metrics\n", cfg.MetricsPort)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down backend service...")
	grpcServer.Stop()
	fmt.Println("Backend service stopped")
}