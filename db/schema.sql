CREATE EXTENSION vector;

CREATE TABLE youtube_channels (
    channel_id TEXT PRIMARY KEY, -- UC1cnByKe24JjTv38tH_7BYw
    handle TEXT NOT NULL UNIQUE -- @izuho_omi
);

CREATE TABLE youtube_playlists (
    channel_id TEXT NOT NULL REFERENCES youtube_channels (channel_id),
    playlist_id TEXT PRIMARY KEY,
    is_uploads BOOLEAN NOT NULL,
    title TEXT,
    CHECK (
        (is_uploads = true AND title IS NULL) OR
        (is_uploads = false AND title IS NOT NULL)
    )
);
