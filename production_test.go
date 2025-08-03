package konfig

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProductionRobustness tests edge cases that occur in production
func TestProductionRobustness_MalformedYAML(t *testing.T) {
	testCases := []struct {
		name          string
		content       string
		errorContains string
	}{
		{
			name: "invalid_yaml_syntax",
			content: `
server:
  port: 8080
    host: localhost  # Wrong indentation
`,
			errorContains: "parse_error",
		},
		{
			name: "unclosed_brackets",
			content: `
server: {
  port: 8080
  host: localhost
# Missing closing bracket
`,
			errorContains: "parse_error",
		},
		{
			name:          "invalid_unicode",
			content:       "server:\n  port: 808\x00\x01\x02", // Null bytes
			errorContains: "parse_error",
		},
		{
			name:          "extremely_long_line",
			content:       fmt.Sprintf("server:\n  very_long_key: %s", strings.Repeat("x", 100000)),
			errorContains: "", // Should handle gracefully
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, "config.yaml")

			err := os.WriteFile(configPath, []byte(tc.content), 0644)
			require.NoError(t, err)

			cfg, err := Load(configPath)
			if tc.errorContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorContains)
				assert.Nil(t, cfg)
			} else {
				// Should handle gracefully even with extreme content
				if err != nil {
					t.Logf("Warning: %v", err)
				}
			}
		})
	}
}

func TestProductionRobustness_FileSystemEdgeCases(t *testing.T) {
	t.Run("empty_file", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "empty.yaml")

		err := os.WriteFile(configPath, []byte(""), 0644)
		require.NoError(t, err)

		cfg, err := Load(configPath)
		require.NoError(t, err)
		assert.Empty(t, cfg.Keys())
	})

	t.Run("readonly_file", func(t *testing.T) {
		if os.Getuid() == 0 {
			t.Skip("Skipping readonly test as root user")
		}

		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "readonly.yaml")

		content := "server:\n  port: 8080"
		err := os.WriteFile(configPath, []byte(content), 0444) // Read-only
		require.NoError(t, err)

		// Should still be able to read
		cfg, err := Load(configPath)
		require.NoError(t, err)
		assert.Equal(t, "8080", cfg.GetString("server.port"))
	})

	t.Run("deeply_nested_paths", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create deeply nested directory structure
		deepPath := tempDir
		for i := 0; i < 10; i++ {
			deepPath = filepath.Join(deepPath, fmt.Sprintf("level%d", i))
		}
		err := os.MkdirAll(deepPath, 0755)
		require.NoError(t, err)

		configPath := filepath.Join(deepPath, "config.yaml")
		content := "deep:\n  nested:\n    value: success"
		err = os.WriteFile(configPath, []byte(content), 0644)
		require.NoError(t, err)

		cfg, err := Load(configPath)
		require.NoError(t, err)
		assert.Equal(t, "success", cfg.GetString("deep.nested.value"))
	})
}

