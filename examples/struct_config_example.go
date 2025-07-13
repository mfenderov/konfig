package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mfenderov/konfig"
)

// Example demonstrating comprehensive struct-based configuration with konfig

// ExampleDatabaseConfig represents database connection settings
type ExampleDatabaseConfig struct {
	Host            string `konfig:"host" default:"localhost"`
	Port            string `konfig:"port" default:"5432"`
	Name            string `konfig:"name" default:"myapp"`
	User            string `konfig:"user" default:"postgres"`
	Password        string `konfig:"password" default:"secret"`
	SSLMode         string `konfig:"ssl_mode" default:"disable"`
	MaxConnections  string `konfig:"max_connections" default:"100"`
	ConnMaxLifetime string `konfig:"conn_max_lifetime" default:"1h"`
}

// ExampleServerConfig represents HTTP server settings
type ExampleServerConfig struct {
	Host           string     `konfig:"host" default:"0.0.0.0"`
	Port           string     `konfig:"port" default:"8080"`
	ReadTimeout    string     `konfig:"read_timeout" default:"30s"`
	WriteTimeout   string     `konfig:"write_timeout" default:"30s"`
	MaxHeaderBytes string     `konfig:"max_header_bytes" default:"1048576"`
	TLS            TLSConfig  `konfig:"tls"`
	CORS           CORSConfig `konfig:"cors"`
}

// TLSConfig represents TLS/SSL settings
type TLSConfig struct {
	Enabled  string `konfig:"enabled" default:"false"`
	CertFile string `konfig:"cert_file" default:""`
	KeyFile  string `konfig:"key_file" default:""`
}

