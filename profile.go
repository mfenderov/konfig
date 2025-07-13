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

// ResetProfileInitialized resets the profile initialization state.
// This function is primarily used in tests to ensure clean state between test runs.
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

// IsProdProfile returns true if the current active profile is "prod".
//
// The profile is determined by command-line flags (-p or --profile).
//
// Example:
//
//	if konfig.IsProdProfile() {
//	    // Production-specific logic
//	    fmt.Println("Running in production mode")
//	}
func IsProdProfile() bool {
	return getProfile() == prodProfile
}

// IsDevProfile returns true if the current active profile is "dev".
//
// The profile is determined by command-line flags (-p or --profile).
//
// Example:
//
//	if konfig.IsDevProfile() {
//	    // Development-specific logic
//	    fmt.Println("Running in development mode")
//	}
func IsDevProfile() bool {
	return getProfile() == devProfile
}

// IsProfile returns true if the current active profile matches the given name.
//
// This is useful for checking custom profile names beyond "dev" and "prod".
//
// Example:
//
//	if konfig.IsProfile("staging") {
//	    // Staging-specific logic
//	    fmt.Println("Running in staging mode")
//	}
func IsProfile(profile string) bool {
	return getProfile() == profile
}

// GetProfile returns the currently active profile name.
//
// The profile is determined by command-line flags (-p or --profile).
// Returns an empty string if no profile is active.
//
// Example:
//
//	profile := konfig.GetProfile()
//	if profile != "" {
//	    fmt.Printf("Active profile: %s\n", profile)
//	} else {
//	    fmt.Println("No active profile")
//	}
func GetProfile() string {
	return getProfile()
}
