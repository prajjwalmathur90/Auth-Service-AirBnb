package app

import (
	dbConfig "AuthInGo/config/db"
	config "AuthInGo/config/env"
	"AuthInGo/controllers"
	db "AuthInGo/db/repositories"
	repo "AuthInGo/db/repositories"
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

	db, err := dbConfig.SetupDB()

	if err != nil {
		fmt.Println("error : Database connection not found")
	}

	ur := repo.NewUserRepository(db)
	urr := repo.NewUserRolesRepository(db)
	us := services.NewUserServiceImpl(ur, urr)
	uc := controllers.NewUserController(us)
	urouter := router.NewUserRouter(uc)

	grouter := router.NewGatewayRouter()

	rr := repo.NewRoleRepository(db)
	rrps := repo.NewRolePermissionsRepository(db)
	rp := repo.NewPermissionRepository(db)
	rrs := services.NewRoleServiceImpl(rr, rrps, rp, ur, urr)
	rc := controllers.NewRoleController(rrs)
	rrouter := router.NewRoleRouter(rc)

	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      router.SetupRouter(urouter, grouter, rrouter),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server is running at port", app.Config.Addr)

	return server.ListenAndServe()
}
