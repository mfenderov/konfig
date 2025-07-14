# konfig

> **Spring-inspired configuration management for Go** - Simple, type-safe, and blazingly fast.

[![Go Reference](https://pkg.go.dev/badge/github.com/mfenderov/konfig.svg)](https://pkg.go.dev/github.com/mfenderov/konfig)
[![Go Report Card](https://goreportcard.com/badge/github.com/mfenderov/konfig)](https://goreportcard.com/report/github.com/mfenderov/konfig)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Simplified](https://img.shields.io/badge/complexity-31%25%20reduced-green.svg)](#project-quality)

## ⚡ Quick Start

```bash
go get github.com/mfenderov/konfig
```

**1. Create your config struct:**
```go
type Config struct {
    Server struct {
        Port string `konfig:"port" default:"8080"`
        Host string `konfig:"host" default:"localhost"`
    } `konfig:"server"`
    Database struct {
        URL  string `konfig:"url" default:"postgres://localhost/myapp"`
        Pool int    `konfig:"pool" default:"10"`
    } `konfig:"database"`
}
```

**2. Load configuration:**
```go
var config Config
err := konfig.LoadInto(&config)
// Done! Your struct is populated with type-safe configuration
```

**3. Use with profiles:**
```bash
go run app.go -p prod  # Loads application-prod.yaml
```

## 🎯 Why konfig?

| Feature | konfig | Others |
|---------|--------|--------|
| **Type Safety** | ✅ Compile-time struct validation | ❌ Runtime string lookups |
| **Profile Support** | ✅ Built-in dev/prod/test profiles | ❌ Manual environment handling |
| **Zero Dependencies** | ✅ Just stdlib + yaml | ❌ Heavy dependency chains |
| **Spring-like** | ✅ Familiar `application.yaml` pattern | ❌ Custom configuration styles |
| **Env Variable Integration** | ✅ Automatic substitution & overrides | ❌ Manual env handling |

## 📁 Configuration Files

```yaml
# application.yaml (default)
server:
  port: ${PORT:8080}
  host: ${HOST:0.0.0.0}

database:
  url: ${DATABASE_URL:postgres://localhost/myapp}
  pool: ${DB_POOL:10}
```

```yaml
# application-prod.yaml (production overrides)
server:
  port: 443
database:
  pool: 50
```

## 🚀 Advanced Features

<details>
<summary><strong>Environment Variable Substitution</strong></summary>

```yaml
database:
  url: ${DATABASE_URL:postgres://localhost/default}  # Uses env var or default
  password: ${DB_PASSWORD}  # Required env var (fails if missing)
```
</details>

<details>
<summary><strong>Nested Configuration</strong></summary>

```go
type Config struct {
    App struct {
        Name    string `konfig:"name" default:"MyApp"`
        Version string `konfig:"version" default:"1.0.0"`
        Features struct {
            Auth     bool `konfig:"auth" default:"true"`
            Metrics  bool `konfig:"metrics" default:"false"`
        } `konfig:"features"`
    } `konfig:"application"`
}
```
</details>

<details>
<summary><strong>Profile-Aware Code</strong></summary>

```go
if konfig.IsProdProfile() {
    // Production-specific logic
    enableHTTPS()
} else if konfig.IsDevProfile() {
    // Development helpers
    enableDebugMode()
}
```
</details>

## 📚 API Reference

### Core Functions
```go
konfig.LoadInto(&config)     // Load into struct (recommended)
konfig.Load()                // Load into environment variables
konfig.LoadFrom("file.yaml") // Load specific file

// Profile helpers
konfig.GetProfile()          // Current active profile
konfig.IsDevProfile()        // Check if dev profile
konfig.IsProdProfile()       // Check if prod profile
```

### Struct Tags
```go
type Config struct {
    Port string `konfig:"server.port" default:"8080"`
    //            ^^^^^^^^^^^^^^^^^^^^^  ^^^^^^^^^^^^^
    //            Configuration key      Default value
}
```

## 🔍 Examples

See [`examples/`](examples/) for complete working examples:

```bash
cd examples
go run simple_example.go                              # Basic usage
go run simple_example.go -p dev                      # With dev profile  
go run simple_example.go struct_config_example.go    # Advanced features
```

## 🏆 Project Quality

konfig maintains high standards through **systematic simplification**:

- **🧪 100% test coverage** - 6 focused test files covering all functionality
- **⚡ Performance optimized** - 3 essential benchmarks measuring real scenarios
- **📖 Clear documentation** - Concise, example-driven approach
- **🎯 31% complexity reduction** - Achieved through [merciless simplification](CLAUDE.md#simplification-results-achieved)
- **✅ Zero functionality loss** - All features preserved during optimization

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes with tests
4. Submit a Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Made with ❤️ for the Go community** | [Documentation](https://pkg.go.dev/github.com/mfenderov/konfig) | [Examples](examples/) | [Contributing](CONTRIBUTING.md)