package konfig

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfile_ShouldReturnEmtpyString(t *testing.T) {
	profile := getProfile()
	assert.Empty(t, profile)
}

func TestProfile_ShouldReturnDevTrue(t *testing.T) {
	os.Args = []string{os.Args[0], "-p", "dev"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	profile := IsDevProfile()
	assert.True(t, profile)
}

func TestProfile_ShouldReturnProdTrue(t *testing.T) {
	os.Args = []string{os.Args[0], "-p", "prod"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	profile := IsProdProfile()
	assert.True(t, profile)
}

func TestProfile_ShouldReturnCustomerProfileTrue(t *testing.T) {
	os.Args = []string{os.Args[0], "-p", "test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	profile := IsProfile("test")
	assert.True(t, profile)
}

func TestProfile_ShouldReturnFalse(t *testing.T) {
	os.Args = []string{os.Args[0]}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	profile := IsProfile("dev123")
	assert.False(t, profile)
}
