package middleware

import (
	"AuthInGo/utils"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Every(1*time.Second), 5)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			utils.WriteJsonErrorResponse(w, http.StatusTooManyRequests, "Too many requests", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}