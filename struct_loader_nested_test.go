package konfig

import (
	"os"
	"testing"
)

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