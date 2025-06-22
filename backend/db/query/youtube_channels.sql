-- name: CreateYouTubeChannel :exec
INSERT INTO youtube_channels (channel_id, handle)
VALUES ($1, $2);

-- name: GetYouTubeChannel :one
SELECT * FROM youtube_channels
WHERE channel_id = $1;

-- name: GetYouTubeChannelByHandle :one
SELECT * FROM youtube_channels
WHERE handle = $1;
