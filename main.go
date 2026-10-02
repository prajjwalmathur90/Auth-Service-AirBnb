package main

import (
	"AuthInGo/app"
	config "AuthInGo/config/env"
	"log"
)

func main() {
	config.Load()

	cfg := app.NewConfig()
	application := app.NewApplication(*cfg)

	if err := application.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
