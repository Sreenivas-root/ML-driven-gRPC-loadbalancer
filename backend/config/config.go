package config

import (
	"os"
	"strconv"
)

// Config holds the configuration for the backend service
type Config struct {
	PodID       string
	GRPCPort    string
	MetricsPort string
	// Performance characteristics for this pod (for ML testing)
	CPUFactor    float64 // Multiplier for CPU simulation (0.5-2.0)
	LatencyBase  int     // Base latency in ms (10-100)
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	config := &Config{
		PodID:       getEnv("POD_ID", "backend-pod-unknown"),
		GRPCPort:    getEnv("GRPC_PORT", "9090"),
		MetricsPort: getEnv("METRICS_PORT", "8080"),
		CPUFactor:   getEnvFloat("CPU_FACTOR", 1.0),
		LatencyBase: getEnvInt("LATENCY_BASE", 50),
	}

	return config
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}