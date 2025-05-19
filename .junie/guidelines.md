# Konfig - Developer Guidelines

## Project Overview

Konfig is a Spring-Framework-like configuration and profile manager for Go applications. It provides:

- YAML configuration loading with profile support
- Environment variable expansion and overrides
- Nested configuration with dot notation
- Default values for properties

## Project Structure

- `resources/` - Default directory for configuration files
  - `application.yaml` - Main configuration file
  - `application-{profile}.yaml` - Profile-specific configuration files
- `test-data/` - Test configuration files
- `test-proj/` - Example project using the library

## Configuration Files

- Use YAML format with `.yaml` or `.yml` extension
- Main configuration file: `application.yaml`
- Profile-specific files: `application-{profile}.yaml`
- Environment variables can be referenced with `${ENV_VAR}`
- Default values can be specified with `${ENV_VAR:default_value}`

## Running Tests

```bash
# Run all tests with race detection
make test

# Run specific tests
go test -v ./... -run TestName
```

## Building and Installing

```bash
# Build the project
make build

# Install the package
make install

# Update dependencies
make deps
```

## Using the Library

### Basic Usage

```go
package main

import (
    "github.com/mfenderov/konfig"
    "log"
)

func main() {
    // Load configuration
    err := konfig.Load()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Access configuration values
    port := konfig.GetEnv("server.port")
    dbHost := konfig.GetEnv("db.host")

    // Use the values in your application
    log.Printf("Server port: %s", port)
    log.Printf("Database host: %s", dbHost)
}
```

### Using Profiles

Profiles can be specified using the `--profile` or `-p` flag:

```bash
# Run with dev profile
go run main.go --profile dev

# Run with prod profile
go run main.go -p prod
```

You can check the current profile in your code:

```go
if konfig.IsDevProfile() {
    // Development-specific code
}

if konfig.IsProdProfile() {
    // Production-specific code
}

if konfig.IsProfile("test") {
    // Test-specific code
}
```

## Best Practices

### Configuration

- Keep configuration files organized and well-commented
- Use nested structures for related configuration
- Provide sensible default values
- Use environment variables for sensitive information

### Error Handling

- Always wrap errors with context using `errors.Wrap`
- Check for nil values in configuration
- Handle missing configuration files gracefully

### Testing

- Write tests for different configuration scenarios
- Test with various profiles
- Test environment variable overrides
- Test error cases (missing files, invalid YAML, etc.)

## References

- [Go Modules Documentation](https://go.dev/ref/mod)
- [YAML Specification](https://yaml.org/spec/)
- [Testify Library](https://github.com/stretchr/testify)
