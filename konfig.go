package konfig

import (
	"fmt"
	"log/slog"
	"os"
	"time"

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
const filePathWithProfileFormat = "%s/%s-%s.%s"
const filePathWithoutProfileFormat = "%s/%s.%s"

var defaultConfigFileExtensions = []string{"yaml", "yml"}

func init() {
	setupLogger()
}
func Load() error {
	now := time.Now()
	slog.Info("Initializing konfig...")
	defer func() { slog.Info("Konfig initialized", "duration", time.Since(now).String()) }()
	err := loadDefault()
	if err != nil {
		return errors.Wrap(err, "error loading default config")
	}

	return loadProfiled()
}

func loadProfiled() error {
	now := time.Now()
	slog.Debug("Loading profiled config")
	defer func() { slog.Debug("Profiled config loaded", "duration", time.Since(now).String()) }()

	profileSuffix := getProfile()
	if profileSuffix == "" {
		return nil
	}

	return loadConfigWithFormat(filePathWithProfileFormat, profileSuffix)
}

func setupLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
}

func LoadFrom(pathToConfigFile string) error {
	configMap, err := localConfigMapFromFile(pathToConfigFile)
	if err != nil {
		return errors.Wrapf(err, "error loading config from %s", pathToConfigFile)
	}
	resultConfigMap, err := buildEnvVariables(configMap)
	if err != nil {
		return errors.Wrap(err, "error building environment variables")
	}
	return postProcessConfig(resultConfigMap)
}

func loadDefault() error {
	now := time.Now()
	slog.Debug("Loading default config", "folder", defaultConfigFolder, "file", defaultConfigFileName, "extensions", defaultConfigFileExtensions)
	defer func() { slog.Debug("Default config loaded", "duration", time.Since(now).String()) }()

	return loadConfigWithFormat(filePathWithoutProfileFormat, "")
}

// loadConfigWithFormat loads configuration using the specified format and profile suffix
func loadConfigWithFormat(format string, profileSuffix string) error {
	path, err := findRootPath()
	if err != nil {
		return errors.Wrap(err, "error finding root path")
	}

	for _, ext := range defaultConfigFileExtensions {
		var filePathWithinProject string
		if profileSuffix == "" {
			filePathWithinProject = fmt.Sprintf(format, defaultConfigFolder, defaultConfigFileName, ext)
		} else {
			filePathWithinProject = fmt.Sprintf(format, defaultConfigFolder, defaultConfigFileName, profileSuffix, ext)
		}

		configFilePath := fmt.Sprintf(fullFileNameFormat, path, filePathWithinProject)
		if _, err := os.Stat(configFilePath); err == nil {
			return LoadFrom(configFilePath)
		}
	}
	return nil
}
