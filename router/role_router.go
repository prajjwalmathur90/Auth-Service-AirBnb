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
	r.Get("/{id}", rr.roleController.GetRoleById)
	r.Get("/", rr.roleController.GetAllRoles)
	r.With(middleware.ValidateCreateRoleRequest).Post("/", rr.roleController.CreateRole)
	r.With(middleware.ValidateUpdateRoleRequest).Put("/{id}", rr.roleController.UpdateRole)
	r.Delete("/{id}", rr.roleController.DeleteRole)

	r.Get("/{id}/permissions", rr.roleController.GetRolePermissions)
	r.With(middleware.ValidateAddPermissionToRoleRequest).Post("/{id}/permissions", rr.roleController.AddPermissionToRole)
	r.Delete("/{id}/permissions/{permissionId}", rr.roleController.RemovePermissionFromRole)
}
