-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (
        id,
        created_at,
        updated_at,
        user_id,
        feed_id
    )
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) RETURNING *
)
SELECT inserted_feed_follow.*,
    feeds.name as feed_name,
    users.name as user_name
FROM inserted_feed_follow
    INNER JOIN users 
    ON inserted_feed_follow.user_id = users.id
    INNER JOIN feeds 
    ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT * 
    FROM feed_follows 
    LEFT JOIN feeds
    ON feeds.id = feed_id
    WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows ff WHERE ff.user_id = $1 AND ff.feed_id = (
    SELECT id FROM feeds f WHERE f.url = $2
);
