package konfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// BenchmarkLoad_SmallConfig benchmarks loading a typical small configuration
func BenchmarkLoad_SmallConfig(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "small.yaml")

	content := `
server:
  port: 8080
  host: localhost
database:
  url: postgres://localhost/test
  pool_size: 10
logging:
  level: info
  format: json
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg, err := Load(configPath)
		if err != nil {
			b.Fatal(err)
		}
		_ = cfg.GetString("server.port")
	}
}

// BenchmarkLoad_LargeConfig benchmarks loading a large configuration (1000+ keys)
func BenchmarkLoad_LargeConfig(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "large.yaml")

	// Generate large configuration
	var content strings.Builder
	content.WriteString("services:\n")
	for i := 0; i < 500; i++ {
		content.WriteString(fmt.Sprintf("  service_%d:\n", i))
		content.WriteString(fmt.Sprintf("    name: service-%d\n", i))
		content.WriteString(fmt.Sprintf("    port: %d\n", 8000+i))
		content.WriteString(fmt.Sprintf("    enabled: %t\n", i%2 == 0))
	}

	err := os.WriteFile(configPath, []byte(content.String()), 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg, err := Load(configPath)
		if err != nil {
			b.Fatal(err)
		}
		_ = cfg.GetString("services.service_0.name")
	}
}

// BenchmarkLoadWithProfile benchmarks profile-based loading
func BenchmarkLoadWithProfile(b *testing.B) {
	tempDir := b.TempDir()

	// Base config
	baseConfigPath := filepath.Join(tempDir, "app.yaml")
	baseConfig := `
server:
  port: 8080
  host: localhost
database:
  url: postgres://localhost/test
`
	err := os.WriteFile(baseConfigPath, []byte(baseConfig), 0644)
	if err != nil {
		b.Fatal(err)
	}

	// Profile config
	profileConfigPath := filepath.Join(tempDir, "app-prod.yaml")
	profileConfig := `
server:
  port: 443
  host: prod.example.com
database:
  url: postgres://prod-db.example.com/app
  ssl_mode: require
`
	err = os.WriteFile(profileConfigPath, []byte(profileConfig), 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg, err := LoadWithProfile(baseConfigPath, "prod")
		if err != nil {
			b.Fatal(err)
		}
		_ = cfg.GetString("server.host")
	}
}

// BenchmarkLoadInto_SmallStruct benchmarks struct loading with a small struct
func BenchmarkLoadInto_SmallStruct(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "struct.yaml")

	content := `
server:
  port: 8080
  host: localhost
  debug: true
timeout: 30s
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		b.Fatal(err)
	}

	type ServerConfig struct {
		Port  int    `konfig:"port"`
		Host  string `konfig:"host"`
		Debug bool   `konfig:"debug"`
	}

	type Config struct {
		Server  ServerConfig `konfig:"server"`
		Timeout string       `konfig:"timeout"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cfg Config
		err := LoadInto(configPath, &cfg)
		if err != nil {
			b.Fatal(err)
		}
		_ = cfg.Server.Port
	}
}

// BenchmarkLoadInto_LargeStruct benchmarks struct loading with a complex struct
func BenchmarkLoadInto_LargeStruct(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "complex.yaml")

	content := `
service:
  name: my-service
  version: 1.0.0
  port: 8080

server:
  host: localhost
  port: 8080
  timeout: 30s
  max_connections: 1000

database:
  host: localhost
  port: 5432
  name: myapp
  username: user
  password: pass
  ssl_mode: require
  pool_size: 10
  max_idle: 5

cache:
  enabled: true
  host: localhost
  port: 6379
  db: 0
  ttl: 300

logging:
  level: info
  format: json
  file: /var/log/app.log

metrics:
  enabled: true
  port: 9090
  path: /metrics

health:
  enabled: true
  port: 8081
  path: /health
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		b.Fatal(err)
	}

	type ServiceConfig struct {
		Name    string `konfig:"name"`
		Version string `konfig:"version"`
		Port    int    `konfig:"port"`
	}

	type ServerConfig struct {
		Host           string `konfig:"host"`
		Port           int    `konfig:"port"`
		Timeout        string `konfig:"timeout"`
		MaxConnections int    `konfig:"max_connections"`
	}

	type DatabaseConfig struct {
		Host     string `konfig:"host"`
		Port     int    `konfig:"port"`
		Name     string `konfig:"name"`
		Username string `konfig:"username"`
		Password string `konfig:"password"`
		SSLMode  string `konfig:"ssl_mode"`
		PoolSize int    `konfig:"pool_size"`
		MaxIdle  int    `konfig:"max_idle"`
	}

	type CacheConfig struct {
		Enabled bool   `konfig:"enabled"`
		Host    string `konfig:"host"`
		Port    int    `konfig:"port"`
		DB      int    `konfig:"db"`
		TTL     int    `konfig:"ttl"`
	}

	type LoggingConfig struct {
		Level  string `konfig:"level"`
		Format string `konfig:"format"`
		File   string `konfig:"file"`
	}

	type MetricsConfig struct {
		Enabled bool   `konfig:"enabled"`
		Port    int    `konfig:"port"`
		Path    string `konfig:"path"`
	}

	type HealthConfig struct {
		Enabled bool   `konfig:"enabled"`
		Port    int    `konfig:"port"`
		Path    string `konfig:"path"`
	}

	type Config struct {
		Service  ServiceConfig  `konfig:"service"`
		Server   ServerConfig   `konfig:"server"`
		Database DatabaseConfig `konfig:"database"`
		Cache    CacheConfig    `konfig:"cache"`
		Logging  LoggingConfig  `konfig:"logging"`
		Metrics  MetricsConfig  `konfig:"metrics"`
		Health   HealthConfig   `konfig:"health"`
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cfg Config
		err := LoadInto(configPath, &cfg)
		if err != nil {
			b.Fatal(err)
		}
		_ = cfg.Database.PoolSize
	}
}

// BenchmarkConfigAccess benchmarks accessing configuration values
func BenchmarkConfigAccess(b *testing.B) {
	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "access.yaml")

	content := `
server:
  port: 8080
  host: localhost
database:
  url: postgres://localhost/test
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		b.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	b.Run("GetString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cfg.GetString("server.host")
		}
	})

	b.Run("GetInt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = cfg.GetInt("server.port")
		}
	})

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = cfg.Get("server.port")
		}
	})
}

// BenchmarkEnvSubstitution benchmarks environment variable substitution
func BenchmarkEnvSubstitution(b *testing.B) {
	// Set test environment variables
	os.Setenv("TEST_HOST", "benchtest.example.com")
	os.Setenv("TEST_PORT", "9999")
	defer func() {
		os.Unsetenv("TEST_HOST")
		os.Unsetenv("TEST_PORT")
	}()

	tempDir := b.TempDir()
	configPath := filepath.Join(tempDir, "env.yaml")

	content := `
server:
  host: ${TEST_HOST:localhost}
  port: ${TEST_PORT:8080}
  protocol: ${UNDEFINED_VAR:http}
database:
  url: postgres://${TEST_HOST:localhost}/app
  ssl: ${SSL_ENABLED:true}
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg, err := Load(configPath)
		if err != nil {
			b.Fatal(err)
		}
		_ = cfg.GetString("server.host")
	}
}
