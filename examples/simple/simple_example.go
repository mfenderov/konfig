package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mfenderov/konfig"
)

// Simple example demonstrating the new explicit konfig API

// DatabaseConfig represents database settings
type DatabaseConfig struct {
	Host     string `konfig:"host" default:"localhost"`
	Port     int    `konfig:"port" default:"5432"`
	Name     string `konfig:"name" default:"myapp"`
	Username string `konfig:"username" default:"postgres"`
	Password string `konfig:"password" default:"secret"`
}

// ServerConfig represents server settings
type ServerConfig struct {
	Host string `konfig:"host" default:"0.0.0.0"`
	Port int    `konfig:"port" default:"8080"`
}

// SimpleConfig represents the complete application configuration
type SimpleConfig struct {
	AppName  string         `konfig:"name" default:"my-app"`
	Version  string         `konfig:"version" default:"1.0.0"`
	Debug    bool           `konfig:"debug" default:"true"`
	Database DatabaseConfig `konfig:"database"`
	Server   ServerConfig   `konfig:"server"`
}

func main() {
	fmt.Println("🚀 Simple Konfig Example")
	fmt.Println("========================")

	// Set some environment variables to demonstrate substitution
	os.Setenv("DB_HOST", "prod-db.example.com")
	os.Setenv("APP_PORT", "9090")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("APP_PORT")
	}()

	// Create example configuration file
	configPath := "./example-config.yaml"
	configContent := `
name: simple-example
version: 2.0.0
debug: false

server:
  host: localhost
  port: ${APP_PORT:8080}

database:
  host: ${DB_HOST:localhost}
  port: 5432
  name: example_db
  user: admin
  password: ${DB_PASSWORD:secret}
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		log.Fatal("Failed to create example config:", err)
	}
	defer os.Remove(configPath) // Clean up

	// Example 1: Direct configuration loading
	fmt.Println("\n📖 Example 1: Direct Configuration Loading")
	cfg, err := konfig.Load(configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	fmt.Printf("   App Name: %s\n", cfg.GetString("name"))
	fmt.Printf("   Server Port: %d\n", cfg.GetInt("server.port"))
	fmt.Printf("   Database Host: %s\n", cfg.GetString("database.host"))

	// Example 2: Type-safe struct loading
	fmt.Println("\n📖 Example 2: Type-Safe Struct Loading")
	var config SimpleConfig
	err = konfig.LoadInto(configPath, &config)
	if err != nil {
		log.Fatal("Failed to load into struct:", err)
	}

	fmt.Printf("   Service: %s v%s (Debug: %t)\n",
		config.AppName, config.Version, config.Debug)
	fmt.Printf("   Server: %s:%d\n",
		config.Server.Host, config.Server.Port)
	fmt.Printf("   Database: %s@%s:%d/%s\n",
		config.Database.Username, config.Database.Host,
		config.Database.Port, config.Database.Name)

	// Example 3: Profile-based loading
	fmt.Println("\n📖 Example 3: Profile-Based Configuration")

	// Create profile configuration
	profileConfigPath := "./example-config-prod.yaml"
	profileContent := `
debug: false

server:
  host: 0.0.0.0
  port: 443

database:
  host: prod-db.example.com
  ssl_mode: require
`

	err = os.WriteFile(profileConfigPath, []byte(profileContent), 0644)
	if err != nil {
		log.Fatal("Failed to create profile config:", err)
	}
	defer os.Remove(profileConfigPath) // Clean up

	prodCfg, err := konfig.LoadWithProfile(configPath, "prod")
	if err != nil {
		log.Fatal("Failed to load prod config:", err)
	}

	fmt.Printf("   Production Server: %s:%s\n",
		prodCfg.GetString("server.host"),
		prodCfg.GetString("server.port"))
	fmt.Printf("   Production Database: %s\n",
		prodCfg.GetString("database.host"))

	// Example 4: Error handling
	fmt.Println("\n📖 Example 4: Error Handling")
	_, err = konfig.Load("./nonexistent.yaml")
	if err != nil {
		fmt.Printf("   Expected error: %v\n", err)
	}

	// Example 5: Performance demonstration
	fmt.Println("\n📖 Example 5: Performance")
	start := time.Now()
	for i := 0; i < 1000; i++ {
		_ = cfg.GetString("name")
	}
	duration := time.Since(start)
	fmt.Printf("   1000 config accesses: %v (%.1f ns/op)\n",
		duration, float64(duration.Nanoseconds())/1000)

	fmt.Println("\n✅ All examples completed successfully!")
	fmt.Println("\nKey benefits of konfig:")
	fmt.Println("• Explicit file paths - no magic discovery")
	fmt.Println("• Type-safe struct loading with defaults")
	fmt.Println("• Environment variable substitution")
	fmt.Println("• Profile-based configuration")
	fmt.Println("• Lightning-fast performance")
	fmt.Println("• Production-ready error handling")
}
