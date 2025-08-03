// Package konfig provides explicit, production-ready configuration management for Go applications.
//
// konfig focuses on simplicity, type safety, and explicit behavior. All configuration
// file paths must be provided explicitly - no magic auto-discovery.
//
// Basic usage:
//
//	cfg, err := konfig.Load("./config/app.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	port := cfg.GetString("server.port")
//
// Type-safe struct loading:
//
//	type Config struct {
//	    Port string `konfig:"server.port" default:"8080"`
//	}
//	var cfg Config
//	err := konfig.LoadInto("./config/app.yaml", &cfg)
//
// Profile-based configuration:
//
//	cfg, err := konfig.LoadWithProfile("./config/app.yaml", "dev")
package konfig

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Config provides type-safe access to configuration values
type Config interface {
	// Get returns the raw value and whether it exists
	Get(key string) (interface{}, bool)

	// Type-safe getters with sensible defaults
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetDuration(key string) time.Duration

	// GetStringWithDefault returns the value or default if not found
	GetStringWithDefault(key, defaultValue string) string
	GetIntWithDefault(key string, defaultValue int) int
	GetBoolWithDefault(key string, defaultValue bool) bool

	// Keys returns all available configuration keys
	Keys() []string
}

// config implements the Config interface
type config struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// ConfigError represents configuration-related errors with context
type ConfigError struct {
	Type    string // "file_not_found", "parse_error", "validation_error", "type_error"
	Path    string // File path or config key path
	Message string
	Cause   error
}

func (e *ConfigError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("konfig %s at %s: %s (%v)", e.Type, e.Path, e.Message, e.Cause)
	}
	return fmt.Sprintf("konfig %s at %s: %s", e.Type, e.Path, e.Message)
}

func (e *ConfigError) Unwrap() error {
	return e.Cause
}

// Load loads configuration from a single YAML file
//
// Example:
//
//	cfg, err := konfig.Load("./config/app.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	port := cfg.GetString("server.port")
func Load(filePath string) (Config, error) {
	if filePath == "" {
		return nil, &ConfigError{
			Type:    "validation_error",
			Path:    filePath,
			Message: "file path cannot be empty",
		}
	}

	return loadFromFile(filePath)
}

// LoadWithProfile loads base configuration and profile-specific overrides
//
// It loads the base file first, then looks for a profile-specific file
// with the pattern: base-{profile}.yaml
//
// Example:
//
//	cfg, err := konfig.LoadWithProfile("./config/app.yaml", "dev")
//	// Loads: ./config/app.yaml, then ./config/app-dev.yaml
func LoadWithProfile(filePath, profile string) (Config, error) {
	if filePath == "" {
		return nil, &ConfigError{
			Type:    "validation_error",
			Path:    filePath,
			Message: "file path cannot be empty",
		}
	}

	if profile == "" {
		return Load(filePath)
	}

	// Load base configuration
	cfg, err := loadFromFile(filePath)
	if err != nil {
		return nil, err
	}

	// Generate profile file path
	profilePath := generateProfilePath(filePath, profile)

	// Load profile configuration if it exists
	if fileExists(profilePath) {
		profileCfg, err := loadFromFile(profilePath)
		if err != nil {
			return nil, &ConfigError{
				Type:    "parse_error",
				Path:    profilePath,
				Message: "failed to load profile configuration",
				Cause:   err,
			}
		}

		// Merge profile config over base config
		cfg = mergeConfigs(cfg, profileCfg)
	}

	return cfg, nil
}

// LoadInto loads configuration into a struct using tags
//
// Struct fields should use `konfig:"key.path"` tags to map configuration keys.
// Default values can be specified with `default:"value"` tags.
//
// Example:
//
//	type Config struct {
//	    Port   int    `konfig:"server.port" default:"8080"`
//	    Host   string `konfig:"server.host" default:"localhost"`
//	    Debug  bool   `konfig:"debug" default:"false"`
//	}
//	var cfg Config
//	err := konfig.LoadInto("./config/app.yaml", &cfg)
func LoadInto(filePath string, target interface{}) error {
	cfg, err := Load(filePath)
	if err != nil {
		return err
	}

	return populateStruct(cfg, target)
}

