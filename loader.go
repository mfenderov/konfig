package konfig

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	rootPathOnce sync.Once
	rootPath     string
	rootPathErr  error
)

func localConfigMapFromFile(pathToConfigFile string) (map[string]any, error) {
	configFile, err := readConfigFile(pathToConfigFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read config file")
	}
	configFile = os.Expand(configFile, enrichValue)

	configMap := make(map[string]any)
	err = yaml.Unmarshal([]byte(configFile), &configMap)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config file")
	}
	return configMap, nil
}

func readConfigFile(pathToConfigFile string) (string, error) {
	if pathToConfigFile == "" {
		return "", errors.New("config file path cannot be empty")
	}

	rootPath, err := findRootPath()
	if err != nil {
		return "", errors.Wrap(err, "failed to find root path")
	}

	fullPath := filepath.Join(rootPath, pathToConfigFile)
	if !strings.HasSuffix(fullPath, ".yml") && !strings.HasSuffix(fullPath, ".yaml") {
		return "", errors.New("config file must have .yml or .yaml extension")
	}

	configFile, err := os.ReadFile(fullPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to read config file")
	}
	return string(configFile), nil
}

func enrichValue(value string) string {
	if value == "" {
		return ""
	}

	index := strings.Index(value, ":")
	if index == -1 {
		// No default value provided
		envValue := GetEnv(value)
		if envValue == "" {
			slog.Warn("Environment variable not found and no default provided", "key", value)
		}
		return envValue
	}

	key := value[:index]
	if key == "" {
		slog.Warn("Empty environment variable key with default value", "value", value)
		return value[index+1:] // Return default value
	}

	defaultValue := value[index+1:]
	envValue := GetEnv(key)
	if envValue != "" {
		return envValue
	}
	return defaultValue
}

func postProcessConfig(resultConfigMap map[string]interface{}) error {
	for key, value := range resultConfigMap {
		if value == nil {
			return fmt.Errorf("property '%s' is nil", key)
		}
		err := SetEnv(key, fmt.Sprintf("%v", value))
		if err != nil {
			return fmt.Errorf("error setting environment variable: %w", err)
		}
	}
	return nil
}

func processConfig(prefix string, configMap map[string]any, resultConfigMap *map[string]any) error {
	if configMap == nil {
		return errors.New("config map cannot be nil")
	}
	if resultConfigMap == nil {
		return errors.New("result config map cannot be nil")
	}

	for key, value := range configMap {
		if key == "" {
			continue // Skip empty keys
		}
		updatedPrefix := updatePrefix(prefix, key)
		switch v := value.(type) {
		case map[string]any:
			if v != nil {
				if err := processConfig(updatedPrefix, v, resultConfigMap); err != nil {
					return errors.Wrapf(err, "failed to process config at key '%s'", key)
				}
			}
		case []any:
			itemValue := make([]any, len(v))
			copy(itemValue, v)
			(*resultConfigMap)[updatedPrefix] = itemValue
		default:
			(*resultConfigMap)[updatedPrefix] = v
		}
	}
	return nil
}

func buildEnvVariables(configMap map[string]any) (map[string]any, error) {
	resultConfigMap := make(map[string]any)
	if err := processConfig("", configMap, &resultConfigMap); err != nil {
		return nil, errors.Wrap(err, "failed to process config")
	}
	return resultConfigMap, nil
}

func updatePrefix(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%s.%s", prefix, key)
}
