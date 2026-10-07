package db

import (
	"AuthInGo/models"
	"database/sql"
)

type PermissionRepository interface {
	GetPermissionById(id int64) (*models.Permissions, error)
	GetPermissionByName(name string) (*models.Permissions, error)
	GetAllPermissions() ([]models.Permissions, error)
	CreatePermission(name string, description string, resource string, action string) (*models.Permissions, error)
	UpdatePermission(id int64, name string, description string, resource string, action string) (*models.Permissions, error)
	DeletePermission(id int64) error
}

type PermissionRepositoryImpl struct {
	db *sql.DB
}

func NewPermissionRepository(db *sql.DB) PermissionRepository {
	return &PermissionRepositoryImpl{
		db: db,
	}
}

func (p *PermissionRepositoryImpl) GetPermissionById(id int64) (*models.Permissions, error) {
	query := "SELECT id, name, description, resource, action, created_at, updated_at FROM permissions WHERE id = ?"

	row := p.db.QueryRow(query, id)
	permission := &models.Permissions{}

	if err := row.Scan(&permission.Id, &permission.Name, &permission.Description, &permission.Resource, &permission.Action, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
		return nil, err
	}

	return permission, nil
}

func (p *PermissionRepositoryImpl) GetPermissionByName(name string) (*models.Permissions, error) {
	query := "SELECT id, name, description, resource, action, created_at, updated_at FROM permissions WHERE name = ?"

	row := p.db.QueryRow(query, name)
	permission := &models.Permissions{}

	if err := row.Scan(&permission.Id, &permission.Name, &permission.Description, &permission.Resource, &permission.Action, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
		return nil, err
	}

	return permission, nil
}

func (p *PermissionRepositoryImpl) GetAllPermissions() ([]models.Permissions, error) {
	query := "SELECT id, name, description, resource, action, created_at, updated_at FROM permissions"

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permissions

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		permission := &models.Permissions{}
		if err := rows.Scan(&permission.Id, &permission.Name, &permission.Description, &permission.Resource, &permission.Action, &permission.CreatedAt, &permission.UpdatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, *permission)
	}

	return permissions, nil
}

func (p *PermissionRepositoryImpl) CreatePermission(name string, description string, resource string, action string) (*models.Permissions, error) {
	query := "INSERT INTO permissions (name, description, resource, action) VALUES (?, ?, ?, ?)"

	result, err := p.db.Exec(query, name, description, resource, action)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return p.GetPermissionById(id)
}

func (p *PermissionRepositoryImpl) UpdatePermission(id int64, name string, description string, resource string, action string) (*models.Permissions, error) {
	query := "UPDATE permissions SET name = ?, description = ?, resource = ?, action = ? WHERE id = ?"

	_, err := p.db.Exec(query, name, description, resource, action, id)
	if err != nil {
		return nil, err
	}

	return p.GetPermissionById(id)
}

func (p *PermissionRepositoryImpl) DeletePermission(id int64) error {
	query := "DELETE FROM permissions WHERE id = ?"

	_, err := p.db.Exec(query, id)
	return err
}