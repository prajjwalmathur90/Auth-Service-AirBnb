package db

import (
	"AuthInGo/models"
	"database/sql"
	"strings"
)

type UserRolesRepository interface {
	GetUserRoles(userId int64) ([]*models.Role, error)
	AssignRoleToUser(userId int64, roleId int64) error
	RemoveRoleFromUser(userId int64, roleId int64) error
	GetUserPermissions(userId int64) ([]*models.Permissions, error)
	HasPermission(userId int64, permissionName string) (bool, error)
	HasRole(userId int64, roleName string) (bool, error)
	HasAllRoles(userId int64, roleNames []string) (bool, error)
}

type UserRolesRepositoryImpl struct {
	db *sql.DB
}

func NewUserRolesRepository(db *sql.DB) UserRolesRepository {
	return &UserRolesRepositoryImpl{
		db: db,
	}
}

func (u *UserRolesRepositoryImpl) GetUserRoles(userId int64) ([]*models.Role, error) {
	query := `SELECT roles.id, roles.name, roles.description, roles.created_at, 
			  roles.updated_at 
			  FROM user_roles JOIN roles ON user_roles.role_id = roles.id 
			  WHERE user_roles.user_id = ?`

	rows, err := u.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*models.Role

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		role := &models.Role{}
		if err := rows.Scan(&role.Id, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (u *UserRolesRepositoryImpl) AssignRoleToUser(userId int64, roleId int64) error {
	query := "INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)"

	_, err := u.db.Exec(query, userId, roleId)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRolesRepositoryImpl) RemoveRoleFromUser(userId int64, roleId int64) error {
	query := "DELETE FROM user_roles WHERE user_id = ? AND role_id = ?"

	_, err := u.db.Exec(query, userId, roleId)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserRolesRepositoryImpl) GetUserPermissions(userId int64) ([]*models.Permissions, error) {
	query := `SELECT permissions.id, permissions.name, permissions.description, permissions.resource, 
			  permissions.action, permissions.created_at, permissions.updated_at 
			  FROM user_roles JOIN roles ON user_roles.role_id = roles.id 
			  JOIN role_permissions ON roles.id = role_permissions.role_id 
			  JOIN permissions ON role_permissions.permission_id = permissions.id 
			  WHERE user_roles.user_id = ?`

	rows, err := u.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*models.Permissions

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	for rows.Next() {
		permission := &models.Permissions{}
		if err := rows.Scan(&permission.Id, &permission.Name, &permission.Description, 
			&permission.Resource, &permission.Action, &permission.CreatedAt, 
			&permission.UpdatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (u *UserRolesRepositoryImpl) HasPermission(userId int64, permissionName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_roles JOIN roles ON user_roles.role_id = roles.id 
			  JOIN role_permissions ON roles.id = role_permissions.role_id 
			  JOIN permissions ON role_permissions.permission_id = permissions.id 
			  WHERE user_roles.user_id = ? AND permissions.name = ?)`

	rows, err := u.db.Query(query, userId, permissionName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var hasPermission bool

	if rows.Err() != nil {
		return false, rows.Err()
	}

	for rows.Next() {
		if err := rows.Scan(&hasPermission); err != nil {
			return false, err
		}
	}

	return hasPermission, nil
}

func (u *UserRolesRepositoryImpl) HasRole(userId int64, roleName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM user_roles JOIN roles ON user_roles.role_id = roles.id 
			  WHERE user_roles.user_id = ? AND roles.name = ?)`

	rows, err := u.db.Query(query, userId, roleName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var hasRole bool

	if rows.Err() != nil {
		return false, rows.Err()
	}

	for rows.Next() {
		if err := rows.Scan(&hasRole); err != nil {
			return false, err
		}
	}

	return hasRole, nil
}

func (u *UserRolesRepositoryImpl) HasAllRoles(userId int64, roleNames []string) (bool, error) {
	
	if len(roleNames) == 0 {
		return true, nil
	}
	
	// database/sql cannot bind a slice to a single "?", so build one placeholder per role
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(roleNames)), ",")

	query := `SELECT COUNT(DISTINCT r.id) = ? FROM user_roles ur INNER JOIN roles r ON ur.role_id = r.id 
			  WHERE ur.user_id = ? AND r.name IN (` + placeholders + `) GROUP BY ur.user_id`

	args := make([]interface{}, 0, len(roleNames)+2)
	args = append(args, len(roleNames), userId)
	for _, name := range roleNames {
		args = append(args, name)
	}

	row := u.db.QueryRow(query, args...)

	var hasAllRoles bool
	if err := row.Scan(&hasAllRoles); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return hasAllRoles, nil
}