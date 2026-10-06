package middleware

import (
	env "AuthInGo/config/env"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}
		
		claims := jwt.MapClaims{}

		_, err := jwt.ParseWithClaims(token, claims, func (token *jwt.Token) (interface{}, error)  {
			return []byte(env.GetString("JWT_SECRET", "TOKEN")), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userId, ok := claims["id"].(float64)
		email, emailOk := claims["email"].(string)

		if !ok || !emailOk {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		fmt.Println("Authenticated user id :", int64(userId), " email : ", email)

		cxt := context.WithValue(r.Context(), "user_id", strconv.FormatFloat(userId, 'f', 0, 64))
		cxt = context.WithValue(cxt, "email", email)

		next.ServeHTTP(w, r.WithContext(cxt))
	})
}