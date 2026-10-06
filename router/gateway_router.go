package router

import (
	env "AuthInGo/config/env"
	"AuthInGo/middleware"
	"AuthInGo/utils"

	"github.com/go-chi/chi/v5"
)

// gatewayRoute maps a public path prefix to a downstream service.
type gatewayRoute struct {
	publicPrefix   string // what the client calls, e.g. /api/bookings
	upstreamURL    string // where the gateway forwards to, e.g. http://localhost:3001
	upstreamPrefix string // the path the service actually serves, e.g. /api/v1/bookings
}

// GatewayRouter registers every proxied microservice behind the JWT middleware.
type GatewayRouter struct {
	routes []gatewayRoute
}

func NewGatewayRouter() Router {
	return &GatewayRouter{
		routes: []gatewayRoute{
			{"/api/bookings", env.GetString("BOOKING_SERVICE_URL", "http://localhost:3001"), "/api/v1/bookings"},
			{"/api/hotels", env.GetString("HOTEL_SERVICE_URL", "http://localhost:3000"), "/api/v1/hotel"},
			{"/api/rooms", env.GetString("HOTEL_SERVICE_URL", "http://localhost:3000"), "/api/v1/room"},
			{"/api/room-categories", env.GetString("HOTEL_SERVICE_URL", "http://localhost:3000"), "/api/v1/room-category"},
			// notification-service is intentionally NOT exposed: it is internal-only
			// and called directly by booking-service when a booking is confirmed
		},
	}
}

func (gr *GatewayRouter) Register(r chi.Router) {
	r.Group(func(pr chi.Router) {
		// authentication is a step of the gateway pipeline: nothing below is
		// proxied unless the JWT is valid
		pr.Use(middleware.JWTAuthMiddleware)

		for _, route := range gr.routes {
			// Mount matches both /api/bookings and /api/bookings/*, all methods.
			// /api/bookings/7 -> {upstream}/api/v1/bookings/7
			pr.Mount(route.publicPrefix, utils.ProxyToService(route.upstreamURL, route.publicPrefix, route.upstreamPrefix))
		}
	})
}
