-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeed :one
SELECT
    *
FROM
    feeds
WHERE
    url = $1;

-- name: DeleteAllFeeds :exec
TRUNCATE TABLE feeds;

-- name: GetFeeds :many
SELECT
    *
FROM
    feeds;

-- name: GetFeedsWithUserName :many
SELECT
    f.id,
    f.name,
    f.url,
    f.created_at,
    f.updated_at,
    u.name
FROM
    feeds f
LEFT JOIN
    users u ON f.user_id = u.id;
