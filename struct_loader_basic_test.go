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
