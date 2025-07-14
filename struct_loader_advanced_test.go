package konfig

import (
	"os"
	"testing"
)

// Advanced tests for the LoadInto functionality: nested structs and error handling

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

// Error handling and edge case tests

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

func TestLoadInto_FieldWithoutKonfigTag(t *testing.T) {
	type Config struct {
		WithTag    string `konfig:"tagged.field" default:"tagged_value"`
		WithoutTag string // No konfig tag, should be ignored
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Fatalf("LoadInto failed: %v", err)
	}

	if cfg.WithTag != "tagged_value" {
		t.Errorf("Expected WithTag 'tagged_value', got '%s'", cfg.WithTag)
	}

	if cfg.WithoutTag != "" {
		t.Errorf("Expected WithoutTag to be empty, got '%s'", cfg.WithoutTag)
	}
}

func TestLoadInto_EmptyStructNoError(t *testing.T) {
	type Config struct{}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Errorf("Expected no error for empty struct, got: %v", err)
	}
}

func TestLoadInto_StructWithOnlyUntaggedFields(t *testing.T) {
	type Config struct {
		Field1 string
		Field2 int
		Field3 bool
	}

	var cfg Config
	err := LoadInto(&cfg)

	if err != nil {
		t.Errorf("Expected no error for struct with only untagged fields, got: %v", err)
	}

	// All fields should remain at their zero values
	if cfg.Field1 != "" {
		t.Errorf("Expected Field1 to be empty, got '%s'", cfg.Field1)
	}

	if cfg.Field2 != 0 {
		t.Errorf("Expected Field2 to be 0, got %d", cfg.Field2)
	}

	if cfg.Field3 != false {
		t.Errorf("Expected Field3 to be false, got %v", cfg.Field3)
	}
}
