package konfig

import "flag"

const devProfile = "dev"
const prodProfile = "prod"

func IsProdProfile() bool {
	return getProfile() == prodProfile
}

func IsDevProfile() bool {
	return getProfile() == devProfile
}

func IsProfile(profile string) bool {
	return getProfile() == profile
}

func getProfile() string {
	profile := flag.String("p", "", "profile to use")
	flag.Parse()
	if profile != nil {
		return *profile
	}
	return ""
}
