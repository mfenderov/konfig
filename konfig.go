package konfig

import (
	"fmt"
	"os"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}

func SetEnv(key string, value string) error {
	return os.Setenv(key, value)
}

func ClearEnv() {
	os.Clearenv()
}

func LoadConfiguration(pathToConfigFile string) error {
	configMap, err := localConfigMapFromFile(pathToConfigFile)
	if err != nil {
		return fmt.Errorf("error reading configuration file: %w", err)
	}
	resultConfigMap := buildEnvVariables(configMap)
	return postProcessConfig(resultConfigMap)
}
