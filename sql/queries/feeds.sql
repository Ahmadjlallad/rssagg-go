-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, created_at, updated_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeed :one
select *
from feeds
where id = $1;

-- name: ListFeeds :many
select *
from feeds
order by id;

-- name: GetNextFeedToFetch :one
select *
from feeds
order by last_fetch_at nulls first
limit 1;

-- name: MarkFeedAsFetched :one
update feeds
set last_fetch_at = now(),
    updated_at=now()
where id = $1
returning *;