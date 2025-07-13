package main

import (
	"os"

	"github.com/mfenderov/konfig"
)

func main() {
	err := konfig.Load()
	if err != nil {
		os.Exit(1)
	}
}
