package konfig

import (
	"flag"
	"sync"
)

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

var activeProfile string
var once sync.Once

func getProfile() string {
	once.Do(func() {
		profile := flag.String("p", "", "profile to use")
		flag.Parse()
		if profile != nil {
			activeProfile = *profile
		}
	})
	return activeProfile
}