// CORSConfig represents CORS settings
type CORSConfig struct {
	Enabled        string `konfig:"enabled" default:"true"`
	AllowedOrigins string `konfig:"allowed_origins" default:"*"`
	AllowedMethods string `konfig:"allowed_methods" default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders string `konfig:"allowed_headers" default:"Content-Type,Authorization"`
}

// LoggingConfig represents logging settings
type LoggingConfig struct {
	Level      string `konfig:"level" default:"info"`
	Format     string `konfig:"format" default:"json"`
	Output     string `konfig:"output" default:"stdout"`
	MaxSize    string `konfig:"max_size" default:"100"`
	MaxBackups string `konfig:"max_backups" default:"3"`
	MaxAge     string `konfig:"max_age" default:"28"`
}

// SecurityConfig represents security settings
type SecurityConfig struct {
	JWT JWTConfig `konfig:"jwt"`
	API APIConfig `konfig:"api"`
}

// JWTConfig represents JWT token settings
type JWTConfig struct {
	Secret          string `konfig:"secret" default:"dev-secret-change-in-production"`
	ExpirationHours string `konfig:"expiration_hours" default:"24"`
	Issuer          string `konfig:"issuer" default:"myapp"`
}

// APIConfig represents API security settings
type APIConfig struct {
	RateLimit RateLimitConfig `konfig:"rate_limit"`
}

// RateLimitConfig represents rate limiting settings
type RateLimitConfig struct {
	RequestsPerMinute string `konfig:"requests_per_minute" default:"100"`
	Burst             string `konfig:"burst" default:"10"`
	Enabled           string `konfig:"enabled" default:"true"`
}

// RedisConfig represents Redis cache settings
type RedisConfig struct {
	Host     string `konfig:"host" default:"localhost"`
	Port     string `konfig:"port" default:"6379"`
	Password string `konfig:"password" default:""`
	DB       string `konfig:"db" default:"0"`
	Enabled  string `konfig:"enabled" default:"false"`
}

// MonitoringConfig represents monitoring and observability settings
type MonitoringConfig struct {
	Metrics MetricsConfig `konfig:"metrics"`
	Health  HealthConfig  `konfig:"health"`
	Tracing TracingConfig `konfig:"tracing"`
}

// MetricsConfig represents metrics collection settings
type MetricsConfig struct {
	Enabled bool   `konfig:"enabled" default:"true"`
	Port    string `konfig:"port" default:"9090"`
	Path    string `konfig:"path" default:"/metrics"`
}

// HealthConfig represents health check settings
type HealthConfig struct {
	Enabled bool   `konfig:"enabled" default:"true"`
	Path    string `konfig:"path" default:"/health"`
}

// TracingConfig represents distributed tracing settings
type TracingConfig struct {
	Enabled        bool   `konfig:"enabled" default:"false"`
	JaegerEndpoint string `konfig:"jaeger_endpoint" default:"http://localhost:14268/api/traces"`
	ServiceName    string `konfig:"service_name" default:"myapp"`
}

// AppConfig represents the complete application configuration
type AppConfig struct {
	Application ApplicationConfig `konfig:"application"`
	Server      ExampleServerConfig      `konfig:"server"`
	Database    ExampleDatabaseConfig    `konfig:"database"`
	Logging     LoggingConfig     `konfig:"logging"`
	Security    SecurityConfig    `konfig:"security"`
	Redis       RedisConfig       `konfig:"redis"`
	Monitoring  MonitoringConfig  `konfig:"monitoring"`
}

// ApplicationConfig represents basic application metadata
type ApplicationConfig struct {
	Name        string `konfig:"name" default:"myapp"`
	Version     string `konfig:"version" default:"1.0.0"`
	Environment string `konfig:"environment" default:"development"`
	Debug       string `konfig:"debug" default:"true"`
	Description string `konfig:"description" default:"My awesome application"`
}

func RunStructExample() {
	fmt.Println("🚀 Konfig Struct-Based Configuration Example")
	fmt.Println("============================================")

	// Demonstrate setting some environment variables to override defaults
	fmt.Println("\n📋 Setting up example environment variables...")
	envVars := map[string]string{
		"application.name":        "example-app",
		"application.version":     "2.0.0",
		"server.port":            "9090",
		"database.host":          "db.example.com",
		"database.port":          "3306",
		"security.jwt.secret":    "super-secret-jwt-key",
		"logging.level":          "debug",
		"redis.enabled":          "true",
		"redis.host":             "redis.example.com",
		"monitoring.metrics.enabled": "true",
		"monitoring.tracing.enabled": "true",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
		fmt.Printf("   %s = %s\n", key, value)
	}

	// Load configuration into struct
	fmt.Println("\n🔧 Loading configuration using LoadInto()...")
	var config AppConfig
	err := konfig.LoadInto(&config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Println("✅ Configuration loaded successfully!")

	// Display configuration sections
	fmt.Println("\n📖 Configuration Summary:")
	fmt.Println("========================")

	// Application info
	fmt.Printf("📱 Application: %s v%s\n", config.Application.Name, config.Application.Version)
	fmt.Printf("   Environment: %s\n", config.Application.Environment)
	fmt.Printf("   Debug: %s\n", config.Application.Debug)
	fmt.Printf("   Description: %s\n", config.Application.Description)

	// Server configuration
	fmt.Printf("\n🌐 Server Configuration:\n")
	fmt.Printf("   Address: %s:%s\n", config.Server.Host, config.Server.Port)
	fmt.Printf("   Read Timeout: %s\n", config.Server.ReadTimeout)
	fmt.Printf("   Write Timeout: %s\n", config.Server.WriteTimeout)
	fmt.Printf("   TLS Enabled: %s\n", config.Server.TLS.Enabled)
	fmt.Printf("   CORS Enabled: %s\n", config.Server.CORS.Enabled)
	fmt.Printf("   CORS Origins: %s\n", config.Server.CORS.AllowedOrigins)

	// Database configuration
	fmt.Printf("\n🗄️  Database Configuration:\n")
	fmt.Printf("   Connection: %s@%s:%s/%s\n", config.Database.User, config.Database.Host, config.Database.Port, config.Database.Name)
	fmt.Printf("   SSL Mode: %s\n", config.Database.SSLMode)
	fmt.Printf("   Max Connections: %s\n", config.Database.MaxConnections)

	// Security configuration
	fmt.Printf("\n🔐 Security Configuration:\n")
	fmt.Printf("   JWT Secret: %s\n", maskSecret(config.Security.JWT.Secret))
	fmt.Printf("   JWT Expiration: %s hours\n", config.Security.JWT.ExpirationHours)
	fmt.Printf("   Rate Limit: %s req/min (burst: %s)\n", 
		config.Security.API.RateLimit.RequestsPerMinute, 
		config.Security.API.RateLimit.Burst)

	// Logging configuration
	fmt.Printf("\n📝 Logging Configuration:\n")
	fmt.Printf("   Level: %s\n", config.Logging.Level)
	fmt.Printf("   Format: %s\n", config.Logging.Format)
	fmt.Printf("   Output: %s\n", config.Logging.Output)

	// Redis configuration
	fmt.Printf("\n📦 Redis Configuration:\n")
	fmt.Printf("   Enabled: %s\n", config.Redis.Enabled)
	if config.Redis.Enabled == "true" {
		fmt.Printf("   Connection: %s:%s (DB: %s)\n", config.Redis.Host, config.Redis.Port, config.Redis.DB)
	}

	// Monitoring configuration
	fmt.Printf("\n📊 Monitoring Configuration:\n")
	fmt.Printf("   Metrics Enabled: %t (Port: %s)\n", config.Monitoring.Metrics.Enabled, config.Monitoring.Metrics.Port)
	fmt.Printf("   Health Checks: %t (Path: %s)\n", config.Monitoring.Health.Enabled, config.Monitoring.Health.Path)
	fmt.Printf("   Tracing Enabled: %t\n", config.Monitoring.Tracing.Enabled)
	if config.Monitoring.Tracing.Enabled {
		fmt.Printf("   Jaeger Endpoint: %s\n", config.Monitoring.Tracing.JaegerEndpoint)
	}

	// Profile-aware logic
	fmt.Printf("\n🏷️  Profile Information:\n")
	profile := konfig.GetProfile()
	if profile != "" {
		fmt.Printf("   Active Profile: %s\n", profile)
	} else {
		fmt.Printf("   Active Profile: default\n")
	}

	if konfig.IsDevProfile() {
		fmt.Println("   🔧 Running in development mode")
		fmt.Println("   - Enhanced logging enabled")
		fmt.Println("   - CORS allows all origins")
		fmt.Println("   - Debug mode active")
	} else if konfig.IsProdProfile() {
		fmt.Println("   🚀 Running in production mode")
		fmt.Println("   - Security hardening active")
		fmt.Println("   - Performance optimizations enabled")
		fmt.Println("   - Monitoring and alerting active")
	}

	// Validation example
	fmt.Printf("\n✅ Configuration Validation:\n")
	validateConfiguration(&config)

	fmt.Println("\n🎉 Example completed successfully!")
	fmt.Println("This demonstrates how konfig can load complex nested configurations")
	fmt.Println("with defaults, environment variable overrides, and type safety.")

	// Clean up environment variables
	for key := range envVars {
		os.Unsetenv(key)
	}
}

// maskSecret masks sensitive information for display
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

// validateConfiguration demonstrates how to add validation to your configuration
func validateConfiguration(config *AppConfig) {
	issues := []string{}

	// Validate required fields
	if config.Application.Name == "" {
		issues = append(issues, "Application name is required")
	}

	// Validate security settings
	if konfig.IsProdProfile() && config.Security.JWT.Secret == "dev-secret-change-in-production" {
		issues = append(issues, "Production JWT secret must be changed from default")
	}

	// Validate database settings
	if config.Database.Host == "" {
		issues = append(issues, "Database host is required")
	}

	// Report validation results
	if len(issues) > 0 {
		fmt.Printf("   ⚠️  Configuration Issues Found:\n")
		for _, issue := range issues {
			fmt.Printf("      - %s\n", issue)
		}
	} else {
		fmt.Printf("   ✅ All validations passed\n")
	}
}