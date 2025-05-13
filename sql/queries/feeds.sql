-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id) 
VALUES ( 
    $1,
    $2,
    $3,
    $4
) 
RETURNING *;

-- name: GetFeeds :many
SELECT f.name as FeedName, f.url, u.name as UserName
    FROM feeds f
    LEFT JOIN users u
    ON f.user_id = u.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds where url = $1;
