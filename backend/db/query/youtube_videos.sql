-- name: CreateYouTubeVideo :exec
INSERT INTO youtube_videos (
    video_id, title, description, duration,
    thumbnail_default_url,thumbnail_medium_url, thumbnail_high_url, thumbnail_standard_url, thumbnail_maxres_url,
    published_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetYouTubeVideo :one
SELECT * FROM youtube_videos
WHERE video_id = $1;

-- name: ListYouTubeVideos :many
SELECT * FROM youtube_videos
WHERE video_id = ANY(@video_ids::text[]);

-- name: CreateYouTubeVideoLiveStreamingDetails :exec
INSERT INTO youtube_video_live_streaming_details (
    video_id, actual_start_time, actual_end_time, scheduled_start_time
)
VALUES ($1, $2, $3, $4);

-- name: GetYouTubeVideoLiveStreamingDetails :one
SELECT * FROM youtube_video_live_streaming_details
WHERE video_id = $1;
