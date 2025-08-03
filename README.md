# konfig

> **Production-ready configuration management for Go** - Explicit, type-safe, and blazingly fast.

[![Go Reference](https://pkg.go.dev/badge/github.com/mfenderov/konfig.svg)](https://pkg.go.dev/github.com/mfenderov/konfig)
[![Go Report Card](https://goreportcard.com/badge/github.com/mfenderov/konfig)](https://goreportcard.com/report/github.com/mfenderov/konfig)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 🚀 **Why konfig?**

**Best-in-class Go configuration library** with zero magic, explicit behavior, and production-proven reliability:

✅ **Explicit Paths** - No magic file discovery, no surprises in production  
✅ **Type-Safe** - Struct-based loading with compile-time validation  
✅ **Profile Support** - Environment-specific configs (dev/prod/test)  
✅ **Environment Variables** - `${VAR:default}` substitution built-in  
✅ **Zero Globals** - No init() side effects or global state  
✅ **Production Ready** - Handles edge cases, concurrent access, large configs  
✅ **Lightning Fast** - 14ns config access, 0 allocations on hot paths

## ⚡ Quick Start

```bash
go get github.com/mfenderov/konfig
```

**Simple Configuration Loading:**

```go
// Load configuration from explicit path
cfg, err := konfig.Load("./config/app.yaml")
if err != nil {
    log.Fatal(err)
}

port := cfg.GetString("server.port")
debug := cfg.GetBool("app.debug")
```

**Type-Safe Struct Loading:**

```go
type Config struct {
    Server struct {
        Port int    `konfig:"port" default:"8080"`
        Host string `konfig:"host" default:"localhost"`
    } `konfig:"server"`
    Database struct {
        URL      string `konfig:"url"`
        PoolSize int    `konfig:"pool_size" default:"10"`
    } `konfig:"database"`
}

var cfg Config
err := konfig.LoadInto("./config/app.yaml", &cfg)
```

**Profile-Based Configuration:**

```go
// Load base config + profile-specific overrides
cfg, err := konfig.LoadWithProfile("./config/app.yaml", "prod")

// Or with structs
var cfg Config
err := konfig.LoadIntoWithProfile("./config/app.yaml", "prod", &cfg)
```

## 📁 Configuration Structure

**Base Configuration** (`config/app.yaml`):
```yaml
server:
  port: 8080
  host: localhost

database:
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
  name: myapp
```

**Profile Configuration** (`config/app-prod.yaml`):
```yaml
server:
  host: 0.0.0.0
  port: 443

database:
  host: prod-db.example.com
  ssl_mode: require
```

Profile values automatically override base values when loaded with `LoadWithProfile()`.

## 🏗️ Real-World Example

```go
package main

import (
    "log"
    "github.com/mfenderov/konfig"
)

type AppConfig struct {
    Service struct {
        Name    string `konfig:"name"`
        Version string `konfig:"version"`
        Port    int    `konfig:"port" default:"8080"`
    } `konfig:"service"`
    
    Database struct {
        Host     string `konfig:"host" default:"localhost"`
        Port     int    `konfig:"port" default:"5432"`
        Name     string `konfig:"name"`
        Username string `konfig:"username"`
        Password string `konfig:"password"`
        SSLMode  string `konfig:"ssl_mode" default:"disable"`
    } `konfig:"database"`
    
    Cache struct {
        Enabled bool   `konfig:"enabled" default:"false"`
        URL     string `konfig:"url" default:"redis://localhost:6379"`
        TTL     int    `konfig:"ttl" default:"300"`
    } `konfig:"cache"`
}

func main() {
    var cfg AppConfig
    
    // Load configuration with profile support
    err := konfig.LoadIntoWithProfile("./config/app.yaml", "prod", &cfg)
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }
    
    // Use type-safe configuration
    log.Printf("Starting %s v%s on port %d", 
        cfg.Service.Name, 
        cfg.Service.Version, 
        cfg.Service.Port)
        
    // Database connection
    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
        cfg.Database.Username,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Name,
        cfg.Database.SSLMode)
}
```

## 🔧 Environment Variables

konfig supports environment variable substitution with defaults:

```yaml
database:
  host: ${DB_HOST:localhost}           # Use DB_HOST env var, fallback to localhost
  port: ${DB_PORT:5432}               # Use DB_PORT env var, fallback to 5432
  password: ${DB_PASSWORD}            # Use DB_PASSWORD env var, no fallback
  ssl_mode: ${DB_SSL_MODE:require}    # Use DB_SSL_MODE env var, fallback to require
```

## 🚀 Performance

konfig is optimized for production use with excellent performance characteristics:

| Operation | Speed | Memory | Notes |
|-----------|-------|---------|-------|
| Small config loading | 38μs | 19KB | Typical microservice config |
| Large config loading | 3.2ms | 2MB | 1000+ configuration keys |
| Config value access | 14ns | 0 allocs | Hot path optimized |
| Type conversion | 45ns | 16B | String/int/bool conversion |
| Struct loading | 35μs | 16KB | Type-safe configuration |

**Benchmark Results:**
```
BenchmarkLoad_SmallConfig         29971    38550 ns/op   19428 B/op    234 allocs/op
BenchmarkConfigAccess/Get         83567641    14.45 ns/op     0 B/op      0 allocs/op
BenchmarkLoadInto_SmallStruct     35079    34828 ns/op   16118 B/op    179 allocs/op
```

## 🛡️ Production Ready

konfig is battle-tested with comprehensive edge case handling:

- ✅ **Concurrent Access** - Thread-safe by design
- ✅ **Large Configurations** - Handles 1000+ keys efficiently  
- ✅ **Malformed YAML** - Graceful error handling with context
- ✅ **File System Edge Cases** - Permissions, deep paths, readonly files
- ✅ **Memory Management** - No memory leaks, bounded allocations
- ✅ **Error Recovery** - Detailed error messages for debugging

## 🔄 Migration Guide

konfig removes architectural flaws while maintaining simplicity:

**❌ What was removed (breaking changes):**
- Project root auto-discovery (fragile in production)
- Global logger setup in init() (side effects)
- Environment variable pollution (global state)
- Overlapping legacy APIs (confusion)

**✅ What's new:**
- Explicit file paths only (predictable)
- Zero global state (safe)
- Structured errors with context (debuggable)
- Production-proven reliability (robust)

**Migration is simple:**
```go
// Old (v1) - magic behavior
err := konfig.Load()
value := os.Getenv("some.key")

// New - explicit behavior  
cfg, err := konfig.Load("./config/app.yaml")
value := cfg.GetString("some.key")
```

## 📚 API Reference

### Loading Functions

```go
// Load single configuration file
func Load(filePath string) (Config, error)

// Load with profile support (base + profile files)
func LoadWithProfile(filePath, profile string) (Config, error)

// Load into struct (type-safe)
func LoadInto(filePath string, target interface{}) error

// Load into struct with profile support
func LoadIntoWithProfile(filePath, profile string, target interface{}) error
```

### Config Interface

```go
type Config interface {
    // Get raw value
    Get(key string) (interface{}, bool)
    
    // Type-safe getters
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    GetFloat64(key string) float64
    GetDuration(key string) time.Duration
    
    // Getters with defaults
    GetStringWithDefault(key, defaultValue string) string
    GetIntWithDefault(key string, defaultValue int) int
    GetBoolWithDefault(key string, defaultValue bool) bool
    
    // Introspection
    Keys() []string
}
```

### Struct Tags

```go
type Config struct {
    Port    int    `konfig:"server.port" default:"8080"`
    Host    string `konfig:"server.host" default:"localhost"`
    Timeout string `konfig:"timeout" default:"30s"`
    Debug   bool   `konfig:"debug" default:"false"`
}
```

## 🧪 Testing

konfig includes comprehensive test coverage:

```bash
go test ./...                 # Run all tests
go test -bench=.             # Run performance benchmarks  
go test -race ./...          # Test concurrent safety
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## ⭐ Star History

If konfig helps you build better Go applications, please consider giving it a star! 

---

**Built with ❤️ for the Go community**