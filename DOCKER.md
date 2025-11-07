# Docker Build Guide

This project supports two Docker build modes:

## 1. Multi-Stage Build (Standalone)

Builds the binary from source inside Docker. Use this for local development or when you don't have a prebuilt binary.

```bash
# Build with default version info
docker build -t go-app:latest .

# Build with custom version info
docker build \
  --build-arg VERSION=1.0.0 \
  --build-arg COMMIT=$(git rev-parse --short HEAD) \
  --target runtime \
  -t go-app:1.0.0 \
  .

# Explicitly target the runtime stage
docker build --target runtime -t go-app:latest .
```

## 2. Prebuilt Binary Mode (CI/CD)

Uses an existing binary from the `bin/` directory. This is faster and used in CI/CD pipelines.

```bash
# First, build the binary
go build -ldflags="-s -w -X main.version=1.0.0 -X main.commit=abc123" \
  -o bin/server ./cmd/server

# Then build the Docker image
docker build --target runtime-prebuilt -t go-app:latest .
```

## Running the Container

```bash
# Run in foreground
docker run -p 8080:8080 go-app:latest

# Run in background
docker run -d -p 8080:8080 --name my-app go-app:latest

# Test health endpoint
curl http://localhost:8080/health

# View logs
docker logs my-app

# Stop and remove
docker stop my-app && docker rm my-app
```

## Environment Variables

Configure the application using environment variables:

```bash
docker run -p 8080:8080 \
  -e PORT=3000 \
  -e HOST=0.0.0.0 \
  -e ENVIRONMENT=production \
  -e READ_TIMEOUT=15s \
  -e WRITE_TIMEOUT=15s \
  -e IDLE_TIMEOUT=60s \
  -e RATE_LIMIT_RPS=50 \
  go-app:latest
```

Available variables:
- `PORT` (default: 8080)
- `HOST` (default: 0.0.0.0)
- `ENVIRONMENT` (default: production)
- `READ_TIMEOUT` (default: 10s)
- `WRITE_TIMEOUT` (default: 10s)
- `IDLE_TIMEOUT` (default: 120s)
- `SHUTDOWN_TIMEOUT` (default: 30s)
- `RATE_LIMIT_RPS` (default: 100)
- `RATE_LIMIT_BURST` (default: 200)

## CI/CD Usage

The GitHub Actions workflow (`.github/workflows/ci.yml`) uses the prebuilt binary mode:

1. Calculates version information from git tags and commits
2. Builds the binary with Go in the CI environment (with semver format)
3. Uses `--target runtime-prebuilt` to create a minimal Docker image
4. Tags the image with a Docker-compatible version (replaces `+` with `-`)
5. Tests the image by running it and checking the health endpoint
6. Saves the image as an artifact

### Version Format

The workflow generates two version strings:

- **Binary version** (semver): `0.0.0-dev.14+6636b22`
  - Used in `-ldflags` for the binary
  - Follows semantic versioning specification

- **Docker tag** (Docker-compatible): `0.0.0-dev.14-6636b22`
  - Used for Docker image tags
  - Replaces `+` with `-` (Docker doesn't allow `+` in tags)

This approach provides:
- **Faster builds**: Binary is cached and reused
- **Smaller images**: Uses `scratch` as base (~10MB final image)
- **Better security**: Minimal attack surface with no shell or OS packages
- **Efficient caching**: Go build cache works across builds
- **Compliant versioning**: Semver for binaries, Docker-safe for images

## Image Details

- **Base**: `scratch` (minimal, no OS)
- **Size**: ~10MB
- **User**: Non-root (UID 65534)
- **Port**: 8080
- **Includes**: CA certificates, timezone data

## Troubleshooting

### Binary not found in prebuilt mode
Ensure `bin/server` exists before building:
```bash
ls -la bin/server
```

### Permission denied
Make sure the binary is executable:
```bash
chmod +x bin/server
```

### Health check fails
Check container logs:
```bash
docker logs <container-name>
```

Verify port mapping:
```bash
docker ps
curl http://localhost:<mapped-port>/health
```
