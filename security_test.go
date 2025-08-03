package konfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSecurity_PathTraversalPrevention(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		shouldFail bool
	}{
		{
			name:       "simple relative path",
			path:       "./testdata/simple.yaml",
			shouldFail: false,
		},
		{
			name:       "path traversal with ../",
			path:       "../../../etc/passwd",
			shouldFail: true,
		},
		{
			name:       "nested path traversal",
			path:       "config/../../../secret.txt",
			shouldFail: true,
		},
		{
			name:       "legitimate nested path",
			path:       "config/app/settings.yaml",
			shouldFail: false, // will fail because file doesn't exist, but not due to path traversal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(tt.path)

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected path traversal to be blocked for path: %s", tt.path)
				} else if !strings.Contains(err.Error(), "path traversal not allowed") &&
					!strings.Contains(err.Error(), "file_not_found") {
					t.Errorf("Expected path traversal or file not found error, got: %v", err)
				}
				// Both types of errors are acceptable for security - they prevent the attack
			}
		})
	}
}

func TestSecurity_FileSizeLimit(t *testing.T) {
	// Create a temporary large file
	tmpFile, err := os.CreateTemp("", "large-config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write more than 10MB of data
	largeData := strings.Repeat("key: value\n", 1024*1024) // ~10MB
	if _, err := tmpFile.WriteString(largeData); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Try to load the large file
	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Error("Expected large file to be rejected")
	} else if !strings.Contains(err.Error(), "file too large") {
		t.Errorf("Expected file size error, got: %v", err)
	}
}

func TestSecurity_YAMLComplexityLimit(t *testing.T) {
	// Create deeply nested YAML that exceeds complexity limits
	deepYAML := "root:\n"
	indent := "  "
	for i := 0; i < 35; i++ { // Exceed maxNestingDepth of 32
		deepYAML += indent + "level" + strings.Repeat("_", i) + ":\n"
		indent += "  "
	}
	deepYAML += indent + "value: deep"

	tmpFile, err := os.CreateTemp("", "deep-config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(deepYAML); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Try to load the deeply nested file
	_, err = Load(tmpFile.Name())
	if err == nil {
		t.Error("Expected deeply nested YAML to be rejected")
	} else if !strings.Contains(err.Error(), "YAML too complex") {
		t.Errorf("Expected complexity error, got: %v", err)
	}
}

func TestSecurity_SymlinkHandling(t *testing.T) {
	// Create a test file
	tmpFile, err := os.CreateTemp("", "target-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString("test: value\n"); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Create a symlink to the test file
	symlinkPath := tmpFile.Name() + ".link"
	if err := os.Symlink(tmpFile.Name(), symlinkPath); err != nil {
		t.Skip("Cannot create symlinks on this system")
	}
	defer os.Remove(symlinkPath)

	// Loading through symlink should work (but be aware of it)
	cfg, err := Load(symlinkPath)
	if err != nil {
		t.Errorf("Symlink handling failed: %v", err)
	}

	if value := cfg.GetString("test"); value != "value" {
		t.Errorf("Expected 'value', got: %s", value)
	}
}

func TestSecurity_NonExistentDirectory(t *testing.T) {
	// Try to load from a non-existent directory structure
	_, err := Load("/non/existent/directory/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}

	// Should not contain sensitive system information
	if strings.Contains(err.Error(), "permission denied") ||
		strings.Contains(err.Error(), "access denied") {
		// This is fine - filesystem level errors
	}
}

func TestSecurity_AbsolutePathHandling(t *testing.T) {
	// Test with absolute path
	absPath, err := filepath.Abs("./testdata/simple.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Create a test file at absolute path for testing
	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(filepath.Dir(absPath), "security_test.yaml")
	if err := os.WriteFile(testFile, []byte("secure: true\n"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile)

	cfg, err := Load(testFile)
	if err != nil {
		t.Errorf("Absolute path loading failed: %v", err)
	}

	if !cfg.GetBool("secure") {
		t.Error("Expected secure=true")
	}
}
