package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mfenderov/konfig"
)

// Simple example demonstrating basic struct-based configuration with konfig

// DatabaseConfig represents database settings
type DatabaseConfig struct {
	Host     string `konfig:"host" default:"localhost"`
	Port     string `konfig:"port" default:"5432"`
	Name     string `konfig:"name" default:"myapp"`
	User     string `konfig:"user" default:"postgres"`
	Password string `konfig:"password" default:"secret"`
}

// ServerConfig represents server settings
type ServerConfig struct {
	Host string `konfig:"host" default:"0.0.0.0"`
	Port string `konfig:"port" default:"8080"`
}

// SimpleConfig represents the complete application configuration
type SimpleConfig struct {
	AppName  string         `konfig:"application.name" default:"my-app"`
	Version  string         `konfig:"application.version" default:"1.0.0"`
	Debug    string         `konfig:"application.debug" default:"true"`
	Database DatabaseConfig `konfig:"database"`
	Server   ServerConfig   `konfig:"server"`
}

func main() {
	fmt.Println("🚀 Simple Konfig Example")
	fmt.Println("========================")

	// Optional: Set some environment variables to demonstrate overrides
	os.Setenv("application.name", "simple-example")
	os.Setenv("server.port", "9090")
	os.Setenv("database.host", "db.example.com")
	defer func() {
		os.Unsetenv("application.name")
		os.Unsetenv("server.port")
		os.Unsetenv("database.host")
	}()

	// Load configuration into struct
	var config SimpleConfig
	err := konfig.LoadInto(&config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Display loaded configuration
	fmt.Printf("📱 Application: %s v%s (Debug: %s)\n", config.AppName, config.Version, config.Debug)
	fmt.Printf("🌐 Server: %s:%s\n", config.Server.Host, config.Server.Port)
	fmt.Printf("🗄️  Database: %s@%s:%s/%s\n", 
		config.Database.User, config.Database.Host, 
		config.Database.Port, config.Database.Name)

	// Check current profile
	profile := konfig.GetProfile()
	if profile != "" {
		fmt.Printf("🏷️  Profile: %s\n", profile)
	} else {
		fmt.Printf("🏷️  Profile: default\n")
	}

	fmt.Println("\n✅ Configuration loaded successfully!")
	fmt.Println("Try running with different profiles:")
	fmt.Println("  go run simple_example.go -p dev")
	fmt.Println("  go run simple_example.go -p prod")
}