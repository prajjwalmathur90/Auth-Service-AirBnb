package router

import (
	"AuthInGo/controllers"
	"AuthInGo/middleware"

	"github.com/go-chi/chi/v5"
)

type Router interface {
	Register(r chi.Router)
}

func SetupRouter(UserRouter Router, GatewayRouter Router) *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.RateLimitMiddleware)
	chiRouter.Get("/ping", controllers.PingHandler)

	// /api/auth/* is served by this process directly (no proxy hop)
	chiRouter.Route("/api/auth", func(r chi.Router) {
		UserRouter.Register(r)
	})

	// /api/bookings/*, /api/hotels/* -> JWT -> reverse proxy
	GatewayRouter.Register(chiRouter)

	return chiRouter
}