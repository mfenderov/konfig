package konfig

import (
	"os"
	"testing"
)

// Error handling and edge case tests for the LoadInto functionality

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

	// Fields with empty konfig tag should remain empty (same as no konfig tag)
	if cfg.EmptyTag != "" {
		t.Errorf("Expected EmptyTag to be empty (empty konfig tag), got '%s'", cfg.EmptyTag)
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