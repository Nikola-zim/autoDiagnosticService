package main

import (
	"log"

	"autoDiagnosticService/config"
	"autoDiagnosticService/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
