package db

import (
	"AuthInGo/models"
	"database/sql"
)

type RolePermissionsRepository interface {
	GetRolePermissionById(id int64) (*models.RolePermission, error)
	GetRolePermissionsByRoleId(roleId int64) ([]*models.RolePermission, error)
	AddPermissionToRole(roleId int64, permissionId int64) (*models.RolePermission, error)
	RemovePermissionFromRole(roleId int64, permissionId int64) error
	GetAllRolePermissions() ([]*models.RolePermission, error)
}

type RolePermissionsRepositoryImpl struct {
	db *sql.DB
}

func NewRolePermissionsRepository(db *sql.DB) RolePermissionsRepository {
	return &RolePermissionsRepositoryImpl{
		db: db,
	}
}

func (rp *RolePermissionsRepositoryImpl) GetRolePermissionById(id int64) (*models.RolePermission, error) {
	query := "SELECT id, role_id, permission_id, created_at, updated_at FROM role_permissions WHERE id = ?"

	row := rp.db.QueryRow(query, id)
	rolePermission := &models.RolePermission{}

	if err := row.Scan(&rolePermission.Id, &rolePermission.RoleId, &rolePermission.PermissionId, &rolePermission.CreatedAt, &rolePermission.UpdatedAt); err != nil {
		return nil, err
	}

	return rolePermission, nil
}

func (rp *RolePermissionsRepositoryImpl) GetRolePermissionsByRoleId(roleId int64) ([]*models.RolePermission, error) {
	query := "SELECT id, role_id, permission_id, created_at, updated_at FROM role_permissions WHERE role_id = ?"

	rows, err := rp.db.Query(query, roleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rolePermissions []*models.RolePermission

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		rolePermission := models.RolePermission{}
		if err := rows.Scan(&rolePermission.Id, &rolePermission.RoleId, &rolePermission.PermissionId, &rolePermission.CreatedAt, &rolePermission.UpdatedAt); err != nil {
			return nil, err
		}
		rolePermissions = append(rolePermissions, &rolePermission)
	}

	return rolePermissions, nil
}

func (rp *RolePermissionsRepositoryImpl) AddPermissionToRole(roleId int64, permissionId int64) (*models.RolePermission, error) {
	query := "INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)"

	result, err := rp.db.Exec(query, roleId, permissionId)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return rp.GetRolePermissionById(id)
}

func (rp *RolePermissionsRepositoryImpl) RemovePermissionFromRole(roleId int64, permissionId int64) error {
	query := "DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?"

	_, err := rp.db.Exec(query, roleId, permissionId)

	if err != nil {
		return err
	}

	return nil
}

func (rp *RolePermissionsRepositoryImpl) GetAllRolePermissions() ([]*models.RolePermission, error) {
	query := "SELECT id, role_id, permission_id, created_at, updated_at FROM role_permissions"

	rows, err := rp.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rolePermissions []*models.RolePermission

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		rolePermission := &models.RolePermission{}
		if err := rows.Scan(&rolePermission.Id, &rolePermission.RoleId, &rolePermission.PermissionId, &rolePermission.CreatedAt, &rolePermission.UpdatedAt); err != nil {
			return nil, err
		}
		rolePermissions = append(rolePermissions, rolePermission)
	}

	return rolePermissions, nil
}