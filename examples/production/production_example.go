package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mfenderov/konfig"
)

// AppConfig represents a real-world microservice configuration
type AppConfig struct {
	Service  ServiceConfig  `konfig:"service"`
	Server   ServerConfig   `konfig:"server"`
	Database DatabaseConfig `konfig:"database"`
	Cache    CacheConfig    `konfig:"cache"`
	Logging  LoggingConfig  `konfig:"logging"`
	Metrics  MetricsConfig  `konfig:"metrics"`
}

type ServiceConfig struct {
	Name        string `konfig:"name" default:"my-service"`
	Version     string `konfig:"version" default:"1.0.0"`
	Environment string `konfig:"environment" default:"development"`
}

type ServerConfig struct {
	Host    string        `konfig:"host" default:"localhost"`
	Port    int           `konfig:"port" default:"8080"`
	Timeout time.Duration `konfig:"timeout" default:"30s"`
}

type DatabaseConfig struct {
	Host     string `konfig:"host" default:"localhost"`
	Port     int    `konfig:"port" default:"5432"`
	Name     string `konfig:"name" default:"myapp"`
	Username string `konfig:"username" default:"user"`
	Password string `konfig:"password" default:"password"`
	SSLMode  string `konfig:"ssl_mode" default:"disable"`
	PoolSize int    `konfig:"pool_size" default:"10"`
}

type CacheConfig struct {
	Enabled bool   `konfig:"enabled" default:"false"`
	URL     string `konfig:"url" default:"redis://localhost:6379"`
	TTL     int    `konfig:"ttl" default:"300"`
}

type LoggingConfig struct {
	Level  string `konfig:"level" default:"info"`
	Format string `konfig:"format" default:"json"`
}

type MetricsConfig struct {
	Enabled bool   `konfig:"enabled" default:"true"`
	Port    int    `konfig:"port" default:"9090"`
	Path    string `konfig:"path" default:"/metrics"`
}

func main() {
	fmt.Println("🚀 konfig Production Example")
	fmt.Println("===========================")

	// Example 1: Simple configuration loading
	fmt.Println("\n1. Loading base configuration...")
	cfg, err := konfig.Load("./resources/application.yaml")
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		return
	}

	fmt.Printf("   Service Name: %s\n", cfg.GetString("service.name"))
	fmt.Printf("   Server Port: %s\n", cfg.GetString("server.port"))

	// Example 2: Profile-based configuration
	fmt.Println("\n2. Loading with development profile...")
	devCfg, err := konfig.LoadWithProfile("./resources/application.yaml", "dev")
	if err != nil {
		log.Printf("Failed to load dev config: %v", err)
		return
	}

	fmt.Printf("   Dev Server Host: %s\n", devCfg.GetString("server.host"))
	fmt.Printf("   Dev Database: %s\n", devCfg.GetString("database.host"))

	// Example 3: Type-safe struct loading
	fmt.Println("\n3. Loading into type-safe struct...")
	var appConfig AppConfig
	err = konfig.LoadInto("./resources/application.yaml", &appConfig)
	if err != nil {
		log.Printf("Failed to load into struct: %v", err)
		return
	}

	fmt.Printf("   Service: %s v%s (%s)\n",
		appConfig.Service.Name,
		appConfig.Service.Version,
		appConfig.Service.Environment)
	fmt.Printf("   Server: %s:%d (timeout: %v)\n",
		appConfig.Server.Host,
		appConfig.Server.Port,
		appConfig.Server.Timeout)
	fmt.Printf("   Database: %s@%s:%d (pool: %d)\n",
		appConfig.Database.Name,
		appConfig.Database.Host,
		appConfig.Database.Port,
		appConfig.Database.PoolSize)

	// Example 4: Profile-based struct loading
	fmt.Println("\n4. Loading production profile into struct...")
	var prodConfig AppConfig
	err = konfig.LoadIntoWithProfile("./resources/application.yaml", "prod", &prodConfig)
	if err != nil {
		log.Printf("Failed to load prod config: %v", err)
		return
	}

	fmt.Printf("   Prod Service: %s (%s)\n",
		prodConfig.Service.Name,
		prodConfig.Service.Environment)
	fmt.Printf("   Prod Database: %s (SSL: %s)\n",
		prodConfig.Database.Host,
		prodConfig.Database.SSLMode)
	fmt.Printf("   Metrics: Enabled=%t, Port=%d\n",
		prodConfig.Metrics.Enabled,
		prodConfig.Metrics.Port)

	// Example 5: Error handling
	fmt.Println("\n5. Demonstrating error handling...")
	_, err = konfig.Load("./nonexistent.yaml")
	if err != nil {
		fmt.Printf("   Expected error: %v\n", err)
	}

	// Example 6: Performance demonstration
	fmt.Println("\n6. Performance test...")
	start := time.Now()
	for i := 0; i < 1000; i++ {
		_ = cfg.GetString("server.port")
	}
	duration := time.Since(start)
	fmt.Printf("   1000 config accesses: %v (%.2f ns/op)\n", duration, float64(duration.Nanoseconds())/1000)

	fmt.Println("\n✅ All examples completed successfully!")
}
