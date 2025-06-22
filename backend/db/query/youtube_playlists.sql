-- name: CreateYouTubePlaylist :exec
INSERT INTO youtube_playlists (playlist_id, channel_id, is_uploads, title)
VALUES ($1, $2, $3, $4);

-- name: GetYouTubePlaylist :one
SELECT * FROM youtube_playlists
WHERE playlist_id = $1;

-- name: ListYouTubePlaylistsByChannel :many
SELECT * FROM youtube_playlists
WHERE channel_id = $1;

-- name: GetUploadsPlaylistByChannel :one
SELECT * FROM youtube_playlists
WHERE channel_id = $1 AND is_uploads = true;
