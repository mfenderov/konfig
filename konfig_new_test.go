package konfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAPI_Load(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "app.yaml")

	configContent := `
server:
  port: 8080
  host: localhost
database:
  name: myapp
  url: postgres://localhost/myapp
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Test Load function
	cfg, err := Load(configPath)
	require.NoError(t, err)

	// Test type-safe getters
	assert.Equal(t, "8080", cfg.GetString("server.port"))
	assert.Equal(t, "localhost", cfg.GetString("server.host"))
	assert.Equal(t, "myapp", cfg.GetString("database.name"))

	// Test integer conversion
	assert.Equal(t, 8080, cfg.GetInt("server.port"))

	// Test non-existent key
	assert.Equal(t, "", cfg.GetString("nonexistent.key"))
	assert.Equal(t, 0, cfg.GetInt("nonexistent.key"))
}

func TestNewAPI_LoadWithProfile(t *testing.T) {
	tempDir := t.TempDir()

	// Base config
	baseConfigPath := filepath.Join(tempDir, "app.yaml")
	baseConfig := `
server:
  port: 8080
  host: localhost
env: base
`
	err := os.WriteFile(baseConfigPath, []byte(baseConfig), 0644)
	require.NoError(t, err)

	// Profile config
	profileConfigPath := filepath.Join(tempDir, "app-dev.yaml")
	profileConfig := `
server:
  port: 3000
env: development
debug: true
`
	err = os.WriteFile(profileConfigPath, []byte(profileConfig), 0644)
	require.NoError(t, err)

	// Test LoadWithProfile
	cfg, err := LoadWithProfile(baseConfigPath, "dev")
	require.NoError(t, err)

	// Base config values should be loaded
	assert.Equal(t, "localhost", cfg.GetString("server.host"))

	// Profile values should override base values
	assert.Equal(t, "3000", cfg.GetString("server.port"))
	assert.Equal(t, "development", cfg.GetString("env"))
	assert.Equal(t, "true", cfg.GetString("debug"))
}

func TestNewAPI_LoadInto(t *testing.T) {
	// Create temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "app.yaml")

	configContent := `
server:
  port: 8080
  host: localhost
  debug: true
timeout: 30s
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Define config struct
	type ServerConfig struct {
		Port  int    `konfig:"port" default:"3000"`
		Host  string `konfig:"host" default:"0.0.0.0"`
		Debug bool   `konfig:"debug" default:"false"`
	}

	type Config struct {
		Server  ServerConfig `konfig:"server"`
		Timeout string       `konfig:"timeout" default:"10s"`
	}

	// Test LoadInto
	var cfg Config
	err = LoadInto(configPath, &cfg)
	require.NoError(t, err)

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, true, cfg.Server.Debug)
	assert.Equal(t, "30s", cfg.Timeout)
}

func TestNewAPI_ErrorHandling(t *testing.T) {
	// Test file not found
	_, err := Load("nonexistent.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "file_not_found")

	// Test empty path
	_, err = Load("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation_error")

	// Test invalid struct target
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "app.yaml")
	err = os.WriteFile(configPath, []byte("key: value"), 0644)
	require.NoError(t, err)

	var notAPointer string
	err = LoadInto(configPath, notAPointer)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "validation_error")
}

func TestNewAPI_EnvSubstitution(t *testing.T) {
	// Set test environment variable
	os.Setenv("TEST_PORT", "9000")
	os.Setenv("TEST_HOST", "0.0.0.0")
	defer func() {
		os.Unsetenv("TEST_PORT")
		os.Unsetenv("TEST_HOST")
	}()

	// Create config with env substitutions
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "app.yaml")

	configContent := `
server:
  port: ${TEST_PORT:8080}
  host: ${TEST_HOST:localhost}
  protocol: ${UNDEFINED_VAR:http}
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	cfg, err := Load(configPath)
	require.NoError(t, err)

	// Environment variables should be substituted
	assert.Equal(t, "9000", cfg.GetString("server.port"))
	assert.Equal(t, "0.0.0.0", cfg.GetString("server.host"))

	// Default should be used for undefined variables
	assert.Equal(t, "http", cfg.GetString("server.protocol"))
}
