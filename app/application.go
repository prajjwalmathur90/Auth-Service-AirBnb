package app

import (
	config "AuthInGo/config/env"
	"AuthInGo/controllers"
	db "AuthInGo/db/repositories"
	"AuthInGo/router"
	"AuthInGo/services"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Addr string
}

func NewConfig() *Config {
	port := config.GetString("PORT", "8080")

	return &Config{
		Addr: ":" + port,
	}
}

type Application struct {
	Config Config
	Store  db.Storage
}

func NewApplication(config Config) *Application {
	return &Application{
		Config: config,
		Store:  *db.NewStorage(),
	}
}

func (app *Application) Run() error {
	ur := db.NewUserRepository()
	us := services.NewUserServiceImpl(ur)
	uc := controllers.NewUserController(us)
	urouter := router.NewUserRouter(uc)

	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      router.SetupRouter(urouter),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server is running at", app.Config.Addr)

	return server.ListenAndServe()
}
