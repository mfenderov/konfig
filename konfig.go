package konfig

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
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

const defaultConfigFolder = "resources"
const defaultConfigFileName = "application"

var defaultConfigFileExtensions = []string{".yaml", ".yml"}
var defaultProfile = prodProfile

func Load() {
	profileSuffix := ""
	if !IsProdProfile() {
		profileSuffix = "-" + getProfile()
	}
	for _, ext := range defaultConfigFileExtensions {
		configFilePath := fmt.Sprintf("%s/%s%s%s", defaultConfigFolder, defaultConfigFileName, profileSuffix, ext)
		if _, err := os.Stat(configFilePath); err == nil {
			LoadFrom(configFilePath)
			return
		}
	}
}

func LoadFrom(pathToConfigFile string) error {
	configMap, err := localConfigMapFromFile(pathToConfigFile)
	if err != nil {
		return errors.Wrapf(err, "error loading config from %s", pathToConfigFile)
	}
	resultConfigMap := buildEnvVariables(configMap)
	return postProcessConfig(resultConfigMap)
}
