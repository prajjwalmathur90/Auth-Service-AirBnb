package app

import (
	config "AuthInGo/config/env"
	"AuthInGo/router"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Addr string
}

func NewConfig() *Config {
	port := config.GetString("PORT", "8080")

	return &Config {
		Addr: ":" + port,
	}
}

type Application struct {
	Config Config
}

func NewApplication(config *Config) *Application {
	return &Application {
		Config: *config,
	}
}

func (app *Application) Run() error {
	server := &http.Server{
		Addr: app.Config.Addr,
		Handler: router.SetupRouter(),
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server is running at", app.Config.Addr)

	return server.ListenAndServe()
}