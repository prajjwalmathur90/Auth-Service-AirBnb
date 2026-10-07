package router

import (
	"AuthInGo/controllers"
	"AuthInGo/middleware"

	"github.com/go-chi/chi/v5"
)

type RoleRouter struct {
	roleController *controllers.RoleController
}

func NewRoleRouter(_roleController *controllers.RoleController) Router {
	return &RoleRouter{
		roleController: _roleController,
	}
}

func (rr *RoleRouter) Register(r chi.Router) {
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.JWTAuthMiddleware)
		pr.Use(middleware.RequireAllRoles("admin"))

		pr.Get("/{id}", rr.roleController.GetRoleById)
		pr.Get("/", rr.roleController.GetAllRoles)
		pr.With(middleware.ValidateCreateRoleRequest).Post("/", rr.roleController.CreateRole)
		pr.With(middleware.ValidateUpdateRoleRequest).Put("/{id}", rr.roleController.UpdateRole)
		pr.Delete("/{id}", rr.roleController.DeleteRole)

		pr.Get("/{id}/permissions", rr.roleController.GetRolePermissions)
		pr.With(middleware.ValidateAddPermissionToRoleRequest).Post("/{id}/permissions", rr.roleController.AddPermissionToRole)
		pr.Delete("/{id}/permissions/{permissionId}", rr.roleController.RemovePermissionFromRole)

		pr.With(middleware.ValidateAssignRoleToUserRequest).Post("/users/{userId}", rr.roleController.AssignRoleToUser)
	})
}
