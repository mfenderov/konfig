// Package konfig provides a Spring Framework-inspired configuration management system for Go applications.
//
// konfig supports YAML-based configuration files, profile-specific configurations,
// environment variable substitution, and struct-based configuration loading with type safety.
//
// Basic usage:
//
//	err := konfig.Load()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	port := os.Getenv("server.port")
//
// Struct-based configuration:
//
//	type Config struct {
//	    Port string `konfig:"server.port" default:"8080"`
//	}
//	var cfg Config
//	err := konfig.LoadInto(&cfg)
package konfig

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/pkg/errors"
)

// GetEnv retrieves the value of the environment variable named by the key.
// This is a convenience wrapper around os.Getenv for consistency with konfig's API.
func GetEnv(key string) string {
	return os.Getenv(key)
}

// SetEnv sets the value of the environment variable named by the key.
// This is a convenience wrapper around os.Setenv for consistency with konfig's API.
func SetEnv(key string, value string) error {
	return os.Setenv(key, value)
}

// ClearEnv deletes all environment variables.
// This is primarily used for testing purposes and should be used with caution.
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
// Load initializes konfig by loading configuration from default sources.
//
// It first loads the base configuration file (resources/application.yaml or .yml),
// then loads any active profile-specific configuration (e.g., resources/application-dev.yaml).
// Profile-specific values override base configuration values.
//
// Configuration keys are flattened into dot notation and set as environment variables,
// making them accessible via os.Getenv() or the GetEnv() function.
//
// The active profile is determined by command-line flags (-p or --profile).
//
// Example:
//
//	err := konfig.Load()
//	if err != nil {
//	    log.Fatal("Failed to load configuration:", err)
//	}
//	
//	port := os.Getenv("server.port")
//	dbURL := os.Getenv("database.url")
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

// LoadFrom loads configuration from a specific YAML file.
//
// Unlike Load(), this function only loads the specified file and does not
// attempt to load default application.yaml or profile-specific files.
//
// The configuration keys are flattened and set as environment variables,
// just like with Load().
//
// Example:
//
//	err := konfig.LoadFrom("config/custom.yaml")
//	if err != nil {
//	    log.Fatal("Failed to load custom config:", err)
//	}
//	
//	customValue := os.Getenv("custom.setting")
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
		var relativeConfigPath string
		if profileSuffix == "" {
			relativeConfigPath = fmt.Sprintf(format, defaultConfigFolder, defaultConfigFileName, ext)
		} else {
			relativeConfigPath = fmt.Sprintf(format, defaultConfigFolder, defaultConfigFileName, profileSuffix, ext)
		}

		configFilePath := fmt.Sprintf(fullFileNameFormat, path, relativeConfigPath)
		if _, err := os.Stat(configFilePath); err == nil {
			return LoadFrom(configFilePath)
		}
	}
	return nil
}
