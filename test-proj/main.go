package main

import (
	"github.com/mfenderov/konfig"
	"os"
)

func main() {
	err := konfig.Load()
	if err != nil {
		os.Exit(1)
	}
}
