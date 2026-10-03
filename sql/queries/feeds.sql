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

-- name: MarkFeedFetched :exec
UPDATE
    feeds
SET
    last_fetched_at = NOW(), updated_at = NOW()
WHERE
    id = $1;

-- name: GetNextFeedToFetch :one
SELECT
    *
FROM
    feeds
ORDER BY
    last_fetched_at ASC NULLS FIRST
LIMIT 1;
