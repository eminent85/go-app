# Dockerfile supports two modes:
# 1. Multi-stage build (default): Builds the binary from source
# 2. Prebuilt binary: Uses existing binary from bin/server (for CI/CD)
#
# Usage:
#   Multi-stage:  docker build .
#   Prebuilt:     docker build --target runtime-prebuilt .

# ============================================
# Base image with runtime dependencies
# ============================================
FROM alpine:latest AS runtime-base

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# ============================================
# Builder stage - builds binary from source
# ============================================
FROM golang:1.25-alpine AS builder

# Build arguments for version information
ARG VERSION=dev
ARG COMMIT=unknown

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application with version information
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static' -X main.version=${VERSION} -X main.commit=${COMMIT}" \
    -a \
    -o /build/bin/server \
    ./cmd/server

# ============================================
# Runtime stage - for multi-stage builds
# ============================================
FROM scratch AS runtime

# Copy CA certificates and timezone data
COPY --from=runtime-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=runtime-base /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary from builder
COPY --from=builder /build/bin/server /server

# Use non-root user
USER 65534:65534

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/server"]

# ============================================
# Runtime stage - for prebuilt binaries (CI/CD)
# ============================================
FROM scratch AS runtime-prebuilt

# Copy CA certificates and timezone data
COPY --from=runtime-base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=runtime-base /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the prebuilt binary from local bin directory
COPY bin/server /server

# Use non-root user
USER 65534:65534

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/server"]
