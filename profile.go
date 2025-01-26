package konfig

import (
	"github.com/spf13/pflag"
)

const devProfile = "dev"
const prodProfile = "prod"

var parsedProfile string

func init() {
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.StringVarP(&parsedProfile, "profile", "p", "", "Application profile")
}

func getProfile() string {
	parsedProfile = ""
	pflag.Parse()
	return parsedProfile
}

func IsProdProfile() bool {
	return getProfile() == prodProfile
}

func IsDevProfile() bool {
	return getProfile() == devProfile
}

func IsProfile(profile string) bool {
	return getProfile() == profile
}
