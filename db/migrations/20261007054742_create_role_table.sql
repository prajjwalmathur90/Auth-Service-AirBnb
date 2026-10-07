-- +goose Up
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (name, description) VALUES
('admin', 'Administrator with full access'),
('user', 'Regular user with standard access'),
('moderator', 'Moderator with elevated privilages');

-- +goose Down
DROP TABLE IF EXISTS roles;