// LoadIntoWithProfile loads configuration with profile support into a struct
func LoadIntoWithProfile(filePath, profile string, target interface{}) error {
	cfg, err := LoadWithProfile(filePath, profile)
	if err != nil {
		return err
	}

	return populateStruct(cfg, target)
}

// Implementation details

func loadFromFile(filePath string) (*config, error) {
	// Check if file exists and is readable
	if !fileExists(filePath) {
		return nil, &ConfigError{
			Type:    "file_not_found",
			Path:    filePath,
			Message: "configuration file not found",
		}
	}

	// Load and parse YAML
	configMap, err := parseYAMLFile(filePath)
	if err != nil {
		return nil, &ConfigError{
			Type:    "parse_error",
			Path:    filePath,
			Message: "failed to parse YAML file",
			Cause:   err,
		}
	}

	// Flatten nested keys into dot notation
	flatMap := flattenMap(configMap, "")

	// Process environment variable substitutions
	processedMap, err := processEnvSubstitutions(flatMap)
	if err != nil {
		return nil, &ConfigError{
			Type:    "parse_error",
			Path:    filePath,
			Message: "failed to process environment variable substitutions",
			Cause:   err,
		}
	}

	return &config{
		data: processedMap,
	}, nil
}

func generateProfilePath(basePath, profile string) string {
	dir := filepath.Dir(basePath)
	filename := filepath.Base(basePath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Try both extensions: same as base file, then the other YAML extension
	extensions := []string{ext}
	if ext == ".yml" {
		extensions = append(extensions, ".yaml")
	} else if ext == ".yaml" {
		extensions = append(extensions, ".yml")
	}

	for _, tryExt := range extensions {
		profileFilename := fmt.Sprintf("%s-%s%s", nameWithoutExt, profile, tryExt)
		profilePath := filepath.Join(dir, profileFilename)
		if fileExists(profilePath) {
			return profilePath
		}
	}

	// Fallback to first extension if nothing found
	profileFilename := fmt.Sprintf("%s-%s%s", nameWithoutExt, profile, extensions[0])
	return filepath.Join(dir, profileFilename)
}

func mergeConfigs(base, override *config) *config {
	result := &config{
		data: make(map[string]interface{}),
	}

	// Copy base config
	base.mu.RLock()
	for key, value := range base.data {
		result.data[key] = value
	}
	base.mu.RUnlock()

	// Override with profile config
	override.mu.RLock()
	for key, value := range override.data {
		result.data[key] = value
	}
	override.mu.RUnlock()

	return result
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Config interface implementation

func (c *config) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.data[key]
	return value, exists
}

func (c *config) GetString(key string) string {
	if value, exists := c.Get(key); exists {
		return fmt.Sprintf("%v", value)
	}
	return ""
}

func (c *config) GetInt(key string) int {
	if value, exists := c.Get(key); exists {
		if str := fmt.Sprintf("%v", value); str != "" {
			if i, err := strconv.Atoi(str); err == nil {
				return i
			}
		}
	}
	return 0
}

func (c *config) GetBool(key string) bool {
	if value, exists := c.Get(key); exists {
		if str := fmt.Sprintf("%v", value); str != "" {
			if b, err := strconv.ParseBool(str); err == nil {
				return b
			}
		}
	}
	return false
}

func (c *config) GetFloat64(key string) float64 {
	if value, exists := c.Get(key); exists {
		if str := fmt.Sprintf("%v", value); str != "" {
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				return f
			}
		}
	}
	return 0.0
}

func (c *config) GetDuration(key string) time.Duration {
	if value, exists := c.Get(key); exists {
		if str := fmt.Sprintf("%v", value); str != "" {
			if d, err := time.ParseDuration(str); err == nil {
				return d
			}
		}
	}
	return 0
}

