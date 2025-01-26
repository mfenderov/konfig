package main

import (
	"flag"
	"github.com/mfenderov/konfig"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_WithoutProfile(t *testing.T) {
	os.Args = []string{os.Args[0]}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	main()

	actual := konfig.GetEnv("some.property.value")
	assert.Equal(t, "123", actual)
}

func Test_WithProfile(t *testing.T) {
	os.Args = []string{os.Args[0], "-p", "dev"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	main()

	actual := konfig.GetEnv("some.property.value")
	assert.Equal(t, "777", actual)
}
