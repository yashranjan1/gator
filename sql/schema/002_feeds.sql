-- +goose Up
CREATE TABLE feeds (
    id uuid PRIMARY KEY,
    name VARCHAR NOT NULL,
    url VARCHAR NOT NULL,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE(url)
);

-- +goose Down
DROP TABLE feeds;
