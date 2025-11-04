# Go Web Application

A production-ready, high-performance Go backend web application designed to handle millions of requests per day. Built with the [Chi router](https://github.com/go-chi/chi) and featuring comprehensive middleware, testing, and CI/CD infrastructure.

## Features

- **High Performance**: Optimized for handling millions of requests per day
- **Production-Ready**: Built with best practices for scalability and reliability
- **Chi Router**: Fast and lightweight HTTP router with middleware support
- **Comprehensive Middleware**:
  - Request logging
  - Panic recovery
  - Metrics collection
  - Rate limiting (per-IP)
  - CORS support
  - Request compression
- **Health Checks**: Multiple health check endpoints (liveness, readiness)
- **Metrics**: Built-in metrics collection and reporting
- **Configuration**: Environment-based configuration with sensible defaults
- **Testing**: Comprehensive unit tests with >80% coverage
- **CI/CD**: GitHub Actions workflows for testing, building, and deployment
- **Containerization**: Docker and docker-compose support
- **Security**: Built-in security scanning and best practices

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerization)
- Make (optional, for convenience commands)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/go-app.git
cd go-app

# Download dependencies
go mod download

# Build the application
go build -o bin/server ./cmd/server

# Run the application
./bin/server
```

Or using Make:

```bash
make build
make run
```

### Using Docker

```bash
# Build and run with Docker
docker build -t go-app .
docker run -p 8080:8080 go-app

# Or use docker-compose
docker-compose up
```

## Configuration

The application is configured through environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `HOST` | `0.0.0.0` | Server host |
| `ENVIRONMENT` | `production` | Environment (production/development) |
| `READ_TIMEOUT` | `10s` | HTTP read timeout |
| `WRITE_TIMEOUT` | `10s` | HTTP write timeout |
| `IDLE_TIMEOUT` | `120s` | HTTP idle timeout |
| `SHUTDOWN_TIMEOUT` | `30s` | Graceful shutdown timeout |
| `RATE_LIMIT_RPS` | `100` | Rate limit requests per second |
| `RATE_LIMIT_BURST` | `200` | Rate limit burst size |

### Example Configuration

```bash
export PORT=9000
export ENVIRONMENT=development
export RATE_LIMIT_RPS=1000
./bin/server
```

## API Endpoints

### Health Checks

- `GET /health` - Full health check with version and uptime
- `GET /health/ready` - Readiness probe
- `GET /health/live` - Liveness probe

### Metrics

- `GET /metrics` - Application metrics (requests, errors, latency, etc.)

### API v1

- `GET /api/v1/hello` - Example endpoint

## Development

### Project Structure

```
.
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Custom middleware
│   └── metrics/         # Metrics collection
├── pkg/
│   └── health/          # Health check functionality
├── test/                # Integration tests
├── .github/
│   └── workflows/       # CI/CD pipelines
├── Dockerfile           # Docker build configuration
├── docker-compose.yml   # Docker Compose configuration
├── Makefile            # Build and dev commands
└── README.md           # This file
```

### Available Make Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Build and run the application
make test              # Run tests
make test-coverage     # Run tests with coverage report
make lint              # Run linters
make format            # Format code
make clean             # Clean build artifacts
make docker-build      # Build Docker image
make docker-run        # Run Docker container
make ci                # Run CI pipeline locally
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -cover ./...

# Run tests with race detection
go test -race ./...

# Or use Make
make test
make test-coverage
```

### Code Quality

```bash
# Run linters
make lint

# Run security scanner
make security

# Format code
make format
```

## CI/CD

The project includes GitHub Actions workflows for:

1. **CI Pipeline** (`.github/workflows/ci.yml`):
   - Runs on push and pull requests
   - Tests on multiple Go versions
   - Runs linters and security scans
   - Generates coverage reports
   - Builds artifacts

2. **Docker Build** (`.github/workflows/docker.yml`):
   - Builds multi-platform Docker images
   - Pushes to GitHub Container Registry
   - Runs vulnerability scanning

3. **Release** (`.github/workflows/release.yml`):
   - Automated releases on version tags
   - Cross-platform binary builds
   - GitHub release creation

## Deployment

### Docker Deployment

```bash
# Build image
docker build -t go-app:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -e PORT=8080 \
  -e ENVIRONMENT=production \
  --name go-app \
  go-app:latest
```

### Kubernetes Deployment

Example Kubernetes deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-app
  template:
    metadata:
      labels:
        app: go-app
    spec:
      containers:
      - name: go-app
        image: ghcr.io/yourusername/go-app:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: ENVIRONMENT
          value: "production"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

## Performance

The application is designed to handle high loads:

- Efficient request routing with Chi
- Connection pooling and keep-alive
- Request/response compression
- Per-IP rate limiting
- Graceful shutdown
- Low memory footprint
- Minimal dependencies

### Benchmarking

Run benchmarks:

```bash
make benchmark
```

## Monitoring

The application exposes metrics at `/metrics` endpoint. You can integrate with:

- Prometheus for metrics collection
- Grafana for visualization
- DataDog, New Relic, or similar APM tools

Example Prometheus configuration is included in the docker-compose file (commented out).

## Security

Security features:

- Panic recovery middleware
- Rate limiting per IP
- Security headers (can be added via middleware)
- Vulnerability scanning in CI/CD
- Static analysis with gosec
- Docker image scanning with Trivy

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues, questions, or contributions, please open an issue on GitHub.

## Acknowledgments

- [Chi Router](https://github.com/go-chi/chi) - Lightweight, idiomatic HTTP router
- [Go](https://golang.org/) - The Go programming language
