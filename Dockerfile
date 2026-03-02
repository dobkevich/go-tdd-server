# Build stage
FROM golang:1.24-alpine AS builder

# Install necessary system dependencies
RUN apk add --no-cache git

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build optimized binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server main.go

# Run stage
FROM alpine:3.21

# Add non-root user for security
RUN adduser -D -u 1001 appuser
USER 1001

WORKDIR /home/appuser/

# Copy binary from builder stage
COPY --from=builder /github.com/project/go-tdd-server/server .

# Default environment variables
ENV PORT=8080

EXPOSE 8080

# Run the server
CMD ["./server"]
