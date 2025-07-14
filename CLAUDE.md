# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

konfig is a Go library for configuration management that provides Spring Framework-like configuration capabilities. It
supports YAML-based configuration files, profile-specific configurations, environment variable substitution, and
struct-based configuration loading with type safety.

## Development Commands

### Testing

- `make test` - Run all tests with race detection and vet checks
- `go test ./...` - Run tests without make
- `go test -v ./...` - Run tests with verbose output
- `go test -run TestSpecificFunction` - Run specific test
- `cd test-proj && go test` - Run critical integration tests (external module usage)

### Building

- `make build` - Build the project
- `go build` - Direct go build

### Quality & Documentation

- `make lint` - Run golangci-lint (informational)
- `make coverage` - Generate test coverage report (HTML + summary)
- `make quality` - Run all quality checks (lint + test + coverage)
- `make ci` - Complete CI pipeline (deps + quality)
- `make docs` - Generate text documentation
- `make docs-serve` - Start documentation server on :6060

### Other Commands

- `make deps` - Tidy dependencies
- `make install` - Install the package
- `make release` - Release process (requires GITHUB_REF_NAME)

## Architecture

### Core Components

1. **Main Configuration Loader** (`konfig.go`):
    - `Load()` - Loads default application.yaml and active profile configurations
    - `LoadFrom(filepath)` - Loads configuration from specific file
    - Environment variable management and YAML processing

2. **Struct-Based Configuration** (`struct_loader.go`):
    - `LoadInto(&config)` - Populates Go structs using reflection and `konfig` tags
    - Supports nested structs and default values via struct tags
    - Type-safe configuration access

3. **Profile Management** (`profile.go`):
    - Profile detection via command-line flags (`-p` or `--profile`)
    - Profile-specific configuration file loading (application-{profile}.yaml)
    - Profile helper functions: `IsDevProfile()`, `IsProdProfile()`, `IsProfile(name)`

4. **Resource Discovery** (`resource_finder.go`):
    - Automatic discovery of configuration files in `resources/` directory
    - Support for both `.yaml` and `.yml` extensions
    - Fallback mechanisms for missing files

5. **Configuration Loading** (`loader.go`):
    - YAML parsing and environment variable substitution
    - Key flattening (nested keys become dot-separated)
    - Environment variable override capabilities

### Configuration File Structure

- Default location: `resources/` directory
- Base configuration: `application.yaml` or `application.yml`
- Profile-specific: `application-{profile}.yaml` (e.g., `application-dev.yaml`)
- Environment variable substitution: `${VAR_NAME:default_value}`

### Key Features

1. **Environment Variable Integration**: All configuration keys are exposed as environment variables with dot notation
2. **Profile Support**: Different configurations for dev, prod, test environments
3. **Struct Tags**: Type-safe configuration using `konfig:"key.path"` and `default:"value"` tags
4. **Nested Configuration**: Deep nesting support for complex configuration structures
5. **Command-line Profile Selection**: Use `-p dev` or `--profile prod` flags

## Testing Strategy

- **Simplified test structure**: 6 focused test files (reduced from 8)
- Unit tests for each component (`*_test.go` files)
- Integration tests with actual YAML files (`test-data/` directory)
- **Essential benchmarks only**: 3 core performance tests (reduced from 6)
- **Critical Integration Tests**: `test-proj/` validates the public API and real-world usage

### Test Files Structure

- `konfig_test.go` - Core functionality tests
- `profile_test.go` - Profile management tests
- `resource_finder_test.go` - Resource discovery tests
- `struct_loader_test.go` - Basic struct-based configuration tests
- `struct_loader_advanced_test.go` - Advanced struct tests (nested, error handling)
- `struct_loader_benchmark_test.go` - Essential performance benchmarks

## Dependencies

- `github.com/pkg/errors` - Enhanced error handling
- `github.com/spf13/pflag` - Command-line flag parsing
- `github.com/stretchr/testify` - Testing utilities
- `gopkg.in/yaml.v3` - YAML parsing

## Development Philosophy

konfig follows a **"merciless simplification"** approach - systematically reducing complexity while maintaining 100%
functionality. This project was simplified using
the Merciless Simplification Methodology.

### Simplification Results Achieved

- **Test files**: 8 → 6 files (25% reduction)
- **Test lines**: 1,075 → 843 lines (22% reduction)
- **Benchmark complexity**: 251 → 87 lines (65% reduction)
- **Documentation**: 31% reduction with improved clarity
- **Functionality**: 100% preserved with improved maintainability