package konfig

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// LoadInto loads configuration into a struct using konfig tags
func LoadInto(config interface{}) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("config must be a pointer to struct")
	}
	
	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to struct")
	}
	
	// First load the YAML configuration (existing functionality)
	err := Load()
	if err != nil {
		return fmt.Errorf("failed to load konfig: %w", err)
	}
	
	// Then populate the struct
	return populateStruct(elem, "")
}

// populateStruct recursively populates struct fields from environment variables
func populateStruct(v reflect.Value, prefix string) error {
	t := v.Type()
	
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		
		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}
		
		konfigTag := field.Tag.Get("konfig")
		if konfigTag == "" {
			continue
		}
		
		fullPath := buildPath(prefix, konfigTag)
		
		if err := setFieldFromEnv(fieldValue, field, fullPath); err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}
	}
	
	return nil
}

// buildPath constructs the full environment variable path
func buildPath(prefix, tag string) string {
	if prefix == "" {
		return tag
	}
	return prefix + "." + tag
}

// setFieldFromEnv sets a field value from environment variable
func setFieldFromEnv(fieldValue reflect.Value, field reflect.StructField, envKey string) error {
	// Get value from environment or use default
	envValue := os.Getenv(envKey)
	if envValue == "" {
		envValue = field.Tag.Get("default")
	}
	
	// Handle different field types
	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(envValue)
	case reflect.Int, reflect.Int32, reflect.Int64:
		if envValue == "" {
			return nil // Leave as zero value
		}
		intVal, err := strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse '%s' as integer: %w", envValue, err)
		}
		fieldValue.SetInt(intVal)
	case reflect.Bool:
		if envValue == "" {
			return nil // Leave as zero value
		}
		boolVal, err := strconv.ParseBool(envValue)
		if err != nil {
			return fmt.Errorf("cannot parse '%s' as boolean: %w", envValue, err)
		}
		fieldValue.SetBool(boolVal)
	case reflect.Float32, reflect.Float64:
		if envValue == "" {
			return nil // Leave as zero value
		}
		floatVal, err := strconv.ParseFloat(envValue, 64)
		if err != nil {
			return fmt.Errorf("cannot parse '%s' as float: %w", envValue, err)
		}
		fieldValue.SetFloat(floatVal)
	case reflect.Struct:
		// Handle nested structs recursively
		return populateStruct(fieldValue, envKey)
	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
	}
	
	return nil
}