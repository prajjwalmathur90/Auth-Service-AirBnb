package middleware

import (
	env "AuthInGo/config/env"
	"AuthInGo/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// extractToken reads the JWT from the auth cookie first and falls back to the
// Authorization: Bearer header (useful for non-browser clients).
func extractToken(r *http.Request) string {
	if cookie, err := r.Cookie(utils.AuthCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)

		if token == "" {
			utils.WriteJsonErrorResponse(w, http.StatusUnauthorized, "Authentication required", nil)
			return
		}

		claims := jwt.MapClaims{}

		// ParseWithClaims also validates the "exp" claim for us.
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(env.GetString("JWT_SECRET", "TOKEN")), nil
		})

		if err != nil {
			utils.WriteJsonErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token", nil)
			return
		}

		userId, ok := claims["id"].(float64)
		email, emailOk := claims["email"].(string)

		if !ok || !emailOk {
			utils.WriteJsonErrorResponse(w, http.StatusUnauthorized, "Invalid token", nil)
			return
		}

		fmt.Println("Authenticated user id :", int64(userId), " email : ", email)

		cxt := context.WithValue(r.Context(), "userId", strconv.FormatFloat(userId, 'f', 0, 64))
		cxt = context.WithValue(cxt, "email", email)

		next.ServeHTTP(w, r.WithContext(cxt))
	})
}