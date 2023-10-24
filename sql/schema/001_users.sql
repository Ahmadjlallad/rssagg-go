-- +goose Up
CREATE TABLE users
(
    id         UUID primary key not null,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    name       VARCHAR(255) not null
);
-- +goose Down
DROP TABLE users;