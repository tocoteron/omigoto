CREATE EXTENSION vector;

CREATE TABLE youtube_channels (
    channel_id TEXT PRIMARY KEY,       -- UC1cnByKe24JjTv38tH_7BYw
    handle TEXT NOT NULL UNIQUE,       -- @izuho_omi
    uploads_playlist_id TEXT NOT NULL, -- UU1cnByKe24JjTv38tH_7BYw
    CONSTRAINT youtube_channels_uploads_playlist_id_fkey FOREIGN KEY (uploads_playlist_id) REFERENCES youtube_playlists (playlist_id) DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE youtube_playlists (
    playlist_id TEXT PRIMARY KEY,
    channel_id TEXT NOT NULL,
    title TEXT NOT NULL,
    CONSTRAINT youtube_playlists_channel_id_fkey FOREIGN KEY (channel_id) REFERENCES youtube_channels (channel_id) DEFERRABLE INITIALLY DEFERRED
);

CREATE TABLE youtube_videos (
    video_id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    duration INTERVAL SECOND NOT NULL,
    thumbnail_default_url TEXT,  -- 120x90
    thumbnail_medium_url TEXT,   -- 320x180
    thumbnail_high_url TEXT,     -- 480x360
    thumbnail_standard_url TEXT, -- 640x480
    thumbnail_maxres_url TEXT,   -- 1280x720
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
