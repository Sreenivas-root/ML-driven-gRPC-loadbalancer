package metrics

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a gRPC unary server interceptor for metrics collection
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		// Record metrics
		duration := time.Since(start)
		statusCode := codes.OK
		if err != nil {
			statusCode = status.Code(err)
		}

		// Extract method and service names from the full method
		service, method := parseFullMethod(info.FullMethod)
		
		RecordGRPCRequest(method, service, statusCode.String(), duration)

		return resp, err
	}
}

// parseFullMethod extracts service and method names from full method string
// e.g., "/loadbalancer.BackendService/ProcessData" -> ("BackendService", "ProcessData")
func parseFullMethod(fullMethod string) (service, method string) {
	// fullMethod format: "/package.Service/Method"
	if len(fullMethod) == 0 || fullMethod[0] != '/' {
		return "unknown", "unknown"
	}

	// Remove leading slash
	fullMethod = fullMethod[1:]

	// Find the last slash
	lastSlash := -1
	for i := len(fullMethod) - 1; i >= 0; i-- {
		if fullMethod[i] == '/' {
			lastSlash = i
			break
		}
	}

	if lastSlash == -1 {
		return "unknown", fullMethod
	}

	servicePart := fullMethod[:lastSlash]
	method = fullMethod[lastSlash+1:]

	// Extract service name from package.Service
	lastDot := -1
	for i := len(servicePart) - 1; i >= 0; i-- {
		if servicePart[i] == '.' {
			lastDot = i
			break
		}
	}

	if lastDot == -1 {
		service = servicePart
	} else {
		service = servicePart[lastDot+1:]
	}

	return service, method
}