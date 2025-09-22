package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"backend/config"
	"backend/metrics"
	"api"
)

// Server wraps the gRPC server
type Server struct {
	grpcServer *grpc.Server
	config     *config.Config
	listener   net.Listener
}

// NewServer creates a new gRPC server
func NewServer(cfg *config.Config) (*Server, error) {
	// Create listener
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(metrics.UnaryServerInterceptor()),
	)

	// Register our service
	backendServer := NewBackendServer(cfg)
	api.RegisterBackendServiceServer(grpcServer, backendServer)

	// Enable reflection for debugging/testing
	reflection.Register(grpcServer)

	return &Server{
		grpcServer: grpcServer,
		config:     cfg,
		listener:   listener,
	}, nil
}

// Start starts the gRPC server
func (s *Server) Start() error {
	fmt.Printf("Starting gRPC server on port %s...\n", s.config.GRPCPort)
	return s.grpcServer.Serve(s.listener)
}

// Stop gracefully stops the gRPC server
func (s *Server) Stop() {
	fmt.Println("Stopping gRPC server...")
	s.grpcServer.GracefulStop()
}

// GetPort returns the port the server is listening on
func (s *Server) GetPort() string {
	return s.config.GRPCPort
}