package konfig

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFrom_SimpleConfiguration(t *testing.T) {
	ClearEnv()

	_ = SetEnv("POSTGRES_PASSWORD", "12345")
	_ = SetEnv("POSTGRES_HOST", "localhost")

	err := LoadFrom("test-data/application.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "12345", GetEnv("db.password"))
	assert.Equal(t, "https://localhost", GetEnv("db.host"))
	assert.Equal(t, "5432", GetEnv("db.port"))
	assert.Equal(t, "user", GetEnv("db.user"))
	assert.Equal(t, "8080", GetEnv("server.port"))
	assert.Equal(t, "debug", GetEnv("log.level"))
	assert.Equal(t, "[one two three]", GetEnv("list"))
}

func TestLoadFrom_ConfigurationWithEmptyProperties(t *testing.T) {
	ClearEnv()

	err := LoadFrom("test-data/application-empty.yaml")
	assert.EqualError(t, err, "property 'db.password' is nil")
}

func TestLoadFrom_ConfigurationWithEmptyPropertyWithDefault(t *testing.T) {
	ClearEnv()

	err := LoadFrom("test-data/application-empty-with-default.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "postgres", GetEnv("db.password"))
	assert.Equal(t, "https://", GetEnv("db.host"))
}

func TestLoad_WithoutProfile(t *testing.T) {
	ClearEnv()

	err := Load()
	assert.NoError(t, err)

	port := GetEnv("server.port")
	assert.Equal(t, "1234", port)
}

func TestLoad_WithTestProfile(t *testing.T) {
	ClearEnv()

	// Set up command-line flags for the test profile
	os.Args = []string{os.Args[0], "-p", "test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	ResetProfileInitialized()

	err := Load()
	assert.NoError(t, err)

	port := GetEnv("server.port")
	assert.Equal(t, "4321", port)
}

func TestLoadFrom_NonExistentFile(t *testing.T) {
	ClearEnv()
	err := LoadFrom("test-data/non-existent.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestLoadFrom_InvalidYAML(t *testing.T) {
	ClearEnv()
	err := LoadFrom("test-data/invalid.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal config file")
}

func TestLoadFrom_NestedConfiguration(t *testing.T) {
	ClearEnv()
	err := LoadFrom("test-data/nested.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "value1", GetEnv("level1.level2.key1"))
	assert.Equal(t, "value2", GetEnv("level1.level2.key2"))
	assert.Equal(t, "[1 2 3]", GetEnv("level1.array"))
}

func TestLoadFrom_EmptyFile(t *testing.T) {
	ClearEnv()
	err := LoadFrom("test-data/empty.yaml")
	assert.NoError(t, err, "Empty config file should be allowed")
}

func TestEnvVariableOverride(t *testing.T) {
	ClearEnv()

	// Set environment variable that should override config
	_ = SetEnv("CUSTOM_VALUE", "override_value")

	err := LoadFrom("test-data/env-override.yaml")
	assert.NoError(t, err)

	// Should use environment variable value instead of default
	assert.Equal(t, "override_value", GetEnv("custom.key"))
	// Should use default value when env var is not set
	assert.Equal(t, "default_value", GetEnv("custom.key2"))
}

func TestLoadFrom_InvalidPath(t *testing.T) {
	ClearEnv()
	err := LoadFrom("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config file path cannot be empty")
}

func TestLoadFrom_InvalidExtension(t *testing.T) {
	ClearEnv()
	err := LoadFrom("test-data/config.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config file must have .yml or .yaml extension")
}
