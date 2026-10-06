package middleware

import (
	dto "AuthInGo/Dto"
	"AuthInGo/utils"
	"context"
	"net/http"
)

type contextKey string

const PayloadKey contextKey = "validatedPayload"

// ValidateLoginRequest validates the login request body
func ValidateLoginRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := &dto.LoginUserRequestDto{}

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

// ValidateCreateUserRequest validates the signup request body
func ValidateCreateUserRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := &dto.CreateUserRequestDto{}

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