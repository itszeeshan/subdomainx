FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o subdomainx ./cmd/subdomainx

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/subdomainx .

# Copy config files
COPY --from=builder /app/configs ./configs

# Make binary executable
RUN chmod +x subdomainx

# Expose port (if needed for web interface in future)
EXPOSE 8080

# Run the binary
ENTRYPOINT ["./subdomainx"]
