-- +goose Up

CREATE TABLE IF NOT EXISTS role_permissions (
    id SERIAL PRIMARY KEY,
    role_id BIGINT UNSIGNED,
    permission_id BIGINT UNSIGNED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);


INSERT INTO role_permissions (role_id, permission_id) 
SELECT 1, id FROM permissions; -- Admin has all the perms

INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions where name IN ('user:read');

-- +goose Down

DROP TABLE IF EXISTS role_permissions;
