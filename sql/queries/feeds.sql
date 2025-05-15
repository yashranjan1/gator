-- name: CreateFeed :one
INSERT INTO feeds (id, name, updated_at, url, user_id) 
VALUES ( 
    $1,
    $2,
    $3,
    $4,
    $5
) 
RETURNING *;

-- name: GetFeeds :many
SELECT f.name as FeedName, f.url, u.name as UserName
    FROM feeds f
    LEFT JOIN users u
    ON f.user_id = u.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds where url = $1;


-- name: MarkFeedFetched :exec
UPDATE feeds
    SET last_fetched_at = $1,
        updated_at = $1
    WHERE id = $2;

-- name: GetNextFeedToFetch :one
SELECT * 
    FROM feeds
    ORDER BY last_fetched_at ASC NULLS FIRST
    LIMIT 1; 
