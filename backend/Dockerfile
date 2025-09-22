# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Copy the API module
COPY ../api ./api

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

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