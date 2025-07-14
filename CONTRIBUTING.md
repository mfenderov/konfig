# Contributing to konfig

Thank you for your interest in contributing to konfig! This document provides guidelines and information for contributors.

## 🚀 Quick Start

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/konfig.git
   cd konfig
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. **Make your changes** with tests
5. **Submit a Pull Request**

## 🧪 Development Setup

### Prerequisites
- Go 1.24 or later
- golangci-lint (optional, for linting)

### Running Tests
```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run linting (informational)
make lint

# Run all quality checks
make quality
```

### Project Structure
```
konfig/
├── konfig.go              # Main configuration loader
├── struct_loader.go       # Struct-based configuration
├── profile.go             # Profile management
├── resource_finder.go     # Configuration file discovery
├── loader.go              # YAML parsing and loading
├── *_test.go              # Test files
├── examples/              # Usage examples
├── test-proj/             # Integration tests
└── resources/             # Test configuration files
```

## 📋 Contribution Guidelines

### Code Style
- Follow Go conventions and best practices
- Use `gofmt` to format your code
- Add comments for exported functions and types
- Keep functions focused and small

### Testing
- **All new code must include tests**
- Aim for high test coverage (>90%)
- Test both success and error cases
- Include integration tests for new features

### Documentation
- Update README.md for new features or API changes
- Add examples for new functionality
- Update struct tag documentation for new tags
- Keep documentation concise and practical

### Commit Messages
Use clear, descriptive commit messages:
```
feat: add support for custom configuration directories
fix: handle empty environment variable defaults correctly
docs: update struct tag examples in README
test: add integration tests for profile loading
```

## 🎯 Areas for Contribution

### High Priority
- **Performance optimizations** - faster configuration loading
- **Additional struct tag features** - validation, transformation
- **Better error messages** - more helpful error reporting
- **Configuration validation** - built-in validation helpers

### Medium Priority
- **Additional file formats** - JSON, TOML support
- **Configuration watching** - reload on file changes
- **Environment-specific optimizations** - production vs development
- **Documentation improvements** - more examples, tutorials

### Low Priority
- **CLI tools** - configuration validation utilities
- **IDE integrations** - better development experience
- **Additional profile features** - profile inheritance, merging

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Go version** and operating system
2. **konfig version** or commit hash
3. **Minimal reproduction case** - ideally a failing test
4. **Expected vs actual behavior**
5. **Configuration files** (if relevant)
6. **Error messages** or logs

### Bug Report Template
```markdown
## Bug Description
[Clear description of the bug]

## Reproduction Steps
1. [Step 1]
2. [Step 2]
3. [Step 3]

## Expected Behavior
[What should happen]

## Actual Behavior
[What actually happens]

## Environment
- Go version: [e.g., 1.24]
- OS: [e.g., macOS 14.0]
- konfig version: [e.g., v0.17.0]

## Additional Context
[Any other relevant information]
```

## ✨ Feature Requests

For feature requests, please:

1. **Check existing issues** to avoid duplicates
2. **Describe the use case** - why is this needed?
3. **Propose an API** - how should it work?
4. **Consider backwards compatibility** - will this break existing code?
5. **Offer to implement** - are you willing to contribute the code?

### Feature Request Template
```markdown
## Feature Description
[Clear description of the desired feature]

## Use Case
[Why is this feature needed? What problem does it solve?]

## Proposed API
[How should the feature work? Include code examples]

## Backwards Compatibility
[Will this change break existing code?]

## Implementation Notes
[Any technical considerations or constraints]
```

## 🔍 Code Review Process

### For Contributors
- **Keep PRs focused** - one feature or fix per PR
- **Write good PR descriptions** - explain what and why
- **Be responsive** to review feedback
- **Update documentation** as needed
- **Ensure CI passes** before requesting review

### Review Criteria
- **Functionality** - does it work as intended?
- **Testing** - is it well tested?
- **Performance** - no significant regressions
- **API design** - is it consistent and intuitive?
- **Documentation** - is it properly documented?
- **Backwards compatibility** - does it break existing code?

## 📚 Project Philosophy

konfig follows a **"merciless simplification"** approach:

- **Keep what provides value, eliminate what doesn't**
- **Every feature must solve a real user problem**
- **Performance and simplicity over complex abstractions**
- **Clear, focused APIs over extensive configuration options**
- **Practical examples over comprehensive documentation**

### Design Principles
1. **Simple by default** - common use cases should be easy
2. **Extensible when needed** - support advanced scenarios
3. **Fast and lightweight** - minimal performance overhead
4. **Type-safe** - leverage Go's type system for safety
5. **Predictable** - consistent behavior across features

## 🎉 Recognition

Contributors are recognized in several ways:

- **GitHub contributor list** - automatic recognition
- **Release notes** - significant contributions mentioned
- **Community acknowledgment** - featured in discussions

### Types of Contributions
- **Code contributions** - new features, bug fixes, optimizations
- **Documentation** - README updates, examples, tutorials
- **Testing** - additional test cases, CI improvements
- **Issue triage** - helping with bug reports and questions
- **Community support** - answering questions, helping users

## 📞 Getting Help

- **GitHub Issues** - for bugs and feature requests
- **GitHub Discussions** - for questions and general discussion
- **Code Review** - detailed feedback on PRs

## 🚀 Release Process

1. **Semantic versioning** - following [semver](https://semver.org/)
2. **Changelog maintenance** - all changes documented
3. **Backwards compatibility** - breaking changes in major versions only
4. **Testing** - comprehensive testing before release
5. **Documentation** - updated for new features

Thank you for contributing to konfig! Every contribution helps make configuration management better for the Go community. 🎉