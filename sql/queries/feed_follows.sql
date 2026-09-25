-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, feed_id, user_id, created_at, updated_at)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedFollow :one
SELECT
    *
FROM
    feed_follows
WHERE
    id = $1;

-- name: GetFeedFollowForUser :many
SELECT
    feed_follows.id,
    feeds.name AS feed_name,
    users.name AS user_name
FROM
    feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE
    feed_follows.user_id = $1;

-- name: GetFeedFollowForFeed :many
SELECT
    *
FROM
    feed_follows
WHERE
    feed_id = $1;

-- name: GetFeedFollows :many
SELECT
    *
FROM
    feed_follows;

-- name: DeleteAllFeedFollows :exec
TRUNCATE TABLE feed_follows;
