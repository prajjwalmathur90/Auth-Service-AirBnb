package middleware

import (
	dto "AuthInGo/dto"
	"AuthInGo/utils"
	"context"
	"net/http"
)

func validateBody(next http.Handler, newPayload func() any) http.Handler {
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

// ValidateCreateRoleRequest validates the create role request body
func ValidateCreateRoleRequest(next http.Handler) http.Handler {
	return validateBody(next, func() any { return &dto.CreateRoleRequestDto{} })
}

// ValidateUpdateRoleRequest validates the update role request body
func ValidateUpdateRoleRequest(next http.Handler) http.Handler {
	return validateBody(next, func() any { return &dto.UpdateRoleRequestDto{} })
}

// ValidateAddPermissionToRoleRequest validates the add permission request body
func ValidateAddPermissionToRoleRequest(next http.Handler) http.Handler {
	return validateBody(next, func() any { return &dto.AddPermissionToRoleRequestDto{} })
}
