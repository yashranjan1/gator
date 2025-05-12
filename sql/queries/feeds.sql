-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id) 
VALUES ( 
    $1,
    $2,
    $3
) 
RETURNING *;

-- name: GetFeeds :many
SELECT f.name as FeedName, f.url, u.name as UserName
    FROM feeds f
    LEFT JOIN users u
    ON f.user_id = u.id;
