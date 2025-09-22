# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy backend source code first (except go.mod files)
COPY backend/ .

# Copy the API module from project root
COPY api/ ./api/

# Fix the replace directive for Docker build
RUN sed -i 's|replace api => ../api|replace api => ./api|' go.mod

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/backend .

# Expose ports
EXPOSE 8080 9090

# Set default environment variables
ENV POD_ID=backend-pod-default
ENV GRPC_PORT=9090
ENV METRICS_PORT=8080
ENV CPU_FACTOR=1.0
ENV LATENCY_BASE=50

# Run the application
CMD ["./backend"]