package konfig

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func localConfigMapFromFile(pathToConfigFile string) (map[string]interface{}, error) {
	rootPath, _ := findRootPath()
	configFile, err := os.ReadFile(filepath.Join(rootPath, pathToConfigFile))
	if err != nil {
		return nil, fmt.Errorf("error reading configuration file: %w", err)
	}
	configFile = []byte(os.ExpandEnv(string(configFile)))

	configMap := make(map[string]interface{})
	err = yaml.Unmarshal(configFile, configMap)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling configuration file: %w", err)
	}
	return configMap, nil
}

func postProcessConfig(resultConfigMap map[string]interface{}) error {
	for key, value := range resultConfigMap {
		if value == nil {
			return fmt.Errorf("value for 'key' %s is nil", key)
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

func findRootPath() (string, error) {
	cfg := &packages.Config{Mode: packages.NeedModule}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return "", errors.Wrap(err, "failed to load packages")
	}

	if len(pkgs) > 0 && pkgs[0].Module != nil {
		return pkgs[0].Module.Dir, nil
	} else {
		return "", errors.Wrap(err, "failed to find root path")
	}
}