func TestProductionRobustness_LargeConfigurations(t *testing.T) {
	t.Run("large_config_file", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "large.yaml")

		// Generate large configuration with 1000+ keys
		var content strings.Builder
		content.WriteString("services:\n")
		for i := 0; i < 1000; i++ {
			content.WriteString(fmt.Sprintf("  service_%d:\n", i))
			content.WriteString(fmt.Sprintf("    name: service-%d\n", i))
			content.WriteString(fmt.Sprintf("    port: %d\n", 8000+i))
			content.WriteString(fmt.Sprintf("    enabled: %t\n", i%2 == 0))
		}

		err := os.WriteFile(configPath, []byte(content.String()), 0644)
		require.NoError(t, err)

		start := time.Now()
		cfg, err := Load(configPath)
		loadTime := time.Since(start)

		require.NoError(t, err)
		assert.NotNil(t, cfg)

		// Should handle large configs reasonably fast (< 100ms for 1000 keys)
		assert.Less(t, loadTime, 100*time.Millisecond, "Large config loading too slow")

		// Verify some keys loaded correctly
		assert.Equal(t, "service-0", cfg.GetString("services.service_0.name"))
		assert.Equal(t, "8999", cfg.GetString("services.service_999.port"))

		// Should have all keys (1000 services * 3 keys each = 3000)
		keys := cfg.Keys()
		assert.GreaterOrEqual(t, len(keys), 3000, "Not all keys loaded")

		t.Logf("Loaded %d keys in %v", len(keys), loadTime)
	})

	t.Run("deeply_nested_structure", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "deep.yaml")

		// Create 20-level deep nesting
		var content strings.Builder
		indent := ""
		for i := 0; i < 20; i++ {
			content.WriteString(fmt.Sprintf("%slevel_%d:\n", indent, i))
			indent += "  "
		}
		content.WriteString(fmt.Sprintf("%svalue: deep_success\n", indent))

		err := os.WriteFile(configPath, []byte(content.String()), 0644)
		require.NoError(t, err)

		cfg, err := Load(configPath)
		require.NoError(t, err)

		// Build expected key path
		var keyPath strings.Builder
		for i := 0; i < 20; i++ {
			if i > 0 {
				keyPath.WriteString(".")
			}
			keyPath.WriteString(fmt.Sprintf("level_%d", i))
		}
		keyPath.WriteString(".value")

		assert.Equal(t, "deep_success", cfg.GetString(keyPath.String()))
	})
}

func TestProductionRobustness_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "concurrent.yaml")

	content := `
server:
  port: 8080
  host: localhost
database:
  url: postgres://localhost/test
  pool_size: 10
cache:
  ttl: 300
  max_size: 1000
`
	err := os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	t.Run("concurrent_loading", func(t *testing.T) {
		const numGoroutines = 50
		const numIterations = 10

		var wg sync.WaitGroup
		errors := make(chan error, numGoroutines*numIterations)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numIterations; j++ {
					cfg, err := Load(configPath)
					if err != nil {
						errors <- err
						return
					}

					// Verify data integrity
					if cfg.GetString("server.port") != "8080" {
						errors <- fmt.Errorf("data corruption: expected 8080, got %s", cfg.GetString("server.port"))
						return
					}
				}
			}()
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		var errs []error
		for err := range errors {
			errs = append(errs, err)
		}

		if len(errs) > 0 {
			t.Fatalf("Concurrent loading failed with %d errors: %v", len(errs), errs[0])
		}

		t.Logf("Successfully completed %d concurrent loads", numGoroutines*numIterations)
	})

	t.Run("concurrent_struct_mapping", func(t *testing.T) {
		type Config struct {
			Server struct {
				Port int    `konfig:"port"`
				Host string `konfig:"host"`
			} `konfig:"server"`
			Database struct {
				URL      string `konfig:"url"`
				PoolSize int    `konfig:"pool_size"`
			} `konfig:"database"`
		}

		const numGoroutines = 25
		var wg sync.WaitGroup
		errors := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				var cfg Config
				err := LoadInto(configPath, &cfg)
				if err != nil {
					errors <- err
					return
				}

				// Verify struct was populated correctly
				if cfg.Server.Port != 8080 {
					errors <- fmt.Errorf("struct corruption: expected 8080, got %d", cfg.Server.Port)
					return
				}
				if cfg.Database.PoolSize != 10 {
					errors <- fmt.Errorf("struct corruption: expected 10, got %d", cfg.Database.PoolSize)
					return
				}
			}()
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		for err := range errors {
			require.NoError(t, err)
		}
	})
}

