package konfig

import (
	"os"
	"testing"
)

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

func TestLoadInto_InvalidInput(t *testing.T) {
	// Test with nil pointer
	err := LoadInto(nil)
	if err == nil {
		t.Error("Expected error for nil input")
	}

	// Test with non-pointer
	var cfg struct{}
	err = LoadInto(cfg)
	if err == nil {
		t.Error("Expected error for non-pointer input")
	}

	// Test with pointer to non-struct
	var str string
	err = LoadInto(&str)
	if err == nil {
		t.Error("Expected error for pointer to non-struct")
	}
}

func TestLoadInto_NestedStruct(t *testing.T) {
	// Set environment values for nested structure
	os.Setenv("nested.server.host", "nestedhost")
	os.Setenv("nested.server.port", "9000")
	os.Setenv("nested.database.name", "testdb")
	defer func() {
		os.Unsetenv("nested.server.host")
		os.Unsetenv("nested.server.port")
		os.Unsetenv("nested.database.name")
	}()

	type ServerConfig struct {
		Host string `konfig:"host" default:"localhost"`
		Port string `konfig:"port" default:"8080"`
	}

	type DatabaseConfig struct {
		Name string `konfig:"name" default:"defaultdb"`
	}

	type Config struct {
		Server   ServerConfig   `konfig:"nested.server"`
		Database DatabaseConfig `konfig:"nested.database"`
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.Server.Host != "nestedhost" {
		t.Errorf("Expected Server.Host 'nestedhost', got '%s'", cfg.Server.Host)
	}

	if cfg.Server.Port != "9000" {
		t.Errorf("Expected Server.Port '9000', got '%s'", cfg.Server.Port)
	}

	if cfg.Database.Name != "testdb" {
		t.Errorf("Expected Database.Name 'testdb', got '%s'", cfg.Database.Name)
	}
}

func TestLoadInto_NestedStructDefaults(t *testing.T) {
	type ServerConfig struct {
		Host string `konfig:"host" default:"defaulthost"`
		Port string `konfig:"port" default:"3333"`
	}

	type Config struct {
		Server ServerConfig `konfig:"nesteddefault.server"`
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.Server.Host != "defaulthost" {
		t.Errorf("Expected nested default Server.Host 'defaulthost', got '%s'", cfg.Server.Host)
	}

	if cfg.Server.Port != "3333" {
		t.Errorf("Expected nested default Server.Port '3333', got '%s'", cfg.Server.Port)
	}
}

// Additional comprehensive tests for the LoadInto functionality

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

func TestLoadInto_EmptyKonfigTag(t *testing.T) {
	// Test struct with empty konfig tag - should be ignored
	type Config struct {
		WithTag    string `konfig:"emptytest.with_tag" default:"with_default"`
		WithoutTag string `default:"without_konfig_tag"`
		EmptyTag   string `konfig:"" default:"empty_tag_default"`
	}

	os.Setenv("emptytest.with_tag", "tag_value")
	defer os.Unsetenv("emptytest.with_tag")

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.WithTag != "tag_value" {
		t.Errorf("Expected WithTag 'tag_value', got '%s'", cfg.WithTag)
	}

	// Fields without konfig tag should remain empty (not populated by LoadInto)
	if cfg.WithoutTag != "" {
		t.Errorf("Expected WithoutTag to be empty (no konfig tag), got '%s'", cfg.WithoutTag)
	}

	// Fields with empty konfig tag should use default value
	if cfg.EmptyTag != "empty_tag_default" {
		t.Errorf("Expected EmptyTag 'empty_tag_default', got '%s'", cfg.EmptyTag)
	}
}

func TestLoadInto_DeepNestedStructs(t *testing.T) {
	// Test deeply nested structures
	os.Setenv("deep.level1.level2.level3.value", "deep_value")
	os.Setenv("deep.level1.level2.simple", "simple_value")
	defer func() {
		os.Unsetenv("deep.level1.level2.level3.value")
		os.Unsetenv("deep.level1.level2.simple")
	}()

	type Level3Config struct {
		Value string `konfig:"value" default:"default_deep"`
	}

	type Level2Config struct {
		Level3 Level3Config `konfig:"level3"`
		Simple string       `konfig:"simple" default:"default_simple"`
	}

	type Level1Config struct {
		Level2 Level2Config `konfig:"level2"`
	}

	type Config struct {
		Deep Level1Config `konfig:"deep.level1"`
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.Deep.Level2.Level3.Value != "deep_value" {
		t.Errorf("Expected deep nested value 'deep_value', got '%s'", cfg.Deep.Level2.Level3.Value)
	}

	if cfg.Deep.Level2.Simple != "simple_value" {
		t.Errorf("Expected simple value 'simple_value', got '%s'", cfg.Deep.Level2.Simple)
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

func TestLoadInto_StructWithPointers(t *testing.T) {
	// Test struct with pointer fields (should be skipped gracefully)
	type Config struct {
		StringVal  string  `konfig:"pointer.string" default:"string_value"`
		PointerVal *string `konfig:"pointer.ptr" default:"ptr_value"`
	}

	os.Setenv("pointer.string", "env_string")
	defer os.Unsetenv("pointer.string")

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.StringVal != "env_string" {
		t.Errorf("Expected StringVal 'env_string', got '%s'", cfg.StringVal)
	}

	// Pointer field should remain nil since we can't handle pointers yet
	if cfg.PointerVal != nil {
		t.Errorf("Expected PointerVal to be nil, got %v", cfg.PointerVal)
	}
}

func TestLoadInto_SpecialCharactersInValues(t *testing.T) {
	// Test values with special characters
	specialValue := "value with spaces, symbols: !@#$%^&*()_+-={}[]|\\:;\"'<>?,./"
	os.Setenv("special.value", specialValue)
	defer os.Unsetenv("special.value")

	type Config struct {
		SpecialValue string `konfig:"special.value" default:"default"`
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.SpecialValue != specialValue {
		t.Errorf("Expected SpecialValue '%s', got '%s'", specialValue, cfg.SpecialValue)
	}
}

func TestLoadInto_EmptyEnvironmentValues(t *testing.T) {
	// Test empty environment values vs default values
	os.Setenv("empty.explicit", "")
	os.Setenv("empty.whitespace", "   ")
	defer func() {
		os.Unsetenv("empty.explicit")
		os.Unsetenv("empty.whitespace")
	}()

	type Config struct {
		ExplicitEmpty string `konfig:"empty.explicit" default:"should_use_default"`
		Whitespace    string `konfig:"empty.whitespace" default:"should_use_default"`
		NotSet        string `konfig:"empty.not_set" default:"default_value"`
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	// Empty env var should override default (empty string is a valid value)
	if cfg.ExplicitEmpty != "" {
		t.Errorf("Expected ExplicitEmpty to be empty, got '%s'", cfg.ExplicitEmpty)
	}

	// Whitespace-only env var should be preserved
	if cfg.Whitespace != "   " {
		t.Errorf("Expected Whitespace to be '   ', got '%s'", cfg.Whitespace)
	}

	// Unset env var should use default
	if cfg.NotSet != "default_value" {
		t.Errorf("Expected NotSet 'default_value', got '%s'", cfg.NotSet)
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
