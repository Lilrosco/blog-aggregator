-- name: CreatePost :one
INSERT INTO posts (id, feed_id, title, url, description, published_at, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostByUrl :one
SELECT
    *
FROM
    posts
WHERE
    url = $1;

-- name: DeleteAllPosts :exec
DELETE FROM posts;

-- name: GetPosts :many
SELECT
    *
FROM
    posts;

-- name: GetPostsForUser :many
SELECT
    posts.*,
    feeds.name AS feed_name
FROM
    posts
INNER JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
INNER JOIN feeds ON posts.feed_id = feeds.id
WHERE
    feed_follows.user_id = $1
ORDER BY
    posts.published_at DESC
LIMIT $2;

-- SELECT
--     *
-- FROM
--     posts
-- WHERE
--     feed_id = ANY($1::uuid[])
-- ORDER BY
--     updated_at DESC
-- LIMIT $2;
