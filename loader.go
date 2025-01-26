package konfig

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func localConfigMapFromFile(pathToConfigFile string) (map[string]interface{}, error) {
	configFile, err := readConfigFile(pathToConfigFile)
	if err != nil {
		return nil, fmt.Errorf("error reading configuration file: %w", err)
	}
	configFile = os.Expand(configFile, enrichValue)

	configMap := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(configFile), configMap)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling configuration file: %w", err)
	}
	return configMap, nil
}

func readConfigFile(pathToConfigFile string) (string, error) {
	rootPath, err := findRootPath()
	if err != nil {
		return "", fmt.Errorf("error finding root path: %w", err)
	}
	configFile, err := os.ReadFile(filepath.Join(rootPath, pathToConfigFile))
	if err != nil {
		return "", fmt.Errorf("error reading configuration file: %w", err)
	}
	return string(configFile), nil
}

func enrichValue(value string) string {
	index := strings.Index(value, ":")
	if index == -1 {
		return GetEnv(value)
	}

	key := value[:index]
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

func processConfig(prefix string, configMap map[string]interface{}, resultConfigMap *map[string]interface{}) {
	for key, value := range configMap {
		updatedPrefix := updatePrefix(prefix, key)
		switch v := value.(type) {
		case map[string]interface{}:
			processConfig(updatedPrefix, v, resultConfigMap)
		case []interface{}:
			itemValue := make([]interface{}, len(v))
			for i, item := range v {
				itemValue[i] = item
			}
			(*resultConfigMap)[updatedPrefix] = itemValue
		default:
			(*resultConfigMap)[updatedPrefix] = v
		}
	}
}

func buildEnvVariables(configMap map[string]interface{}) map[string]interface{} {
	resultConfigMap := make(map[string]interface{})
	processConfig("", configMap, &resultConfigMap)
	return resultConfigMap
}

func updatePrefix(prefix string, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%s.%s", prefix, key)
}

type root struct {
	path string
	err  error
	once sync.Once
}

func findRootPath() (string, error) {
	root := root{}
	root.once.Do(func() {
		dir, err := os.Getwd()
		if err != nil {
			root.path, root.err = "", errors.Wrap(err, "failed to get current directory")
		}

		for {
			if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
				root.path, root.err = dir, nil
				return
			}

			parentDir := filepath.Dir(dir)
			if parentDir == dir {
				break
			}
			dir = parentDir
		}

		root.path, root.err = "", errors.New("failed to find root path")
	})
	return root.path, root.err
}
