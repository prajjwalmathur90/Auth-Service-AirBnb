package middleware

import (
	"AuthInGo/utils"
	"context"
	"net/http"
)

type contextKey string

const PayloadKey contextKey = "validatedPayload"

// ValidateBody returns a route-specific middleware that:
// 1. Parses the JSON body into the struct returned by newPayload()
// 2. Validates the struct using the validator
// 3. Stores the validated payload in the request context
func ValidateBody(newPayload func() any) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload := newPayload()

			if err := utils.ReadJsonBody(r, payload); err != nil {
				utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Json Body", err)
				return
			}

			if err := utils.Validator.Struct(payload); err != nil {
				utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Request Body", err)
				return
			}

			ctx := context.WithValue(r.Context(), PayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}