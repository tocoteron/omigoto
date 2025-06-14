CREATE EXTENSION vector;

CREATE TABLE youtube_channels (
    channel_id TEXT PRIMARY KEY, -- UC1cnByKe24JjTv38tH_7BYw
    handle TEXT NOT NULL UNIQUE -- @izuho_omi
);

CREATE TABLE youtube_playlists (
    playlist_id TEXT PRIMARY KEY,
    channel_id TEXT NOT NULL REFERENCES youtube_channels (channel_id),
    is_uploads BOOLEAN NOT NULL,
    title TEXT,
    CONSTRAINT youtube_playlists_check CHECK (
        (is_uploads = true AND title IS NULL) OR
        (is_uploads = false AND title IS NOT NULL)
    )
);

CREATE TABLE youtube_videos (
    video_id TEXT PRIMARY KEY,
    playlist_id TEXT NOT NULL REFERENCES youtube_playlists (playlist_id),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    duration INTERVAL SECOND NOT NULL,
    published_at TIMESTAMPTZ NOT NULL
);
