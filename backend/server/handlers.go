package server

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"backend/config"
	"api"
)

// BackendServer implements the BackendService gRPC interface
type BackendServer struct {
	api.UnimplementedBackendServiceServer
	config *config.Config
}

// NewBackendServer creates a new backend server instance
func NewBackendServer(cfg *config.Config) *BackendServer {
	return &BackendServer{
		config: cfg,
	}
}

// ProcessData simulates CPU-intensive work
func (s *BackendServer) ProcessData(ctx context.Context, req *api.ProcessDataRequest) (*api.ProcessDataResponse, error) {
	start := time.Now()

	// Simulate CPU work based on complexity and pod's CPU factor
	complexity := req.Complexity
	if complexity <= 0 {
		complexity = 1
	}
	if complexity > 10 {
		complexity = 10
	}

	// CPU simulation - mathematical operations
	result := s.simulateCPUWork(req.Data, int(complexity))

	// Add some random variation based on pod characteristics
	baseDelay := time.Duration(s.config.LatencyBase) * time.Millisecond
	variableDelay := time.Duration(rand.Intn(50)) * time.Millisecond
	cpuDelay := time.Duration(float64(complexity*10)*s.config.CPUFactor) * time.Millisecond

	totalDelay := baseDelay + variableDelay + cpuDelay
	time.Sleep(totalDelay)

	processingTime := time.Since(start)

	return &api.ProcessDataResponse{
		Result:           result,
		ProcessingTimeMs: processingTime.Milliseconds(),
		PodId:            s.config.PodID,
	}, nil
}

// GetResource simulates I/O operations
func (s *BackendServer) GetResource(ctx context.Context, req *api.GetResourceRequest) (*api.GetResourceResponse, error) {
	start := time.Now()

	// Simulate I/O delay
	baseDelay := time.Duration(s.config.LatencyBase/2) * time.Millisecond
	if req.SimulateSlowOperation {
		baseDelay *= 3 // Make it slower for testing
	}

	// Add random jitter
	jitter := time.Duration(rand.Intn(30)) * time.Millisecond
	time.Sleep(baseDelay + jitter)

	// Generate response data
	resourceData := fmt.Sprintf("Resource-%s-Data-%d", req.ResourceId, time.Now().Unix())
	metadata := fmt.Sprintf("Processed by %s in %v", s.config.PodID, time.Since(start))

	return &api.GetResourceResponse{
		ResourceData: resourceData,
		Metadata:     metadata,
		PodId:        s.config.PodID,
	}, nil
}

// HealthCheck returns the health status of the pod
func (s *BackendServer) HealthCheck(ctx context.Context, req *api.HealthCheckRequest) (*api.HealthCheckResponse, error) {
	// Simple health check - always return SERVING for now
	// In a real scenario, you might check database connections, external services, etc.
	
	return &api.HealthCheckResponse{
		Status: api.HealthCheckResponse_SERVING,
		PodId:  s.config.PodID,
	}, nil
}

// simulateCPUWork performs mathematical operations to simulate CPU load
func (s *BackendServer) simulateCPUWork(data string, complexity int) string {
	// Perform some CPU-intensive operations
	var result float64 = 1.0
	iterations := complexity * 1000 * int(s.config.CPUFactor)

	for i := 0; i < iterations; i++ {
		result = math.Sin(result) + math.Cos(float64(i))
	}

	// String processing to add more CPU work
	processed := strings.ToUpper(data)
	for i := 0; i < complexity; i++ {
		processed = fmt.Sprintf("%s-%d-%.2f", processed, i, result)
		if len(processed) > 1000 {
			processed = processed[:1000] // Prevent excessive memory usage
		}
	}

	return fmt.Sprintf("Processed[%s] by pod %s (complexity=%d, iterations=%d, result=%.2f)", 
		processed, s.config.PodID, complexity, iterations, result)
}