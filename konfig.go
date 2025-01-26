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
const fullFileNameFormat = "%s/%s"
const filePathWithinProjectFormat = "%s/%s%s%s"

var defaultConfigFileExtensions = []string{".yaml", ".yml"}
var defaultProfile = prodProfile

func Load() error {
	profileSuffix := ""
	if getProfile() == defaultProfile {
		profileSuffix = "-" + getProfile()
	}
	path, err := findRootPath()
	if err != nil {
		return errors.Wrap(err, "error finding root path")
	}

	for _, ext := range defaultConfigFileExtensions {
		filePathWithinProject := fmt.Sprintf(filePathWithinProjectFormat, defaultConfigFolder, defaultConfigFileName, profileSuffix, ext)
		configFilePath := fmt.Sprintf(fullFileNameFormat, path, filePathWithinProject)
		if _, err := os.Stat(configFilePath); err == nil {
			return LoadFrom(filePathWithinProject)
		}
	}
	return errors.Errorf("config file not found in '%s'. Default file name is '%s', and suppoted extensions are %v", defaultConfigFolder, defaultConfigFileName, defaultConfigFileExtensions)
}

func LoadFrom(pathToConfigFile string) error {
	configMap, err := localConfigMapFromFile(pathToConfigFile)
	if err != nil {
		return errors.Wrapf(err, "error loading config from %s", pathToConfigFile)
	}
	resultConfigMap := buildEnvVariables(configMap)
	return postProcessConfig(resultConfigMap)
}
