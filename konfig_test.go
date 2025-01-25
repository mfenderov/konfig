package konfig_test

import (
	"testing"

	"github.com/mfenderov/konfig"

	"github.com/stretchr/testify/assert"
)

func TestLoadSimpleConfiguration(t *testing.T) {
	konfig.ClearEnv()

	_ = konfig.SetEnv("POSTGRES_PASSWORD", "12345")
	_ = konfig.SetEnv("POSTGRES_HOST", "localhost")

	err := konfig.LoadFrom("test-data/application.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "12345", konfig.GetEnv("db.password"))
	assert.Equal(t, "https://localhost", konfig.GetEnv("db.host"))
	assert.Equal(t, "5432", konfig.GetEnv("db.port"))
	assert.Equal(t, "user", konfig.GetEnv("db.user"))
	assert.Equal(t, "8080", konfig.GetEnv("server.port"))
	assert.Equal(t, "debug", konfig.GetEnv("log.level"))
	assert.Equal(t, "[one two three]", konfig.GetEnv("list"))
}

func TestLoadConfigurationWithEmptyProperties(t *testing.T) {
	konfig.ClearEnv()

	err := konfig.LoadFrom("test-data/application-empty.yaml")
	assert.EqualError(t, err, "property 'db.password' is nil")
}

func TestLoadConfigurationWithEmptyPropertyWithDefault(t *testing.T) {
	konfig.ClearEnv()

	err := konfig.LoadFrom("test-data/application-empty-with-default.yaml")
	assert.NoError(t, err)

	assert.Equal(t, "postgres", konfig.GetEnv("db.password"))
	assert.Equal(t, "https://", konfig.GetEnv("db.host"))
}
