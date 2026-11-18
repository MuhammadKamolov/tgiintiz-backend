-- +goose Up
CREATE TABLE roles(
    id BIGINT PRIMARY KEY,
    name VARCHAR(30) UNIQUE NOT NULL
);

INSERT INTO roles(id, name) VALUES
(1, 'client'),
(2, 'admin');

ALTER TABLE users
ADD COLUMN role_id BIGINT NOT NULL DEFAULT 1
REFERENCES roles(id) ON DELETE RESTRICT;

UPDATE users
SET role_id = 2
WHERE role = 'admin';

ALTER TABLE users
DROP COLUMN role;

-- +goose Down
ALTER TABLE users
ADD COLUMN role VARCHAR(30) NOT NULL DEFAULT 'client';

UPDATE users
SET role = 'admin'
WHERE role_id = 2;

ALTER TABLE users
DROP COLUMN role_id;

DROP TABLE IF NOT EXISTS roles;




