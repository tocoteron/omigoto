-- name: CreateYouTubePlaylist :exec
INSERT INTO youtube_playlists (playlist_id, channel_id, title)
VALUES ($1, $2, $3);

-- name: GetYouTubePlaylist :one
SELECT * FROM youtube_playlists
WHERE playlist_id = $1;

-- name: ListPlaylists :many
SELECT * FROM youtube_playlists
WHERE playlist_id = ANY(@playlist_ids::text[]);

-- name: ListPlaylistIDsByChannel :many
SELECT playlist_id FROM youtube_playlists
WHERE channel_id = $1;
