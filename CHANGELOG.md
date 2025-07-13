# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive godoc documentation for all exported functions
- Development tooling: golangci-lint, test coverage reporting, quality checks
- Pre-commit hooks configuration
- Organized test suite with focused test files
- Go workspace support with proper configuration
- Enhanced Makefile with quality, coverage, and CI targets
- CLAUDE.md for AI assistant guidance

### Changed  
- Improved project organization with better test file structure
- Updated CI/CD workflows to use quality pipeline
- Cleaned up .gitignore to be Go-specific
- Removed vendor directory in favor of Go modules

### Fixed
- Package name conflicts in examples directory
- Compilation issues with Go workspace mode
- Makefile compatibility with workspace configuration

## [0.17.0] - Previous Release

### Added
- Struct-based configuration system with `LoadInto()` function
- Support for nested structs with automatic prefix handling
- Default value support using struct tags
- Comprehensive test suite for struct loading functionality
- Benchmark tests for performance measurement

### Changed
- Enhanced configuration loading with type safety
- Improved error handling and validation

## [Earlier Versions]

For earlier version history, please refer to the git commit history or GitHub releases.

---

## Release Guidelines

### Version Types
- **Major** (X.0.0): Breaking changes to public API
- **Minor** (0.X.0): New features, backward compatible
- **Patch** (0.0.X): Bug fixes, backward compatible

### Adding New Entries
When adding new changes:
1. Add entries under `[Unreleased]` section
2. Use categories: `Added`, `Changed`, `Deprecated`, `Removed`, `Fixed`, `Security`
3. Write clear, concise descriptions
4. Include references to issues/PRs when applicable

### Creating a Release
1. Move entries from `[Unreleased]` to new version section
2. Add release date in format `[X.Y.Z] - YYYY-MM-DD`
3. Update version references in code
4. Create git tag and GitHub release