package services

import (
	repositories "AuthInGo/db/repositories"
	"AuthInGo/models"
)

type RoleService interface {
	GetRoleById(id int64) (*models.Role, error)
	GetRoleByName(name string) (*models.Role, error)
	GetAllRoles() ([]*models.Role, error)
	CreateRole(name string, description string) (*models.Role, error)
	UpdateRole(id int64, name string, description string) (*models.Role, error)
	DeleteRole(id int64) error
	GetRolePermissions(roleId int64) ([]*models.RolePermission, error)
	AddPermissionToRole(roleId int64, permissionId int64) (*models.RolePermission, error)
	RemovePermissionFromRole(roleId int64, permissionId int64) error
}

type RoleServiceImpl struct {
	roleRepository          repositories.RoleRepository
	rolePermissionRepsitory repositories.RolePermissionsRepository
}

func NewRoleServiceImpl(roleRepo repositories.RoleRepository, rolePermissionRepo repositories.RolePermissionsRepository) RoleService {
	return &RoleServiceImpl{
		roleRepository:          roleRepo,
		rolePermissionRepsitory: rolePermissionRepo,
	}
}

func (r *RoleServiceImpl) GetRoleById(id int64) (*models.Role, error) {
	return r.roleRepository.GetRoleById(id)
}

func (r *RoleServiceImpl) GetRoleByName(name string) (*models.Role, error) {
	return r.roleRepository.GetRoleByName(name)
}

func (r *RoleServiceImpl) GetAllRoles() ([]*models.Role, error) {
	return r.roleRepository.GetAllRoles()
}

func (r *RoleServiceImpl) CreateRole(name string, description string) (*models.Role, error) {
	return r.roleRepository.CreateRole(name, description)
}

func (r *RoleServiceImpl) UpdateRole(id int64, name string, description string) (*models.Role, error) {
	return r.roleRepository.UpdateRole(id, name, description)
}

func (r *RoleServiceImpl) DeleteRole(id int64) error {
	return r.roleRepository.DeleteRole(id)
}

func (r *RoleServiceImpl) GetRolePermissions(roleId int64) ([]*models.RolePermission, error) {
	return r.rolePermissionRepsitory.GetRolePermissionsByRoleId(roleId)
}

func (r *RoleServiceImpl) AddPermissionToRole(roleId int64, permissionId int64) (*models.RolePermission, error) {
	return r.rolePermissionRepsitory.AddPermissionToRole(roleId, permissionId)
}

func (r *RoleServiceImpl) RemovePermissionFromRole(roleId int64, permissionId int64) error {
	return r.rolePermissionRepsitory.RemovePermissionFromRole(roleId, permissionId)
}