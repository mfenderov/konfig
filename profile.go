package konfig

import (
	"github.com/spf13/pflag"
)

const devProfile = "dev"
const prodProfile = "prod"

var parsedProfile string
var profileInitialized bool

func init() {
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.StringVarP(&parsedProfile, "profile", "p", "", "Application profile")
}

// ResetProfileInitialized is used in tests to reset the profile initialization state
func ResetProfileInitialized() {
	profileInitialized = false
	parsedProfile = ""
}

func getProfile() string {
	if !profileInitialized {
		pflag.Parse()
		profileInitialized = true
	}
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

// GetProfile returns the currently active profile name
func GetProfile() string {
	return getProfile()
}
