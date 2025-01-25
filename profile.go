package konfig

import (
	"flag"
	"os"
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

func resetProfile() {
	activeProfile = ""
	once = sync.Once{}
	os.Args = []string{os.Args[0]}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}