func TestProductionRobustness_MemoryManagement(t *testing.T) {
	t.Run("no_memory_leaks", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "memory_test.yaml")

		content := `
data:
  large_string: ` + strings.Repeat("x", 10000) + `
  numbers:
    - 1
    - 2
    - 3
`
		err := os.WriteFile(configPath, []byte(content), 0644)
		require.NoError(t, err)

		// Load and discard many configs to test for memory leaks
		for i := 0; i < 100; i++ {
			cfg, err := Load(configPath)
			require.NoError(t, err)

			// Use the config briefly
			_ = cfg.GetString("data.large_string")
			_ = cfg.Keys()

			// Let Go runtime know we're done with this config
			cfg = nil
		}

		// Force garbage collection
		for i := 0; i < 3; i++ {
			runtime.GC()
			time.Sleep(10 * time.Millisecond)
		}

		// If we get here without OOM, the test passes
		t.Log("Memory leak test completed successfully")
	})
}

func TestProductionRobustness_RealWorldPatterns(t *testing.T) {
	t.Run("microservice_config", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "microservice.yaml")

		content := `
service:
  name: user-service
  version: 1.2.3
  port: 8080

database:
  host: ${DB_HOST:localhost}
  port: ${DB_PORT:5432}
  name: users
  ssl_mode: ${DB_SSL:require}

redis:
  url: ${REDIS_URL:redis://localhost:6379}
  pool_size: ${REDIS_POOL:10}

logging:
  level: ${LOG_LEVEL:info}
  format: json

metrics:
  enabled: true
  port: 9090
  path: /metrics

health:
  enabled: true
  port: 8081
  path: /health
`

		// Set some environment variables
		os.Setenv("DB_HOST", "prod-db.example.com")
		os.Setenv("LOG_LEVEL", "warn")
		defer func() {
			os.Unsetenv("DB_HOST")
			os.Unsetenv("LOG_LEVEL")
		}()

		err := os.WriteFile(configPath, []byte(content), 0644)
		require.NoError(t, err)

		cfg, err := Load(configPath)
		require.NoError(t, err)

		// Verify environment variable substitution
		assert.Equal(t, "prod-db.example.com", cfg.GetString("database.host"))
		assert.Equal(t, "5432", cfg.GetString("database.port")) // Default used
		assert.Equal(t, "warn", cfg.GetString("logging.level"))
		assert.Equal(t, "redis://localhost:6379", cfg.GetString("redis.url")) // Default used

		// Verify regular values
		assert.Equal(t, "user-service", cfg.GetString("service.name"))
		assert.Equal(t, true, cfg.GetBool("metrics.enabled"))
	})

	t.Run("kubernetes_deployment_config", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "k8s.yaml")

		content := `
deployment:
  replicas: ${REPLICAS:3}
  namespace: ${NAMESPACE:default}
  
resources:
  requests:
    memory: "${MEMORY_REQUEST:256Mi}"
    cpu: "${CPU_REQUEST:100m}"
  limits:
    memory: "${MEMORY_LIMIT:512Mi}"
    cpu: "${CPU_LIMIT:500m}"

ingress:
  enabled: ${INGRESS_ENABLED:true}
  host: ${APP_HOST:app.example.com}
  tls: ${TLS_ENABLED:true}
`

		// Simulate K8s environment variables
		os.Setenv("REPLICAS", "5")
		os.Setenv("MEMORY_LIMIT", "1Gi")
		os.Setenv("APP_HOST", "prod.example.com")
		defer func() {
			os.Unsetenv("REPLICAS")
			os.Unsetenv("MEMORY_LIMIT")
			os.Unsetenv("APP_HOST")
		}()

		err := os.WriteFile(configPath, []byte(content), 0644)
		require.NoError(t, err)

		cfg, err := Load(configPath)
		require.NoError(t, err)

		// Verify K8s-style configuration
		assert.Equal(t, 5, cfg.GetInt("deployment.replicas"))
		assert.Equal(t, "256Mi", cfg.GetString("resources.requests.memory"))
		assert.Equal(t, "1Gi", cfg.GetString("resources.limits.memory"))
		assert.Equal(t, "prod.example.com", cfg.GetString("ingress.host"))
		assert.Equal(t, true, cfg.GetBool("ingress.tls"))
	})
}
