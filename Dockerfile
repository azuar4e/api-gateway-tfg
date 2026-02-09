# Build Stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# We build the binary named 'main' from the cmd/main.go file
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Runtime Stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 9090

# Command to run the executable
CMD ["./main"]