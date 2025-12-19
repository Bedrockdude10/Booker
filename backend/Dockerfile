# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git (needed for go mod download with private repos)
RUN apk add --no-cache git

# Copy go mod files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
# CGO_ENABLED=0 for static binary, GOOS=linux for Linux container
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Runtime stage - minimal image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS (needed for MongoDB Atlas)
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/server .

# Cloud Run uses PORT env var
ENV PORT=8080
EXPOSE 8080

# Run the binary
CMD ["./server"]