func (c *config) GetStringWithDefault(key, defaultValue string) string {
	if value := c.GetString(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *config) GetIntWithDefault(key string, defaultValue int) int {
	if value, exists := c.Get(key); exists && fmt.Sprintf("%v", value) != "" {
		return c.GetInt(key)
	}
	return defaultValue
}

func (c *config) GetBoolWithDefault(key string, defaultValue bool) bool {
	if value, exists := c.Get(key); exists && fmt.Sprintf("%v", value) != "" {
		return c.GetBool(key)
	}
	return defaultValue
}

func (c *config) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

// populateStruct fills a struct using konfig tags
func populateStruct(cfg Config, target interface{}) error {
	if target == nil {
		return &ConfigError{
			Type:    "validation_error",
			Path:    "struct",
			Message: "target struct cannot be nil",
		}
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return &ConfigError{
			Type:    "validation_error",
			Path:    "struct",
			Message: "target must be a pointer to struct",
		}
	}

	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return &ConfigError{
			Type:    "validation_error",
			Path:    "struct",
			Message: "target must be a pointer to struct",
		}
	}

	return populateStructFields(cfg, elem, elem.Type(), "")
}

func populateStructFields(cfg Config, v reflect.Value, t reflect.Type, prefix string) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// Get konfig tag
		tag := field.Tag.Get("konfig")
		if tag == "" {
			// Handle nested structs without explicit tags
			if fieldValue.Kind() == reflect.Struct {
				nestedPrefix := prefix
				if prefix != "" {
					nestedPrefix = prefix + "."
				}
				nestedPrefix += strings.ToLower(field.Name)

				if err := populateStructFields(cfg, fieldValue, fieldValue.Type(), nestedPrefix); err != nil {
					return err
				}
			}
			continue
		}

		// Build full config key path
		configKey := tag
		if prefix != "" {
			configKey = prefix + "." + tag
		}

		// Handle nested structs
		if fieldValue.Kind() == reflect.Struct {
			// For nested structs, recursively populate using the config key as prefix
			if err := populateStructFields(cfg, fieldValue, fieldValue.Type(), configKey); err != nil {
				return err
			}
		} else {
			// Get default value
			defaultValue := field.Tag.Get("default")

			// Set scalar field value
			if err := setFieldValue(cfg, fieldValue, configKey, defaultValue); err != nil {
				return &ConfigError{
					Type:    "type_error",
					Path:    fmt.Sprintf("%s.%s", t.Name(), field.Name),
					Message: fmt.Sprintf("failed to set field from config key '%s'", configKey),
					Cause:   err,
				}
			}
		}
	}

	return nil
}

func setFieldValue(cfg Config, fieldValue reflect.Value, configKey, defaultValue string) error {
	// Get value from config or use default
	var strValue string
	if value, exists := cfg.Get(configKey); exists && value != nil {
		strValue = fmt.Sprintf("%v", value)
	} else {
		strValue = defaultValue
	}

	// Skip if no value available
	if strValue == "" {
		return nil
	}

	// Set value based on field type
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(strValue)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Handle time.Duration specially
		if fieldValue.Type() == reflect.TypeOf(time.Duration(0)) {
			if d, err := time.ParseDuration(strValue); err == nil {
				fieldValue.Set(reflect.ValueOf(d))
			} else {
				return fmt.Errorf("cannot convert '%s' to duration: %w", strValue, err)
			}
		} else if i, err := strconv.ParseInt(strValue, 10, 64); err == nil {
			fieldValue.SetInt(i)
		} else {
			return fmt.Errorf("cannot convert '%s' to int: %w", strValue, err)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if u, err := strconv.ParseUint(strValue, 10, 64); err == nil {
			fieldValue.SetUint(u)
		} else {
			return fmt.Errorf("cannot convert '%s' to uint: %w", strValue, err)
		}

	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(strValue, fieldValue.Type().Bits()); err == nil {
			fieldValue.SetFloat(f)
		} else {
			return fmt.Errorf("cannot convert '%s' to float: %w", strValue, err)
		}

	case reflect.Bool:
		if b, err := strconv.ParseBool(strValue); err == nil {
			fieldValue.SetBool(b)
		} else {
			return fmt.Errorf("cannot convert '%s' to bool: %w", strValue, err)
		}

	case reflect.Struct:
		// Handle time.Duration specially
		if fieldValue.Type() == reflect.TypeOf(time.Duration(0)) {
			if d, err := time.ParseDuration(strValue); err == nil {
				fieldValue.Set(reflect.ValueOf(d))
			} else {
				return fmt.Errorf("cannot convert '%s' to duration: %w", strValue, err)
			}
		} else {
			// Nested struct - recursive population
			return populateStructFields(cfg, fieldValue, fieldValue.Type(), configKey)
		}

	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
	}

	return nil
}
