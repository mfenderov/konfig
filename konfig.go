package konfig

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
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
	slog.Debug("Initializing konfig...")
	defer func() { slog.Debug("Konfig initialized", "duration", time.Since(now).String()) }()
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
	path, err := findRootPath()
	if err != nil {
		return errors.Wrap(err, "error finding root path")
	}

	for _, ext := range defaultConfigFileExtensions {
		filePathWithinProject := fmt.Sprintf(filePathWithProfileFormat, defaultConfigFolder, defaultConfigFileName, profileSuffix, ext)
		configFilePath := fmt.Sprintf(fullFileNameFormat, path, filePathWithinProject)
		if _, err := os.Stat(configFilePath); err == nil {
			return LoadFrom(filePathWithinProject)
		}
	}
	return nil
}

var setupLoggerOnce sync.Once

func setupLogger() {
	setupLoggerOnce.Do(func() {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	})
}

func LoadFrom(pathToConfigFile string) error {
	configMap, err := localConfigMapFromFile(pathToConfigFile)
	if err != nil {
		return errors.Wrapf(err, "error loading config from %s", pathToConfigFile)
	}
	resultConfigMap := buildEnvVariables(configMap)
	return postProcessConfig(resultConfigMap)
}

func loadDefault() error {
	now := time.Now()
	slog.Debug("Loading default config", "folder", defaultConfigFolder, "file", defaultConfigFileName, "extensions", defaultConfigFileExtensions)
	defer func() { slog.Debug("Default config loaded", "duration", time.Since(now).String()) }()
	path, err := findRootPath()
	if err != nil {
		return errors.Wrap(err, "error finding root path")
	}

	for _, ext := range defaultConfigFileExtensions {
		filePathWithinProject := fmt.Sprintf(filePathWithoutProfileFormat, defaultConfigFolder, defaultConfigFileName, ext)
		absoluteConfigFilePath := fmt.Sprintf(fullFileNameFormat, path, filePathWithinProject)
		if _, err := os.Stat(absoluteConfigFilePath); err == nil {
			return LoadFrom(filePathWithinProject)
		}
	}
	return nil
}
