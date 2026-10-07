package controllers

import (
	dto "AuthInGo/dto"
	"AuthInGo/middleware"
	"AuthInGo/services"
	"AuthInGo/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type RoleController struct {
	RoleService services.RoleService
}

func NewRoleController(roleService services.RoleService) *RoleController {
	return &RoleController{
		RoleService: roleService,
	}
}

// writeServiceError maps "not found" service errors to 404, everything else to 500.
func writeServiceError(w http.ResponseWriter, msg string, err error) {
	if err != nil && strings.Contains(err.Error(), "not found") {
		utils.WriteJsonErrorResponse(w, http.StatusNotFound, err.Error(), err)
		return
	}
	utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, msg, err)
}

func (rc *RoleController) GetRoleById(w http.ResponseWriter, r *http.Request) {
	roleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Role Id", err)
		return
	}

	role, err := rc.RoleService.GetRoleById(roleId)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Failed to fetch role", err)
		return
	}

	if role == nil {
		utils.WriteJsonErrorResponse(w, http.StatusNotFound, "Role not found", nil)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "Role fetched successfully", role)
}

func (rc *RoleController) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := rc.RoleService.GetAllRoles()

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Failed to fetch roles", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "Roles fetched successfully", roles)
}

func (rc *RoleController) CreateRole(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.PayloadKey).(*dto.CreateRoleRequestDto)

	role, err := rc.RoleService.CreateRole(payload.Name, payload.Description)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Failed to create role", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusCreated, "Role created successfully", role)
}

func (rc *RoleController) UpdateRole(w http.ResponseWriter, r *http.Request) {
	roleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Role Id", err)
		return
	}

	payload := r.Context().Value(middleware.PayloadKey).(*dto.UpdateRoleRequestDto)

	role, err := rc.RoleService.UpdateRole(roleId, payload.Name, payload.Description)

	if err != nil {
		writeServiceError(w, "Failed to update role", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "Role updated successfully", role)
}

func (rc *RoleController) DeleteRole(w http.ResponseWriter, r *http.Request) {
	roleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Role Id", err)
		return
	}

	if err := rc.RoleService.DeleteRole(roleId); err != nil {
		writeServiceError(w, "Failed to delete role", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "Role deleted successfully", nil)
}

func (rc *RoleController) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Role Id", err)
		return
	}

	permissions, err := rc.RoleService.GetRolePermissions(roleId)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Failed to fetch role permissions", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "Role permissions fetched successfully", permissions)
}

func (rc *RoleController) AddPermissionToRole(w http.ResponseWriter, r *http.Request) {
	roleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Role Id", err)
		return
	}

	payload := r.Context().Value(middleware.PayloadKey).(*dto.AddPermissionToRoleRequestDto)

	rolePermission, err := rc.RoleService.AddPermissionToRole(roleId, payload.PermissionId)

	if err != nil {
		writeServiceError(w, "Failed to add permission to role", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusCreated, "Permission added to role successfully", rolePermission)
}

func (rc *RoleController) RemovePermissionFromRole(w http.ResponseWriter, r *http.Request) {
	roleId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Role Id", err)
		return
	}

	permissionId, err := strconv.ParseInt(chi.URLParam(r, "permissionId"), 10, 64)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Permission Id", err)
		return
	}

	if err := rc.RoleService.RemovePermissionFromRole(roleId, permissionId); err != nil {
		writeServiceError(w, "Failed to remove permission from role", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "Permission removed from role successfully", nil)
}

func (rc *RoleController) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid User Id", err)
		return
	}

	payload := r.Context().Value(middleware.PayloadKey).(*dto.AssignRoleToUserRequestDto)

	if err := rc.RoleService.AssignRoleToUser(userId, payload.RoleId); err != nil {
		writeServiceError(w, "Failed to assign role to user", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "Role assigned to user successfully", nil)
}
