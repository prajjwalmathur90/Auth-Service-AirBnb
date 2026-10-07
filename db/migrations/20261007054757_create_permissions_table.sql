-- +goose Up
CREATE TABLE IF NOT EXISTS permissions(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO permissions (name, description, resource, action) VALUES
('user:read', 'Permission to read user data', 'user', 'read'),
('user:write', 'Permission to create/update user data', 'user', 'write'),
('user:delete', 'Permission to delete user data', 'user', 'delete'),
('role:read', 'Permission to read role data', 'role', 'read'),
('role:write', 'Permission to create/update role data', 'role', 'write'),
('role:delete', 'Permission to delete role data', 'role', 'delete'),
('role:manage', 'Permission to manage role data', 'role', 'manage'),
('permission:read', 'Permission to read permission data', 'permission', 'read'),
('permission:write', 'Permission to create/update permission data', 'permission', 'write'),
('permission:delete', 'Permission to delete permission data', 'permission', 'delete'),
('permission:manage', 'Permission to manage permission data', 'permission', 'manage');


-- +goose Down
DROP TABLE IF EXISTS permissions;
