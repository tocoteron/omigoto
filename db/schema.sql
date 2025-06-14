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
    thumbnail_default_url TEXT NOT NULL,  -- 120x90
    thumbnail_medium_url TEXT NOT NULL,   -- 320x180
    thumbnail_high_url TEXT NOT NULL,     -- 480x360
    thumbnail_standard_url TEXT NOT NULL, -- 640x480
    thumbnail_maxres_url TEXT NOT NULL,   -- 1280x720
    published_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE youtube_video_live_streaming_details (
    video_id TEXT PRIMARY KEY REFERENCES youtube_videos (video_id),
    actual_start_time TIMESTAMPTZ NOT NULL,
    actual_end_time TIMESTAMPTZ NOT NULL,
    scheduled_start_time TIMESTAMPTZ NOT NULL
);

CREATE TABLE youtube_playlist_videos (
    playlist_id TEXT NOT NULL REFERENCES youtube_playlists (playlist_id),
    video_id TEXT NOT NULL REFERENCES youtube_videos (video_id),
    PRIMARY KEY (playlist_id, video_id)
);
