-- name: CreateYouTubeChannel :exec
INSERT INTO youtube_channels (channel_id, handle, uploads_playlist_id)
VALUES ($1, $2, $3);

-- name: GetYouTubeChannel :one
SELECT * FROM youtube_channels
WHERE channel_id = $1;

-- name: GetYouTubeChannelByHandle :one
SELECT * FROM youtube_channels
WHERE handle = $1;
