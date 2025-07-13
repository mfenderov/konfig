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