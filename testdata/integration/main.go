package main

import (
	"log"
	"os"

	"github.com/mfenderov/konfig"
)

// config represents our application configuration
var appConfig konfig.Config

func main() {
	// Load configuration explicitly from resources directory
	cfg, err := konfig.Load("./resources/application.yml")
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
		os.Exit(1)
	}

	appConfig = cfg
}

// GetConfig returns the loaded configuration for testing access
func GetConfig() konfig.Config {
	return appConfig
}
