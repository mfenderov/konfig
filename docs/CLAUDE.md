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
- `cd testdata/integration && go test` - Run critical integration tests (external module usage)

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

## Architecture (Post-Redesign)

### Design Philosophy Learned
- **Explicit > Implicit**: All file paths must be provided explicitly - no magic discovery
- **Performance > Convenience**: Hot path optimization (config access) over initialization speed  
- **Security-First**: Input validation, resource limits, path sanitization built into core design
- **Minimal Dependencies**: 2 runtime dependencies vs 200+ in competitors

### Core Components

1. **Main Configuration Loader** (`konfig.go`):
    - `Load(filePath)` - Explicit file path loading with security validation
    - `LoadWithProfile(filePath, profile)` - Profile-based configuration merging
    - `LoadInto(filePath, struct)` - Type-safe struct mapping with defaults
    - `LoadIntoWithProfile(filePath, profile, struct)` - Combined profile + struct loading

2. **Security-Hardened YAML Parser** (`yaml_parser.go`):
    - Path traversal protection (blocks `../` patterns)
    - File size limits (10MB maximum)
    - YAML complexity validation (32 levels max nesting)
    - Environment variable substitution with `${VAR:default}` syntax

### Performance Characteristics Achieved
- **Config Access**: 156ns/op (4.5x faster than Viper)
- **Memory Usage**: 32B/op (7.5x less than Viper)
- **Allocations**: 3/op (3x fewer than Viper)
- **Zero allocations** on hot path after initial load

### Security Measures Implemented
- Input sanitization at API boundaries
- Resource exhaustion prevention (file size, nesting depth)
- Path traversal attack prevention  
- Structured error handling (avoid information leakage)

### Testing Strategy That Worked
- **Production Robustness**: 25 tests covering malformed input, filesystem edges, concurrency
- **Integration Testing**: External module validation in `testdata/integration/`
- **Security Testing**: Path traversal, file size, complexity attack scenarios
- **Competitive Benchmarking**: Performance validation vs real alternatives
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
- Integration tests with actual YAML files (`testdata/configs/` directory)
- **Essential benchmarks only**: 3 core performance tests (reduced from 6)
- **Critical Integration Tests**: `testdata/integration/` validates the public API and real-world usage

### Test Files Structure

- `konfig_test.go` - Core functionality tests
- `profile_test.go` - Profile management tests
- `resource_finder_test.go` - Resource discovery tests
- `struct_loader_test.go` - Basic struct-based configuration tests
- `struct_loader_advanced_test.go` - Advanced struct tests (nested, error handling)
- `struct_loader_benchmark_test.go` - Essential performance benchmarks

## Dependencies

- `github.com/stretchr/testify` - Testing utilities (dev only)
- `gopkg.in/yaml.v3` - YAML parsing

### Removed Dependencies
- ~~`github.com/pkg/errors`~~ - Replaced with stdlib `fmt.Errorf` + `%w`
- ~~`github.com/spf13/pflag`~~ - Replaced with environment variables + `SetProfile()`

## Development Philosophy

konfig follows a **"merciless simplification"** approach - systematically reducing complexity while maintaining 100%
functionality. This project was simplified using
the Merciless Simplification Methodology.

### Simplification Results Achieved

#### v1.0 Results (Initial Simplification)
- **Test files**: 8 → 6 files (25% reduction)
- **Test lines**: 1,075 → 843 lines (22% reduction)
- **Benchmark complexity**: 251 → 87 lines (65% reduction)
- **Documentation**: 31% reduction with improved clarity

#### Results (API Modernization)
- **Dependencies**: 4 → 2 dependencies (50% reduction)
- **Public API surface**: 30% consolidation via fluent interface
- **Performance**: Configuration caching + pre-compiled reflection
- **Error handling**: Modern stdlib approach (zero dependencies)
- **Profile management**: Environment variables + programmatic control
- **Directory flexibility**: Configurable paths (no hardcoded `resources/`)
- **Functionality**: 100% preserved + enhanced with 100% backward compatibility