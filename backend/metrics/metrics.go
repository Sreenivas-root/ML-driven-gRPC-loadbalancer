package metrics

import (
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// gRPC server handling duration histogram (required for load balancer)
	GrpcServerHandlingSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_handling_seconds",
			Help:    "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"grpc_method", "grpc_service", "pod"},
	)

	// Request counter
	GrpcServerRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_server_requests_total",
			Help: "Total number of gRPC requests processed",
		},
		[]string{"grpc_method", "grpc_service", "pod", "status"},
	)

	// Container CPU usage simulation (for ML model)
	ContainerCpuUsageSecondsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "container_cpu_usage_seconds_total",
			Help: "Cumulative cpu time consumed by the container in seconds",
		},
		[]string{"pod"},
	)

	// Container memory working set (for ML model)
	ContainerMemoryWorkingSetBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "container_memory_working_set_bytes",
			Help: "Current working set memory of the container in bytes",
		},
		[]string{"pod"},
	)

	// Pod identification
	podID string
)

// Initialize sets up all Prometheus metrics
func Initialize(podIdentifier string) {
	podID = podIdentifier

	// Register metrics
	prometheus.MustRegister(GrpcServerHandlingSeconds)
	prometheus.MustRegister(GrpcServerRequestsTotal)
	prometheus.MustRegister(ContainerCpuUsageSecondsTotal)
	prometheus.MustRegister(ContainerMemoryWorkingSetBytes)

	// Start background goroutine to update system metrics
	go updateSystemMetrics()
}

// updateSystemMetrics updates CPU and memory metrics periodically
func updateSystemMetrics() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	cpuAccumulator := 0.0

	for range ticker.C {
		// Simulate CPU usage based on actual runtime + some randomness
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// Update memory metric (convert to MB for easier reading)
		ContainerMemoryWorkingSetBytes.WithLabelValues(podID).Set(float64(m.Alloc))

		// Simulate CPU accumulation (this is a simple simulation)
		// In a real scenario, you'd read from /proc/stat or cgroups
		cpuAccumulator += 0.1 + (float64(time.Now().UnixNano()%100) / 1000.0) // Add some variance
		ContainerCpuUsageSecondsTotal.WithLabelValues(podID).Add(cpuAccumulator)
	}
}

// StartMetricsServer starts the HTTP server for Prometheus metrics
func StartMetricsServer(port string) error {
	http.Handle("/metrics", promhttp.Handler())
	
	// Add a health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return http.ListenAndServe(":"+port, nil)
}

// RecordGRPCRequest records metrics for a gRPC request
func RecordGRPCRequest(method, service, status string, duration time.Duration) {
	GrpcServerHandlingSeconds.WithLabelValues(method, service, podID).Observe(duration.Seconds())
	GrpcServerRequestsTotal.WithLabelValues(method, service, podID, status).Inc()
}