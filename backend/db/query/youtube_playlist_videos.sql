-- name: CreateYouTubePlaylistVideo :exec
INSERT INTO youtube_playlist_videos (playlist_id, video_id)
VALUES ($1, $2);

-- name: ListYouTubePlaylistVideoIDs :many
SELECT video_id FROM youtube_playlist_videos
WHERE playlist_id = $1;
