package middleware

import (
	dbConfig "AuthInGo/config/db"
	env "AuthInGo/config/env"
	repo "AuthInGo/db/repositories"
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

func RequireAllRoles(roles ...string) func(http.Handler) http.Handler {
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userIdStr := r.Context().Value("userId").(string)
			userId, err := strconv.ParseInt(userIdStr, 10, 64)
			if err != nil {
				utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid user id", nil)
				return
			}

			dbConn, dbErr := dbConfig.SetupDB()
			if dbErr != nil {
				utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Database connection error", nil)
				return
			}

			defer dbConn.Close()

			urr := repo.NewUserRolesRepository(dbConn)
			hasAllRoles, hasAllRolesErr := urr.HasAllRoles(userId, roles)
			if hasAllRolesErr != nil {
				utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Internal Server error", nil)
				return
			}
			if !hasAllRoles {
				utils.WriteJsonErrorResponse(w, http.StatusForbidden, "Insufficient permissions", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}