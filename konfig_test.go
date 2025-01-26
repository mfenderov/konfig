package konfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFrom_SimpleConfiguration(t *testing.T) {
	ClearEnv()

	_ = SetEnv("POSTGRES_PASSWORD", "12345")
	_ = SetEnv("POSTGRES_HOST", "localhost")

	err := LoadFrom("test-data/application.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "12345", GetEnv("db.password"))
	assert.Equal(t, "https://localhost", GetEnv("db.host"))
	assert.Equal(t, "5432", GetEnv("db.port"))
	assert.Equal(t, "user", GetEnv("db.user"))
	assert.Equal(t, "8080", GetEnv("server.port"))
	assert.Equal(t, "debug", GetEnv("log.level"))
	assert.Equal(t, "[one two three]", GetEnv("list"))
}

func TestLoadFrom_ConfigurationWithEmptyProperties(t *testing.T) {
	ClearEnv()

	err := LoadFrom("test-data/application-empty.yaml")
	assert.EqualError(t, err, "property 'db.password' is nil")
}

func TestLoadFrom_ConfigurationWithEmptyPropertyWithDefault(t *testing.T) {
	ClearEnv()

	err := LoadFrom("test-data/application-empty-with-default.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "postgres", GetEnv("db.password"))
	assert.Equal(t, "https://", GetEnv("db.host"))
}

func TestLoad_WithoutProfile(t *testing.T) {
	ClearEnv()

	err := Load()
	assert.NoError(t, err)

	port := GetEnv("server.port")
	assert.Equal(t, "1234", port)
}

func TestLoad_WithTestProfile(t *testing.T) {
	ClearEnv()
	resetCommandLineFlags()

	setCommandLineFlag("test")

	err := Load()
	assert.NoError(t, err)

	port := GetEnv("server.port")
	assert.Equal(t, "4321", port)
}
