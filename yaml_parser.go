package konfig

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	maxFileSize     = 10 * 1024 * 1024 // 10MB max file size
	maxNestingDepth = 32               // Maximum YAML nesting depth
)

// parseYAMLFile reads and parses a YAML file into a map with security validations
func parseYAMLFile(filePath string) (map[string]interface{}, error) {
	// Security: Prevent path traversal attacks before cleaning
	if strings.Contains(filePath, "..") {
		return nil, fmt.Errorf("path traversal not allowed: %s", filePath)
	}

	// Security: Clean the file path after validation
	cleanPath := filepath.Clean(filePath)

	// Security: Check file info before reading
	fileInfo, err := os.Stat(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access file: %w", err)
	}

	// Security: Enforce file size limit
	if fileInfo.Size() > maxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max: %d)", fileInfo.Size(), maxFileSize)
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Security: Validate YAML complexity
	if err := validateYAMLComplexity(result, 0); err != nil {
		return nil, fmt.Errorf("YAML too complex: %w", err)
	}

	return result, nil
}

// validateYAMLComplexity prevents deeply nested YAML from causing stack overflow
func validateYAMLComplexity(data interface{}, depth int) error {
	if depth > maxNestingDepth {
		return fmt.Errorf("nesting depth exceeds maximum of %d", maxNestingDepth)
	}

	switch v := data.(type) {
	case map[string]interface{}:
		for _, value := range v {
			if err := validateYAMLComplexity(value, depth+1); err != nil {
				return err
			}
		}
	case []interface{}:
		for _, item := range v {
			if err := validateYAMLComplexity(item, depth+1); err != nil {
				return err
			}
		}
	}

	return nil
}

// flattenMap converts nested maps into dot-notation keys
func flattenMap(m map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range m {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// Recursively flatten nested maps
			nested := flattenMap(v, fullKey)
			for nestedKey, nestedValue := range nested {
				result[nestedKey] = nestedValue
			}
		default:
			result[fullKey] = value
		}
	}

	return result
}

// processEnvSubstitutions processes ${VAR} and ${VAR:default} substitutions
func processEnvSubstitutions(m map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Regular expression to match ${VAR} or ${VAR:default}
	envVarRegex := regexp.MustCompile(`\$\{([^}:]+)(?::([^}]*))?\}`)

	for key, value := range m {
		strValue := fmt.Sprintf("%v", value)

		// Process all environment variable substitutions in the string
		processedValue := envVarRegex.ReplaceAllStringFunc(strValue, func(match string) string {
			matches := envVarRegex.FindStringSubmatch(match)
			if len(matches) < 2 {
				return match // Should not happen, but safety first
			}

			envVar := matches[1]
			defaultVal := ""
			if len(matches) > 2 {
				defaultVal = matches[2]
			}

			// Get environment variable value
			if envValue := os.Getenv(envVar); envValue != "" {
				return envValue
			}

			// Use default value if environment variable is not set
			return defaultVal
		})

		// Convert back to appropriate type if possible
		if processedValue != strValue {
			// String was modified, keep as string
			result[key] = processedValue
		} else {
			// String was not modified, keep original type
			result[key] = value
		}
	}

	return result, nil
}
