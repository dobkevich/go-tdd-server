# Build stage
FROM golang:1.24 AS builder

# Install necessary system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build optimized binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server main.go

# Run stage
FROM alpine:3.23

# Add non-root user for security
RUN adduser -D -u 1001 appuser
USER 1001

WORKDIR /home/appuser/

# Copy binary from builder stage
COPY --from=builder /app/server .

# Default environment variables
ENV PORT=8080

EXPOSE 8080

# Run the server
CMD ["./server"]
