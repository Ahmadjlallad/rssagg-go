-- +goose Up
CREATE TABLE feed_follows
(
    id         UUID primary key not null,
    created_at TIMESTAMP        not null,
    updated_at TIMESTAMP        not null,
    user_id    UUID             NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    feed_id    UUID             NOT NULL REFERENCES feeds (id) ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;