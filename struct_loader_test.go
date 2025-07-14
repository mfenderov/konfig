package konfig

import (
	"os"
	"testing"
)

// Comprehensive tests for the LoadInto functionality

func TestLoadInto_BasicStruct(t *testing.T) {
	// Use unique test keys to avoid conflicts with existing YAML
	os.Setenv("test.port", "3000")
	os.Setenv("test.host", "testhost")
	defer func() {
		os.Unsetenv("test.port")
		os.Unsetenv("test.host")
	}()

	type Config struct {
		ServerPort string `konfig:"test.port" default:"8080"`
		ServerHost string `konfig:"test.host" default:"localhost"`
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.ServerPort != "3000" {
		t.Errorf("Expected ServerPort '3000', got '%s'", cfg.ServerPort)
	}

	if cfg.ServerHost != "testhost" {
		t.Errorf("Expected ServerHost 'testhost', got '%s'", cfg.ServerHost)
	}
}

func TestLoadInto_DefaultValues(t *testing.T) {
	// Use unique test keys with no environment values set
	type Config struct {
		ServerPort string `konfig:"defaulttest.port" default:"8080"`
		ServerHost string `konfig:"defaulttest.host" default:"localhost"`
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.ServerPort != "8080" {
		t.Errorf("Expected default ServerPort '8080', got '%s'", cfg.ServerPort)
	}

	if cfg.ServerHost != "localhost" {
		t.Errorf("Expected default ServerHost 'localhost', got '%s'", cfg.ServerHost)
	}
}

func TestLoadInto_DifferentDataTypes(t *testing.T) {
	// Set up test environment variables
	os.Setenv("datatypes.string_val", "test_string")
	os.Setenv("datatypes.int_val", "42")
	os.Setenv("datatypes.bool_val", "true")
	os.Setenv("datatypes.float_val", "3.14")
	defer func() {
		os.Unsetenv("datatypes.string_val")
		os.Unsetenv("datatypes.int_val")
		os.Unsetenv("datatypes.bool_val")
		os.Unsetenv("datatypes.float_val")
	}()

	type Config struct {
		StringVal string `konfig:"datatypes.string_val" default:"default_string"`
		IntVal    string `konfig:"datatypes.int_val" default:"0"`
		BoolVal   string `konfig:"datatypes.bool_val" default:"false"`
		FloatVal  string `konfig:"datatypes.float_val" default:"0.0"`
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.StringVal != "test_string" {
		t.Errorf("Expected StringVal 'test_string', got '%s'", cfg.StringVal)
	}

	if cfg.IntVal != "42" {
		t.Errorf("Expected IntVal '42', got '%s'", cfg.IntVal)
	}

	if cfg.BoolVal != "true" {
		t.Errorf("Expected BoolVal 'true', got '%s'", cfg.BoolVal)
	}

	if cfg.FloatVal != "3.14" {
		t.Errorf("Expected FloatVal '3.14', got '%s'", cfg.FloatVal)
	}
}

func TestLoadInto_MixedEnvAndDefaults(t *testing.T) {
	// Test mix of environment variables and default values
	os.Setenv("mixed.provided", "env_value")
	// Don't set mixed.default - should use default
	defer os.Unsetenv("mixed.provided")

	type Config struct {
		Provided   string `konfig:"mixed.provided" default:"default_provided"`
		UseDefault string `konfig:"mixed.default" default:"default_value"`
		NoDefault  string `konfig:"mixed.no_default"`
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.Provided != "env_value" {
		t.Errorf("Expected Provided 'env_value', got '%s'", cfg.Provided)
	}

	if cfg.UseDefault != "default_value" {
		t.Errorf("Expected UseDefault 'default_value', got '%s'", cfg.UseDefault)
	}

	if cfg.NoDefault != "" {
		t.Errorf("Expected NoDefault to be empty, got '%s'", cfg.NoDefault)
	}
}

func TestLoadInto_LargeConfiguration(t *testing.T) {
	// Test with a large, realistic configuration structure
	type DatabaseConfig struct {
		Host     string `konfig:"host" default:"localhost"`
		Port     string `konfig:"port" default:"5432"`
		Name     string `konfig:"name" default:"myapp"`
		User     string `konfig:"user" default:"postgres"`
		Password string `konfig:"password" default:"secret"`
		SSLMode  string `konfig:"ssl_mode" default:"disable"`
	}

	type ServerConfig struct {
		Host           string `konfig:"host" default:"0.0.0.0"`
		Port           string `konfig:"port" default:"8080"`
		ReadTimeout    string `konfig:"read_timeout" default:"30s"`
		WriteTimeout   string `konfig:"write_timeout" default:"30s"`
		MaxHeaderBytes string `konfig:"max_header_bytes" default:"1048576"`
	}

	type LoggingConfig struct {
		Level  string `konfig:"level" default:"info"`
		Format string `konfig:"format" default:"json"`
		Output string `konfig:"output" default:"stdout"`
	}

	type Config struct {
		App      string         `konfig:"large.app.name" default:"myapp"`
		Version  string         `konfig:"large.app.version" default:"1.0.0"`
		Debug    string         `konfig:"large.app.debug" default:"false"`
		Database DatabaseConfig `konfig:"large.database"`
		Server   ServerConfig   `konfig:"large.server"`
		Logging  LoggingConfig  `konfig:"large.logging"`
	}

	// Set some env vars, leave others to use defaults
	os.Setenv("large.app.name", "test-app")
	os.Setenv("large.database.host", "db.example.com")
	os.Setenv("large.database.port", "3306")
	os.Setenv("large.server.port", "9090")
	os.Setenv("large.logging.level", "debug")

	defer func() {
		os.Unsetenv("large.app.name")
		os.Unsetenv("large.database.host")
		os.Unsetenv("large.database.port")
		os.Unsetenv("large.server.port")
		os.Unsetenv("large.logging.level")
	}()

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	// Check overridden values
	if cfg.App != "test-app" {
		t.Errorf("Expected App 'test-app', got '%s'", cfg.App)
	}
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("Expected Database.Host 'db.example.com', got '%s'", cfg.Database.Host)
	}
	if cfg.Database.Port != "3306" {
		t.Errorf("Expected Database.Port '3306', got '%s'", cfg.Database.Port)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("Expected Server.Port '9090', got '%s'", cfg.Server.Port)
	}
	if cfg.Logging.Level != "debug" {
		t.Errorf("Expected Logging.Level 'debug', got '%s'", cfg.Logging.Level)
	}

	// Check default values
	if cfg.Version != "1.0.0" {
		t.Errorf("Expected Version '1.0.0', got '%s'", cfg.Version)
	}
	if cfg.Database.Name != "myapp" {
		t.Errorf("Expected Database.Name 'myapp', got '%s'", cfg.Database.Name)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected Server.Host '0.0.0.0', got '%s'", cfg.Server.Host)
	}
	if cfg.Logging.Format != "json" {
		t.Errorf("Expected Logging.Format 'json', got '%s'", cfg.Logging.Format)
	}
}
