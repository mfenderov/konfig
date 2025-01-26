package konfig

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfile_ShouldReturnEmptyString(t *testing.T) {
	resetCommandLineFlags()

	profile := getProfile()
	assert.Empty(t, profile)
}

func TestProfile_ShouldReturnDevTrue(t *testing.T) {
	resetCommandLineFlags()
	setCommandLineFlag("dev")

	profile := IsDevProfile()
	assert.True(t, profile)
}

func TestProfile_ShouldReturnProdTrue(t *testing.T) {
	resetCommandLineFlags()
	setCommandLineFlag("prod")

	profile := IsProdProfile()
	assert.True(t, profile)
}

func TestProfile_ShouldReturnCustomerProfileTrue(t *testing.T) {
	resetCommandLineFlags()
	setCommandLineFlag("test")

	profile := IsProfile("test")
	assert.True(t, profile)
}

func TestProfile_ShouldReturnFalse(t *testing.T) {
	resetCommandLineFlags()
	profile := IsProfile("dev123")
	assert.False(t, profile)
}

func resetCommandLineFlags() {
	os.Args = []string{os.Args[0]}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func setCommandLineFlag(p string) {
	os.Args = []string{os.Args[0], "-p", p}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
