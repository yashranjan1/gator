-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    created_at date,
    updated_at date,
    name varchar
);

-- +goose Down
DROP TABLE users;
