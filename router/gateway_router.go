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
	requireAuth    bool
}

// GatewayRouter registers proxied microservices behind JWT middleware.
type GatewayRouter struct {
	routes []gatewayRoute
}

func NewGatewayRouter() Router {
	return &GatewayRouter{
		routes: []gatewayRoute{
			{"/api/bookings", env.GetString("BOOKING_SERVICE_URL", "http://localhost:3001"), "/api/v1/bookings", true},
			{"/api/hotels", env.GetString("HOTEL_SERVICE_URL", "http://localhost:3000"), "/api/v1/hotel", false},
			{"/api/rooms", env.GetString("HOTEL_SERVICE_URL", "http://localhost:3000"), "/api/v1/room", false},
			{"/api/room-categories", env.GetString("HOTEL_SERVICE_URL", "http://localhost:3000"), "/api/v1/room-category", false},
			{"/api/room-generation", env.GetString("HOTEL_SERVICE_URL", "http://localhost:3000"), "/api/v1/room-generation", false},
			// notification-service is intentionally NOT exposed: it is internal-only
			// and called directly by booking-service when a booking is confirmed
		},
	}
}

func (gr *GatewayRouter) Register(r chi.Router) {
	// 1. Mandatory authenticated routes
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.JWTAuthMiddleware)

		for _, route := range gr.routes {
			if route.requireAuth {
				pr.Mount(route.publicPrefix, utils.ProxyToService(route.upstreamURL, route.publicPrefix, route.upstreamPrefix))
			}
		}
	})

	// 2. Public / Optional auth routes (e.g. browsing rooms, checking availability)
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.OptionalJWTAuthMiddleware)

		for _, route := range gr.routes {
			if !route.requireAuth {
				pr.Mount(route.publicPrefix, utils.ProxyToService(route.upstreamURL, route.publicPrefix, route.upstreamPrefix))
			}
		}
	})
}
