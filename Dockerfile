# Stage 1: Build stage
FROM golang:1.23.3 as builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /myapp ./cmd

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary tools
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /myapp .

# Copy the .env file if it exists in the root directory
COPY .env .env

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./myapp"]
