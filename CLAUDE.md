# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go application designed to demonstrate CI/CD pipeline capabilities.

## Development Commands

### Building
```bash
go build -o bin/go-app ./cmd/...
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -cover ./...

# Run a specific test
go test -v -run TestName ./path/to/package
```

### Linting
```bash
# Run go vet
go vet ./...

# Run golangci-lint (if configured)
golangci-lint run
```

### Running the Application
```bash
# Run directly
go run ./cmd/...

# Run built binary
./bin/go-app
```

## Project Structure

This is a new project. As it grows, follow standard Go project layout conventions:
- `cmd/` - Main applications for this project
- `pkg/` - Library code that's ok to use by external applications
- `internal/` - Private application and library code
- `api/` - API definitions (OpenAPI/Swagger specs, protocol definitions)
- `test/` - Additional external test apps and test data

## CI/CD Considerations

Since this project is designed for CI/CD pipeline demonstrations:
- Ensure all builds are reproducible
- Tests should be fast and deterministic
- Consider containerization (Dockerfile) for consistent deployments
- Include health check endpoints for deployment verification
