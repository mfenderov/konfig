package main

import (
	"testing"

	"github.com/mfenderov/konfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WithoutProfile(t *testing.T) {
	// Test loading base configuration without profile
	cfg, err := konfig.Load("./resources/application.yml")
	require.NoError(t, err)

	actual := cfg.GetString("some.property.value")
	assert.Equal(t, "123", actual)
}

func Test_WithProfile(t *testing.T) {
	// Test loading configuration with dev profile
	cfg, err := konfig.LoadWithProfile("./resources/application.yml", "dev")
	require.NoError(t, err)

	actual := cfg.GetString("some.property.value")
	assert.Equal(t, "777", actual)
}
